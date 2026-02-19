package errors

import (
	base "errors"
	"net/http"
)

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

// Error implements the error interface for ApiError.
func (f ApiError) Error() string {
	return f.Message
}

func New(text string) error {
	return base.New(text)
}

// This is a helper function to create a firewall blocking error with a custom message
// This is a spread parameter to make it optional, if the length of the message is greater than 0, it will use the first element as the message
func FirewallBlocking(message ...string) ApiError {
	e := Unauthorized
	if len(message) > 0 {
		e.Message = message[0]
	}
	return e
}

// This is a helper function to create a firewall blocking error with a custom message
// This is a spread parameter to make it optional, if the length of the message is greater than 0, it will use the first element as the message
func FirewallNonBlocking(message ...string) ApiError {
	e := Unauthorized
	if len(message) > 0 {
		e.Message = message[0]
	}
	e.Blocking = false
	return e
}

// Default HTTP errors ...
var (
	//400
	BadRequest = ApiError{http.StatusBadRequest, "bad request", true}
	BadFormat  = ApiError{http.StatusBadRequest, "could not parse format", true}
	//Unused ? needs to investigate, another error must have been used instead of this one
	BadSchema = ApiError{http.StatusBadRequest, "bad schema", true}

	//401
	Unauthorized = ApiError{http.StatusUnauthorized, "unauthorized", true}

	//403 Unused ? needs to investigate, another error must have been used instead of this one
	Forbidden = ApiError{http.StatusForbidden, "forbidden", true}

	//404 is there a more iconic one
	NotFound = ApiError{http.StatusNotFound, "not found", true}

	//405 Unused ? needs to investigate, another error must have been used instead of this one
	BadMethod = ApiError{http.StatusMethodNotAllowed, "method not allowed", true}

	// 406
	NotAcceptable = ApiError{http.StatusNotAcceptable, "not acceptable", true}

	//409 Unused ? needs to investigate, another error must have been used instead of this one
	Conflict = ApiError{http.StatusConflict, "conflict", true}

	// 500
	// db is cause we dont want client to know about database ...
	DatabaseIssue = ApiError{http.StatusInternalServerError, "database issue", true}

	//This one is the default I DONT WANNA TELL, which can be normal on a prod level but
	//I should investigate each case to see if we cannot add a more specific one
	//or adding a DISPLAY MORE PRECISE STUFF ON NON PRODUCTION ENVIRONMENTS mechanism
	Internal = ApiError{http.StatusInternalServerError, "internal error", true}
)
