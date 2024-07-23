package model

// ErrorResponse is returned when 400+ status
// codes are returned from API handlers
type ErrorResponse struct {
	Error string
}
