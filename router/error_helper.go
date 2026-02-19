package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
)

// abortWithError aborts the request with the appropriate error code and message.
// If the error is an ApiError, it uses the error's code and message.
// Otherwise, it returns a generic 500 Internal Server Error.
func abortWithError(c *gin.Context, err error) {
	if apiErr, ok := err.(errors.ApiError); ok {
		c.AbortWithStatusJSON(apiErr.Code, apiErr.Message)
	} else {
		c.AbortWithStatusJSON(errors.Internal.Code, errors.Internal.Message)
	}
}

// abortWithApiError aborts the request with the given ApiError.
// This is a convenience function for predefined errors like errors.NotFound.
func abortWithApiError(c *gin.Context, apiErr errors.ApiError) {
	c.AbortWithStatusJSON(apiErr.Code, apiErr.Message)
}
