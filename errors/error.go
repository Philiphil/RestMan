package errors

import "net/http"

type ApiError struct {
	Code     int
	Message  string
	Blocking bool
}

func (f ApiError) Error() string {
	return f.Message
}

var (
	ErrNotFound      = ApiError{http.StatusNotFound, "Not found", true}
	ErrUnauthorized  = ApiError{http.StatusUnauthorized, "Unauthorized", true}
	ErrBadSchema     = ApiError{http.StatusBadRequest, "Bad schema", true}
	ErrBadFormat     = ApiError{http.StatusBadRequest, "Could not parse format", true}
	ErrDatabaseIssue = ApiError{http.StatusInternalServerError, "Database issue", true}
	ErrUnsupported   = ApiError{http.StatusTeapot, "Unsupported", false}
)
