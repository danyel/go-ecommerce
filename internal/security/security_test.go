package security

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	ApplicationMiddleware "github.com/danyel/ecommerce/cmd/middleware"
	"github.com/stretchr/testify/require"
)

type testProvider struct{}

func (testProvider) Name() string { return "test" }
func (testProvider) Verify(context.Context, string) (Identity, error) {
	return Identity{Provider: "test", Subject: "subject"}, nil
}

func TestProviderRegistryRegistersProviders(t *testing.T) {
	registry := NewProviderRegistry(testProvider{})
	provider, ok := registry.Provider("test")
	require.True(t, ok)
	require.Equal(t, "test", provider.Name())
	_, ok = registry.Provider("missing")
	require.False(t, ok)
}

func TestRequireAuthenticationProtectsMutationsWithoutAffectingPublicRoutes(t *testing.T) {
	issuer := NewEncryptedSessionIssuer("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	protected := RequireAuthentication(issuer)(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusNoContent)
	}))
	public := http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusOK)
	})

	publicResponse := httptest.NewRecorder()
	public.ServeHTTP(publicResponse, httptest.NewRequest(http.MethodGet, "/api/product/v1/products", nil))
	require.Equal(t, http.StatusOK, publicResponse.Code)

	protectedResponse := httptest.NewRecorder()
	protected.ServeHTTP(protectedResponse, httptest.NewRequest(http.MethodPut, "/api/shopping-basket/v1/shopping-baskets/id", nil))
	require.Equal(t, http.StatusUnauthorized, protectedResponse.Code)

	token, err := issuer.Issue(&ApplicationMiddleware.UserClaims{UserID: "user", Roles: []string{"USER"}})
	require.NoError(t, err)
	request := httptest.NewRequest(http.MethodPut, "/api/shopping-basket/v1/shopping-baskets/id", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	authenticatedResponse := httptest.NewRecorder()
	protected.ServeHTTP(authenticatedResponse, request)
	require.Equal(t, http.StatusNoContent, authenticatedResponse.Code)
}
