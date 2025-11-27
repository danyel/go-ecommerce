package commonHandler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ResponseHandler interface {
	WriteResponse(status int, w http.ResponseWriter, v any)
	StatusBadRequest(w http.ResponseWriter)
	StatusNotFound(w http.ResponseWriter)
	StatusInternalServerError(w http.ResponseWriter)
	StatusOK(w http.ResponseWriter)
}

type responseHandler struct {
}

func GetId(r *http.Request, key string) (uuid.UUID, error) {
	productId := chi.URLParam(r, key)
	return uuid.Parse(productId)
}

func ValidateRequest[T any](req *http.Request, model *T) error {
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(model); err != nil {
		return err
	}

	return nil
}

func (r *responseHandler) WriteResponse(status int, w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		r.StatusInternalServerError(w)
	}
}

func (r *responseHandler) StatusOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (r *responseHandler) StatusBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func (r *responseHandler) StatusNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func (r *responseHandler) StatusInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func NewResponseHandler() ResponseHandler {
	return &responseHandler{}
}
