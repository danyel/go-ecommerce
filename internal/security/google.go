package security

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	JWT "github.com/golang-jwt/jwt/v5"
)

const googleCertificatesURL = "https://www.googleapis.com/oauth2/v1/certs"

type GoogleProvider struct {
	clientID   string
	httpClient *http.Client
	mu         sync.RWMutex
	keys       map[string]*rsa.PublicKey
	expiry     time.Time
}

func NewGoogleProvider(clientID string) *GoogleProvider {
	return &GoogleProvider{
		clientID:   strings.TrimSpace(clientID),
		httpClient: &http.Client{Timeout: 5 * time.Second},
		keys:       make(map[string]*rsa.PublicKey),
	}
}

func (provider *GoogleProvider) Name() string { return "google" }

type googleClaims struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	EmailVerified bool   `json:"email_verified"`
	JWT.RegisteredClaims
}

func (provider *GoogleProvider) Verify(ctx context.Context, token string) (Identity, error) {
	if provider.clientID == "" {
		return Identity{}, errors.New("google client id is not configured")
	}
	var claims googleClaims
	parsed, err := JWT.ParseWithClaims(token, &claims, func(parsed *JWT.Token) (any, error) {
		if _, ok := parsed.Method.(*JWT.SigningMethodRSA); !ok {
			return nil, errors.New("google token is not RSA signed")
		}
		return provider.key(ctx, parsed.Header["kid"])
	}, JWT.WithValidMethods([]string{JWT.SigningMethodRS256.Alg()}))
	if err != nil || !parsed.Valid {
		return Identity{}, fmt.Errorf("invalid google id token: %w", err)
	}
	if claims.Issuer != "https://accounts.google.com" && claims.Issuer != "accounts.google.com" {
		return Identity{}, errors.New("invalid google token issuer")
	}
	if len(claims.Audience) != 1 || claims.Audience[0] != provider.clientID {
		return Identity{}, errors.New("invalid google token audience")
	}
	if !claims.EmailVerified || claims.Subject == "" {
		return Identity{}, errors.New("google account is not verified")
	}
	return Identity{Provider: provider.Name(), Subject: claims.Subject, Email: claims.Email, Name: claims.Name}, nil
}

func (provider *GoogleProvider) key(ctx context.Context, kid any) (*rsa.PublicKey, error) {
	keyID, ok := kid.(string)
	if !ok || keyID == "" {
		return nil, errors.New("google token has no key id")
	}
	provider.mu.RLock()
	key, valid := provider.keys[keyID], time.Now().Before(provider.expiry)
	provider.mu.RUnlock()
	if valid {
		return key, nil
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, googleCertificatesURL, nil)
	if err != nil {
		return nil, err
	}
	response, err := provider.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google certificate endpoint returned %s", response.Status)
	}
	var certificates map[string]string
	if err := json.NewDecoder(response.Body).Decode(&certificates); err != nil {
		return nil, err
	}
	keys := make(map[string]*rsa.PublicKey, len(certificates))
	for id, certificate := range certificates {
		block, _ := pem.Decode([]byte(certificate))
		if block == nil {
			return nil, errors.New("invalid google certificate")
		}
		parsed, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		publicKey, ok := parsed.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("google certificate is not RSA")
		}
		keys[id] = publicKey
	}
	provider.mu.Lock()
	provider.keys = keys
	provider.expiry = time.Now().Add(30 * time.Minute)
	key = provider.keys[keyID]
	provider.mu.Unlock()
	if key == nil {
		return nil, errors.New("unknown google certificate key id")
	}
	return key, nil
}
