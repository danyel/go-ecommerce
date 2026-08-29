package service

import (
	Fmt "fmt"
)

type ServiceError struct {
	Status  Status
	Fields  map[string]string
	Message string
}

type Status string

func (e *ServiceError) Error() string {
	return Fmt.Sprintf("Error %s with fields: %v", e.Message, e.Fields)
}
