package middleware

import (
	Context "context"
	Http "net/http"

	Logger "github.com/danyel/ecommerce/cmd/logger"
	Uuid "github.com/google/uuid"
)

func CorrelationIDMiddleware(next Http.Handler) Http.Handler {
	return Http.HandlerFunc(func(response Http.ResponseWriter, request *Http.Request) {
		Logger.Log.Info("Entering CorrelationIDMiddleware")
		ID := request.Header.Get("X-Correlation-ID")
		if ID == "" {
			ID = request.Header.Get("X-Request-ID")
		}

		if ID == "" {
			ID = Uuid.New().String()
		}

		response.Header().Set("X-Correlation-ID", ID)

		ctx := Context.WithValue(request.Context(), Logger.CorrelationIDKey, ID)

		next.ServeHTTP(response, request.WithContext(ctx))
	})
}
