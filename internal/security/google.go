package security

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	IO "io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SSOProvider handles generic OpenID Connect / Enterprise SSO authentication
type SSOProvider struct {
	issuerURL  string
	jwksURL    string
	clientID   string
	httpClient *http.Client
	mu         sync.RWMutex
	keys       map[string]*rsa.PublicKey
	expiry     time.Time
}

// NewSSOProvider creates an SSO validator.
// issuerURL: The base URL of your IdP (e.g., "https://yourcompany.com")
// jwksURL: The IdP's JSON Web Key Set URL (e.g., "https://yourcompany.com/oauth2/v1/keys")
func NewSSOProvider(issuerURL, jwksURL, clientID string) *SSOProvider {
	return &SSOProvider{
		issuerURL:  strings.TrimSuffix(strings.TrimSpace(issuerURL), "/"),
		jwksURL:    strings.TrimSpace(jwksURL),
		clientID:   strings.TrimSpace(clientID),
		httpClient: &http.Client{Timeout: 5 * time.Second},
		keys:       make(map[string]*rsa.PublicKey),
	}
}

func (provider *SSOProvider) Name() string { return "sso" }

type ssoClaims struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	EmailVerified bool   `json:"email_verified"`
	jwt.RegisteredClaims
}

// Verify decodes and validates the enterprise SSO ID Token
func (provider *SSOProvider) Verify(ctx context.Context, token string) (Identity, error) {
	if provider.clientID == "" || provider.jwksURL == "" {
		return Identity{}, errors.New("sso provider is not fully configured")
	}

	var claims ssoClaims
	parsed, err := jwt.ParseWithClaims(token, &claims, func(parsed *jwt.Token) (any, error) {
		if _, ok := parsed.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", parsed.Header["alg"])
		}
		return provider.key(ctx, parsed.Header["kid"])
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))

	if err != nil || !parsed.Valid {
		return Identity{}, fmt.Errorf("invalid sso token: %w", err)
	}

	// Validate dynamic issuer matching our configured enterprise IdP
	cleanIssuer := strings.TrimSuffix(claims.Issuer, "/")
	if cleanIssuer != provider.issuerURL {
		return Identity{}, errors.New("invalid token issuer for enterprise sso")
	}

	// Validate target audience match
	if len(claims.Audience) != 1 || claims.Audience[0] != provider.clientID {
		return Identity{}, errors.New("invalid token audience mapping")
	}

	return Identity{
		Provider: provider.Name(),
		Subject:  claims.Subject,
		Email:    claims.Email,
		Name:     claims.Name,
	}, nil
}

// Key dynamic retrieval from the enterprise IdP JWKS endpoint
func (provider *SSOProvider) key(ctx context.Context, kid any) (*rsa.PublicKey, error) {
	keyID, ok := kid.(string)
	if !ok || keyID == "" {
		return nil, errors.New("sso token lacks a key ID (kid)")
	}

	provider.mu.RLock()
	key, valid := provider.keys[keyID], time.Now().Before(provider.expiry)
	provider.mu.RUnlock()
	if valid && key != nil {
		return key, nil
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, provider.jwksURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := provider.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body IO.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jwks endpoint returned status: %s", response.Status)
	}

	// JWKS formats vary between standard cert arrays and JWK JSON web keys.
	// This parser handles standard raw PEM cert mappings typically exposed at x509 cert endpoints.
	var certificates map[string]string
	if err := json.NewDecoder(response.Body).Decode(&certificates); err != nil {
		return nil, err
	}

	keys := make(map[string]*rsa.PublicKey, len(certificates))
	for id, certificate := range certificates {
		block, _ := pem.Decode([]byte(certificate))
		if block == nil {
			continue
		}
		parsed, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			continue
		}
		publicKey, ok := parsed.PublicKey.(*rsa.PublicKey)
		if !ok {
			continue
		}
		keys[id] = publicKey
	}

	provider.mu.Lock()
	provider.keys = keys
	provider.expiry = time.Now().Add(15 * time.Minute) // Lowered to 15m for enterprise key-rotation agility
	key = provider.keys[keyID]
	provider.mu.Unlock()

	if key == nil {
		return nil, errors.New("matching signing key not found in identity provider metadata")
	}
	return key, nil
}
