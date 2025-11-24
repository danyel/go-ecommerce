package request_util

import (
	"encoding/json"
	"net/http"
)

func ValidateRequest[T any](r *http.Request, model *T) error {
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		return err
	}

	return nil
}
