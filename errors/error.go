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
	ErrNotFound      = ApiError{http.StatusNotFound, "not found", true}
	ErrUnauthorized  = ApiError{http.StatusUnauthorized, "unauthorized", true}
	ErrBadSchema     = ApiError{http.StatusBadRequest, "bad schema", true}
	ErrNotAcceptable = ApiError{http.StatusNotAcceptable, "not acceptable", true}
	ErrBadFormat     = ApiError{http.StatusBadRequest, "could not parse format", true}
	ErrDatabaseIssue = ApiError{http.StatusInternalServerError, "database issue", true}
	ErrUnsupported   = ApiError{http.StatusTeapot, "unsupported", false}

	ErrBadMethod  = ApiError{http.StatusMethodNotAllowed, "method not allowed", true}
	ErrBadRequest = ApiError{http.StatusBadRequest, "bad request", true}
	ErrForbidden  = ApiError{http.StatusForbidden, "forbidden", true}
	ErrConflict   = ApiError{http.StatusConflict, "conflict", true}
	ErrInternal   = ApiError{http.StatusInternalServerError, "internal error", true}
)
