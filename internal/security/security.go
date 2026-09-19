package security

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	ApplicationMiddleware "github.com/danyel/ecommerce/cmd/middleware"
	JWT "github.com/golang-jwt/jwt/v5"
)

type Identity struct {
	Provider string
	Subject  string
	Email    string
	Name     string
}

type IdentityProvider interface {
	Name() string
	Verify(context.Context, string) (Identity, error)
}

type ProviderRegistry struct {
	providers map[string]IdentityProvider
}

func NewProviderRegistry(providers ...IdentityProvider) *ProviderRegistry {
	registry := &ProviderRegistry{providers: make(map[string]IdentityProvider, len(providers))}
	for _, provider := range providers {
		registry.Register(provider)
	}
	return registry
}

func (registry *ProviderRegistry) Register(provider IdentityProvider) {
	if provider != nil && provider.Name() != "" {
		registry.providers[provider.Name()] = provider
	}
}

func (registry *ProviderRegistry) Provider(name string) (IdentityProvider, bool) {
	provider, ok := registry.providers[name]
	return provider, ok
}

type UserLinker interface {
	Link(context.Context, Identity) (string, error)
}

// MemoryUserLinker is deliberately non-persistent. It is safe for development
// because it never accepts an unverified identity and makes its limitation explicit.
type MemoryUserLinker struct {
	users map[string]string
}

func NewMemoryUserLinker() *MemoryUserLinker {
	return &MemoryUserLinker{users: make(map[string]string)}
}

func (linker *MemoryUserLinker) Link(_ context.Context, identity Identity) (string, error) {
	if identity.Provider == "" || identity.Subject == "" {
		return "", errors.New("identity is incomplete")
	}
	key := identity.Provider + ":" + identity.Subject
	if userID, ok := linker.users[key]; ok {
		return userID, nil
	}
	userID := key
	linker.users[key] = userID
	return userID, nil
}

type SessionIssuer interface {
	Issue(*ApplicationMiddleware.UserClaims) (string, error)
	Parse(string) (*ApplicationMiddleware.UserClaims, error)
}

type EncryptedSessionIssuer struct {
	secret string
}

func NewEncryptedSessionIssuer(secret string) *EncryptedSessionIssuer {
	return &EncryptedSessionIssuer{secret: secret}
}

func (issuer *EncryptedSessionIssuer) Issue(claims *ApplicationMiddleware.UserClaims) (string, error) {
	key, err := hex.DecodeString(issuer.secret)
	if err != nil || len(key) < 32 {
		return "", errors.New("session secret must be at least 32 bytes of hex")
	}
	now := time.Now()
	tokenClaims := sessionClaims{
		UserID: claims.UserID,
		Roles:  claims.Roles,
		RegisteredClaims: JWT.RegisteredClaims{
			Issuer:    "go-commerce",
			Subject:   claims.UserID,
			IssuedAt:  JWT.NewNumericDate(now),
			ExpiresAt: JWT.NewNumericDate(now.Add(8 * time.Hour)),
		},
	}
	return JWT.NewWithClaims(JWT.SigningMethodHS256, tokenClaims).SignedString(key)
}

func (issuer *EncryptedSessionIssuer) Parse(token string) (*ApplicationMiddleware.UserClaims, error) {
	key, err := hex.DecodeString(issuer.secret)
	if err != nil || len(key) < 32 {
		return nil, errors.New("session secret must be at least 32 bytes of hex")
	}
	var claims sessionClaims
	parsed, err := JWT.ParseWithClaims(token, &claims, func(parsed *JWT.Token) (any, error) {
		if parsed.Method != JWT.SigningMethodHS256 {
			return nil, errors.New("unexpected session signing method")
		}
		return key, nil
	}, JWT.WithValidMethods([]string{JWT.SigningMethodHS256.Alg()}), JWT.WithIssuer("go-commerce"))
	if err != nil || !parsed.Valid {
		// Existing development clients used the authenticated legacy payload.
		// Accept it only after the signed session parser fails; malformed and
		// unauthenticated values are still rejected.
		legacyClaims, legacyErr := ApplicationMiddleware.DecryptClaims(token, issuer.secret)
		if legacyErr == nil && legacyClaims.UserID != "" {
			return legacyClaims, nil
		}
		return nil, fmt.Errorf("invalid application session: %w", err)
	}
	return &ApplicationMiddleware.UserClaims{UserID: claims.UserID, Roles: claims.Roles}, nil
}

type sessionClaims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	JWT.RegisteredClaims
}

type contextKey string

const ClaimsContextKey contextKey = "authenticated_claims"

func ClaimsFromContext(ctx context.Context) (*ApplicationMiddleware.UserClaims, bool) {
	claims, ok := ctx.Value(ClaimsContextKey).(*ApplicationMiddleware.UserClaims)
	return claims, ok
}

func RequireAuthentication(issuer SessionIssuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			header := strings.TrimSpace(request.Header.Get(ApplicationMiddleware.Authorization))
			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(response, "Unauthorized: Missing or malformed token", http.StatusUnauthorized)
				return
			}
			token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
			if token == "" {
				http.Error(response, "Unauthorized: Missing or malformed token", http.StatusUnauthorized)
				return
			}
			claims, err := issuer.Parse(token)
			if err != nil || claims.UserID == "" {
				http.Error(response, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(request.Context(), ClaimsContextKey, claims)
			next.ServeHTTP(response, request.WithContext(ctx))
		})
	}
}
