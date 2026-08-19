package commonandler

import (
	JSON "encoding/json"
	Http "net/http"

	Types "github.com/danyel/ecommerce/internal/common/types"
	Router "github.com/go-chi/chi/v5"
	Uuid "github.com/google/uuid"
)

func GetID(r *Http.Request, key string) (Types.ID, error) {
	productID := Router.URLParam(r, key)
	var err error
	var newID Uuid.UUID
	if newID, err = Uuid.Parse(productID); err != nil {
		return Types.ID{}, err
	}
	return Types.ID{ID: newID}, nil
}

func GetHeader(r *Http.Request, key string) string {
	return r.Header.Get(key)
}

func GetPathParam(r *Http.Request, key string) string {
	return Router.URLParam(r, key)
}

func GetRequestParam(r *Http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func ValidateRequest[T any](req *Http.Request, model *T) error {
	decoder := JSON.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(model); err != nil {
		return err
	}

	return nil
}

func WriteResponse(status int, w Http.ResponseWriter, r *Http.Request, v any) {
	setHeaders(w, r)
	w.WriteHeader(status)
	encoder := JSON.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		StatusInternalServerError(w, r)
	}
}

//goland:noinspection GoUnusedExportedFunction
func StatusOK(w Http.ResponseWriter, r *Http.Request) {
	setHeaders(w, r)
	w.WriteHeader(Http.StatusOK)
}

func StatusNoContent(w Http.ResponseWriter, r *Http.Request) {
	setHeaders(w, r)
	w.WriteHeader(Http.StatusNoContent)
}

func StatusBadRequest(w Http.ResponseWriter, r *Http.Request) {
	setHeaders(w, r)
	w.WriteHeader(Http.StatusBadRequest)
}

func setHeaders(w Http.ResponseWriter, r *Http.Request) {
	header := GetHeader(r, "X-Correlation-Id")
	if header == "" {
		header = Uuid.NewString()
	}
	w.Header().Set("X-Correlation-ID", header)
	w.Header().Set("Content-Type", "application/json")
}

func StatusNotFound(w Http.ResponseWriter, r *Http.Request) {
	setHeaders(w, r)
	w.WriteHeader(Http.StatusNotFound)
}

func StatusInternalServerError(w Http.ResponseWriter, r *Http.Request) {
	setHeaders(w, r)
	w.WriteHeader(Http.StatusInternalServerError)
}
