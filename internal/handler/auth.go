package handler

import (
	JSON "encoding/json"
	Http "net/http"

	ApplicationMiddleware "github.com/danyel/ecommerce/cmd/middleware"
	Security "github.com/danyel/ecommerce/internal/security"
)

type AuthWebHandler struct {
	registry *Security.ProviderRegistry
	linker   Security.UserLinker
	issuer   Security.SessionIssuer
}

// Renamed to represent a generic SSO login structure
type ssoLoginRequest struct {
	IDToken string `json:"id_token"`
}

type loginResponse struct {
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	Provider string `json:"provider"`
}

func NewAuthWebHandler(registry *Security.ProviderRegistry, linker Security.UserLinker, issuer Security.SessionIssuer) *AuthWebHandler {
	return &AuthWebHandler{registry: registry, linker: linker, issuer: issuer}
}

// HandleSSOLogin replaces HandleGoogleLogin to support any configured enterprise provider
func (handler *AuthWebHandler) HandleSSOLogin(response Http.ResponseWriter, request *Http.Request) {
	var login ssoLoginRequest
	if err := JSON.NewDecoder(request.Body).Decode(&login); err != nil || login.IDToken == "" {
		Http.Error(response, "Invalid SSO login request", Http.StatusBadRequest)
		return
	}

	// Looks for the generic "sso" registry configuration we established earlier
	provider, ok := handler.registry.Provider("sso")
	if !ok {
		Http.Error(response, "SSO login is not configured", Http.StatusServiceUnavailable)
		return
	}

	identity, err := provider.Verify(request.Context(), login.IDToken)
	if err != nil {
		Http.Error(response, "SSO identity could not be verified", Http.StatusUnauthorized)
		return
	}

	userID, err := handler.linker.Link(request.Context(), identity)
	if err != nil {
		Http.Error(response, "Could not link account", Http.StatusInternalServerError)
		return
	}

	token, err := handler.issuer.Issue(&ApplicationMiddleware.UserClaims{UserID: userID, Roles: []string{"USER"}})
	if err != nil {
		Http.Error(response, "Could not create application session", Http.StatusInternalServerError)
		return
	}

	WriteResponse(Http.StatusOK, response, request, loginResponse{
		Token:    token,
		UserID:   userID,
		Provider: identity.Provider,
	})
}
