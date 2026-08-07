package commonHandler

import (
	JSON "encoding/json"
	Http "net/http"

	Types "github.com/danyel/ecommerce/internal/common/types"
	Router "github.com/go-chi/chi/v5"
	Uuid "github.com/google/uuid"
)

func GetId(r *Http.Request, key string) (Types.Id, error) {
	productId := Router.URLParam(r, key)
	var err error
	var newId Uuid.UUID
	if newId, err = Uuid.Parse(productId); err != nil {
		return Types.Id{}, err
	}
	return Types.Id{ID: newId}, nil
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

func WriteResponse(status int, w Http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := JSON.NewEncoder(w).Encode(v); err != nil {
		StatusInternalServerError(w)
	}
}

//goland:noinspection GoUnusedExportedFunction
func StatusOK(w Http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(Http.StatusOK)
}

func StatusNoContent(w Http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(Http.StatusNoContent)
}

func StatusBadRequest(w Http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(Http.StatusBadRequest)
}

func StatusNotFound(w Http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(Http.StatusNotFound)
}

func StatusInternalServerError(w Http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(Http.StatusInternalServerError)
}
