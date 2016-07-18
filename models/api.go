package models

import "time"

type APIResponse struct {
	Status string
	Message string
	Result interface{}
}

// Structure for collection of search string for frontend request.
type APISearch struct {
	Name string
}

// Base structure contains fields that are common to objects
// returned by the nava's REST API.
type Base struct {
	ID        uint64 `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}