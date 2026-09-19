package handler

import (
	"encoding/json"
	"net/http"

	ApplicationMiddleware "github.com/danyel/ecommerce/cmd/middleware"
	"github.com/danyel/ecommerce/internal/security"
)

type AuthWebHandler struct {
	registry *security.ProviderRegistry
	linker   security.UserLinker
	issuer   security.SessionIssuer
}

type googleLoginRequest struct {
	IDToken string `json:"id_token"`
}

type loginResponse struct {
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	Provider string `json:"provider"`
}

func NewAuthWebHandler(registry *security.ProviderRegistry, linker security.UserLinker, issuer security.SessionIssuer) *AuthWebHandler {
	return &AuthWebHandler{registry: registry, linker: linker, issuer: issuer}
}

func (handler *AuthWebHandler) HandleGoogleLogin(response http.ResponseWriter, request *http.Request) {
	var login googleLoginRequest
	if err := json.NewDecoder(request.Body).Decode(&login); err != nil || login.IDToken == "" {
		http.Error(response, "Invalid Google login request", http.StatusBadRequest)
		return
	}
	provider, ok := handler.registry.Provider("google")
	if !ok {
		http.Error(response, "Google login is not configured", http.StatusServiceUnavailable)
		return
	}
	identity, err := provider.Verify(request.Context(), login.IDToken)
	if err != nil {
		http.Error(response, "Google identity could not be verified", http.StatusUnauthorized)
		return
	}
	userID, err := handler.linker.Link(request.Context(), identity)
	if err != nil {
		http.Error(response, "Could not link account", http.StatusInternalServerError)
		return
	}
	token, err := handler.issuer.Issue(&ApplicationMiddleware.UserClaims{UserID: userID, Roles: []string{"USER"}})
	if err != nil {
		http.Error(response, "Could not create application session", http.StatusInternalServerError)
		return
	}
	WriteResponse(http.StatusOK, response, request, loginResponse{Token: token, UserID: userID, Provider: identity.Provider})
}
