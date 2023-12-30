package model

type GenericResponse[T any] struct {
	Code   int    `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
	Data   T      `json:"data,omitempty"`
}
