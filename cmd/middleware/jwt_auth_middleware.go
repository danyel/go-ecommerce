package middleware

import (
	Context "context"
	Http "net/http"
	Array "slices"
	Strings "strings"

	JWT "github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserContextKey contextKey = "user_claims"
	Authorization  string     = "Authorization"
)

type UserClaims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
}

type CustomClaims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	JWT.RegisteredClaims
}

func JwtAuthMiddleware(secretKey string) func(handler Http.Handler) Http.Handler {
	return func(next Http.Handler) Http.Handler {
		return Http.HandlerFunc(func(response Http.ResponseWriter, request *Http.Request) {
			authHeader := request.Header.Get(Authorization)
			if authHeader == "" || !Strings.HasPrefix(authHeader, "Bearer ") {
				Http.Error(response, "Unauthorized: Missing or malformed token", Http.StatusUnauthorized)
				return
			}
			tokenString := Strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := DecryptClaims(tokenString, secretKey)
			if err != nil {
				Http.Error(response, "Unauthorized: Invalid token", Http.StatusUnauthorized)
				return
			}

			ctx := Context.WithValue(request.Context(), UserContextKey, claims)
			next.ServeHTTP(response, request.WithContext(ctx))
		})
	}
}

func RequireRole(requiredRole string) func(Http.Handler) Http.Handler {
	return func(next Http.Handler) Http.Handler {
		return Http.HandlerFunc(func(response Http.ResponseWriter, request *Http.Request) {
			claims, ok := request.Context().Value(UserContextKey).(*UserClaims)
			if !ok || !Array.Contains(claims.Roles, requiredRole) {
				Http.Error(response, "Forbidden: Insufficient permissions", Http.StatusForbidden)
				return
			}
			next.ServeHTTP(response, request)
		})
	}
}
