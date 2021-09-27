package models

type ErrorResponse struct {
	Success bool
	Error ErrorEntry
}

type ErrorEntry struct {
	Code string
}