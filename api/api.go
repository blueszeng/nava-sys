package api

type ResponseStatus int

const (
	SUCCESS ResponseStatus = iota
	FAIL
	ERROR
)

type Response struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Link    Link           `json:"links"`
}

type Link struct {
	Self    string `json:"self"`
	Related string `json:"relateed"`
	Next    string `json:"next"`
	Last    string `json:"last"`
}

// Structure for collection of search string for frontend request.
type Search struct {
	Name string
}
