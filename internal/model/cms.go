package model

import Uuid "github.com/google/uuid"

type Translation struct {
	Code     string `json:"code"`
	Value    string `json:"value"`
	Language string `json:"language"`
}

type CreateCms struct {
	Code     string `json:"code"`
	Value    string `json:"value"`
	Language string `json:"language"`
}

type CmsID struct {
	ID Uuid.UUID `json:"id"`
}
