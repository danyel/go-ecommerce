package middleware

import (
	Context "context"
	Errors "errors"
	Fmt "fmt"
	Http "net/http"
	Strings "strings"

	JWT "github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user_claims"

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
		return Http.HandlerFunc(func(w Http.ResponseWriter, r *Http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !Strings.HasPrefix(authHeader, "Bearer ") {
				Http.Error(w, "Unauthorized: Missing or malformed token", Http.StatusUnauthorized)
				return
			}
			tokenString := Strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := validateJWT(tokenString, secretKey)

			if err != nil {
				Http.Error(w, "Unauthorized: Invalid token", Http.StatusUnauthorized)
				return
			}

			ctx := Context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func validateJWT(tokenString string, secretKey string) (*UserClaims, error) {
	token, err := JWT.ParseWithClaims(tokenString, &CustomClaims{}, func(token *JWT.Token) (interface{}, error) {
		if _, ok := token.Method.(*JWT.SigningMethodHMAC); !ok {
			return nil, Fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		if Errors.Is(err, JWT.ErrTokenExpired) {
			return nil, Errors.New("token has expired")
		}
		return nil, Fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, Errors.New("invalid token claims")
	}

	return &UserClaims{
		UserID: claims.UserID,
		Roles:  claims.Roles,
	}, nil
}
