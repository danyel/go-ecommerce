package commonandler

import (
	JSON "encoding/json"
	Http "net/http"

	Logger "github.com/danyel/ecommerce/cmd/logger"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Router "github.com/go-chi/chi/v5"
	Uuid "github.com/google/uuid"
)

type APIError struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"` // For validation-specific details (maps, slices)
}

func GetID(request *Http.Request) (Types.ID, error) {
	ID := Router.URLParam(request, "ID")
	var err error
	var newID Uuid.UUID
	if newID, err = Uuid.Parse(ID); err != nil {
		return Types.ID{}, err
	}
	return Types.ID{ID: newID}, nil
}

func GetHeader(request *Http.Request, key string) string {
	return request.Header.Get(key)
}

func GetPathParam(request *Http.Request, key string) string {
	return Router.URLParam(request, key)
}

func GetRequestParam(request *Http.Request, key string) string {
	return request.URL.Query().Get(key)
}

func ValidateRequest[T any](request *Http.Request, model *T) error {
	decoder := JSON.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(model); err != nil {
		return err
	}

	return nil
}

func WriteResponse(status int, response Http.ResponseWriter, request *Http.Request, body any) {
	setHeaders(response, request)
	response.WriteHeader(status)
	encoder := JSON.NewEncoder(response)
	if err := encoder.Encode(body); err != nil {
		StatusInternalServerError(response, request)
	}
}

//goland:noinspection GoUnusedExportedFunction
func StatusOK(response Http.ResponseWriter, request *Http.Request) {
	setHeaders(response, request)
	response.WriteHeader(Http.StatusOK)
}

func StatusNoContent(response Http.ResponseWriter, request *Http.Request) {
	setHeaders(response, request)
	response.WriteHeader(Http.StatusNoContent)
}

// Deprecated: use BadRequest(response Http.response, request *Http.Request, message string, details any)
func StatusBadRequest(response Http.ResponseWriter, request *Http.Request) {
	setHeaders(response, request)
	response.WriteHeader(Http.StatusBadRequest)
}

func BadRequest(response Http.ResponseWriter, request *Http.Request, message string, details any) {
	WriteJSONError(response, request, Http.StatusBadRequest, message, details)
}

func setHeaders(response Http.ResponseWriter, request *Http.Request) {
	header := GetHeader(request, "X-Correlation-Id")
	if header == "" {
		header = Uuid.NewString()
	}
	response.Header().Set("X-Correlation-ID", header)
	response.Header().Set("Content-Type", "application/json")
}

func StatusNotFound(response Http.ResponseWriter, request *Http.Request) {
	setHeaders(response, request)
	response.WriteHeader(Http.StatusNotFound)
}

// StatusInternalServerError deprecated
// Deprecated: use InternalServerError(response Http.ResponseWriter, request *Http.Request, message string, details any)
func StatusInternalServerError(response Http.ResponseWriter, request *Http.Request) {
	setHeaders(response, request)
	response.WriteHeader(Http.StatusInternalServerError)
}

func InternalServerError(response Http.ResponseWriter, request *Http.Request, message string, details any) {
	WriteJSONError(response, request, Http.StatusInternalServerError, message, details)
}

func WriteJSONError(response Http.ResponseWriter, request *Http.Request, statusCode int, message string, details any) {
	setHeaders(response, request)
	response.WriteHeader(statusCode)

	resp := APIError{
		Success: false,
		Status:  statusCode,
		Message: message,
		Errors:  details,
	}

	err := JSON.NewEncoder(response).Encode(resp)
	if err != nil {
		Logger.Log.Fatalf("Error writing response: %v", err)
		return
	}
}
