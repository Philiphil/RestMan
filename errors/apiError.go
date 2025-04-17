package errors

import "net/http"

// ApiError is a struct that represents an error in RestMan
// it has a code, a message and a boolean to indicate if the error is blocking
// blocking errors should stop the execution of the request
// An use case for a non-blocking error is when using multiple firewalls
// We  could have user-password firewall and a token firewall coexisting
type ApiError struct {
	Code     int
	Message  string
	Blocking bool
}

func (f ApiError) Error() string {
	return f.Message
}

// Default HTTP errors ...
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
