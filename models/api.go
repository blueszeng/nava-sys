package models

import "time"

type APIResponse struct {
	Status string
	Message string
	Result interface{}
}

type APISearch struct {
	Name string
}

// Base structure contains fields that are common to objects returned by the nava's REST
// API.
type Base struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt bool
}