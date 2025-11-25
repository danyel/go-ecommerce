package commonHandler

import (
	"encoding/json"
	"net/http"
)

type ResponseHandler interface {
	WriteResponse(w http.ResponseWriter, v any)
	StatusBadRequest(w http.ResponseWriter)
	StatusNotFound(w http.ResponseWriter)
	StatusInternalServerError(w http.ResponseWriter)
	StatusOK(w http.ResponseWriter)
}

type responseHandler struct {
}

func (r *responseHandler) WriteResponse(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
