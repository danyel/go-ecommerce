package handler

import (
	JSON "encoding/json"
	Errors "errors"
	Fmt "fmt"
	Http "net/http"
	Strings "strings"

	Types "github.com/danyel/ecommerce/internal/types"
	Validator "github.com/go-playground/validator"

	Logger "github.com/danyel/ecommerce/cmd/logger"
	Domain "github.com/danyel/ecommerce/internal/domain"
	Router "github.com/go-chi/chi/v5"
	Uuid "github.com/google/uuid"
)

var validate = Validator.New()

type Title = string

const (
	NotFoundTitle            Title = "Not Found"
	BadRequestTitle          Title = "Not Found"
	IdNotFoundTitle          Title = "Id Not Found In Request"
	TodoTitle                Title = "TODO"
	InvalidRequestTitle      Title = "Invalid Request"
	ValidationTitle          Title = "Validation"
	InternalServerErrorTitle Title = "Internal Server Error"
	NotParseableTitle        Title = "Not Parseable"
)

type ProblemDetail = Domain.ProblemDetail

//func (p ProblemDetail) MarshalJSON() ([]byte, error) {
//	type Alias ProblemDetail
//	b, err := JSON.Marshal(Alias(p))
//	if err != nil {
//		return nil, err
//	}
//
//	if len(p.Errors) == 0 {
//		return b, nil
//	}
//
//	// Merge standard fields and custom extra fields
//	var res map[string]any
//	_ = JSON.Unmarshal(b, &res)
//	Maps.Copy(res, p.Errors)
//	return JSON.Marshal(res)
//}

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

//goland:noinspection GoUnusedExportedFunction
func GetRequestParam(request *Http.Request, key string) string {
	return request.URL.Query().Get(key)
}

func ValidateRequest[T any](request *Http.Request, model *T) (map[string]any, error) {
	decoder := JSON.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	// 1. Check for JSON syntax or unknown field errors
	if err := decoder.Decode(model); err != nil {
		return map[string]any{"json": err.Error()}, err
	}

	// 2. Validate struct fields against validation tags
	if err := validate.Struct(model); err != nil {
		var validationErrors Validator.ValidationErrors
		if Errors.As(err, &validationErrors) {
			errorDetails := make(map[string]any)

			for _, fieldErr := range validationErrors {
				// Convert structural field names into JSON property naming conventions (lowercase)
				fieldName := Strings.ToLower(fieldErr.Field())

				// Generate clean, readable error strings based on the tag triggered
				switch fieldErr.Tag() {
				case "required":
					errorDetails[fieldName] = "This field is required and cannot be empty."
				case "min":
					errorDetails[fieldName] = Fmt.Sprintf("Must be at least %s characters long.", fieldErr.Param())
				default:
					errorDetails[fieldName] = Fmt.Sprintf("Failed validation rule: %s", fieldErr.Tag())
				}
			}
			return errorDetails, err
		}
		return map[string]any{"validation": err.Error()}, err
	}

	return nil, nil
}

func WriteResponse(status int, response Http.ResponseWriter, request *Http.Request, body any) {
	setHeaders(response, request)
	response.WriteHeader(status)
	encoder := JSON.NewEncoder(response)
	if err := encoder.Encode(body); err != nil {
		errorDetails := make(map[string]any)
		errorDetails["error"] = err.Error()
		InternalServerError(response, request, NotParseableTitle, errorDetails)
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

func BadRequest(response Http.ResponseWriter, request *Http.Request, message Title, details map[string]any) {
	WriteJSONError(response, request, Http.StatusBadRequest, message, details)
}

func NotFound(response Http.ResponseWriter, request *Http.Request, message Title, details map[string]any) {
	WriteJSONError(response, request, Http.StatusNotFound, message, details)
}

func ProblemDetailResponse(response Http.ResponseWriter, request *Http.Request, problemDetail ProblemDetail) {
	WriteJSONError(response, request, problemDetail.Status, problemDetail.Title, problemDetail.Errors)
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

func InternalServerError(response Http.ResponseWriter, request *Http.Request, message Title, details map[string]any) {
	WriteJSONError(response, request, Http.StatusInternalServerError, message, details)
}

func WriteJSONError(response Http.ResponseWriter, request *Http.Request, statusCode int, message Title, details map[string]any) {
	setHeaders(response, request)
	response.Header().Set("Content-Type", "application/problem+json")
	response.WriteHeader(statusCode)
	path := request.URL.Path
	resp := ProblemDetail{
		Status:   statusCode,
		Instance: path,
		Title:    message,
		Details:  request.Header.Get("X-Correlation-Id"),
		Errors:   details,
	}

	err := JSON.NewEncoder(response).Encode(resp)
	if err != nil {
		Logger.Log.Fatalf("Error writing response: %v", err)
		return
	}
}

func StatusCreated(response Http.ResponseWriter, request *Http.Request, body any) {
	WriteResponse(Http.StatusCreated, response, request, body)
}
