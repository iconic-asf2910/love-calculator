package main

// CalculateRequest represents the incoming love calculation request.
type CalculateRequest struct {
	Name1 string `json:"name1"`
	Name2 string `json:"name2"`
}

// CalculateResponse represents the love calculation result.
type CalculateResponse struct {
	Percentage int    `json:"percentage"`
	Message    string `json:"message"`
}

// ErrorResponse represents an error returned by the API.
type ErrorResponse struct {
	Error string `json:"error"`
}
