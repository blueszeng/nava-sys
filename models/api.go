package models

type APIResponse struct {
	Status string
	Message string
	Result interface{}
}

type APISearch struct {
	Name string
}