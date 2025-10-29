package router

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/serializer"
)

// UnserializeBodyAndMerge deserializes the request body and merges it with the provided entity.
// If groups are provided, only fields with matching group tags will be deserialized.
func UnserializeBodyAndMerge[T any](c *gin.Context, e *T, groups ...string) error {
	serializedData, err := io.ReadAll(c.Request.Body)
	bodyReader := bytes.NewReader(serializedData)
	c.Request.Body = io.NopCloser(bodyReader)
	if err != nil {
		return errors.ErrBadFormat
	}
	serializer_ := serializer.NewSerializer(ParseTypeFromString(c.GetHeader("Content-type")))
	err = serializer_.DeserializeAndMerge(string(serializedData), e, groups...)
	if err != nil {
		return errors.ErrBadFormat
	}
	return nil
}

// UnserializeBodyAndMerge_A deserializes the request body as an array and merges it with the provided entity slice.
// If groups are provided, only fields with matching group tags will be deserialized.
func UnserializeBodyAndMerge_A[T any](c *gin.Context, e *[]*T, groups ...string) error {
	serializedData, err := io.ReadAll(c.Request.Body)
	bodyReader := bytes.NewReader(serializedData)
	c.Request.Body = io.NopCloser(bodyReader)
	if err != nil {
		return errors.ErrBadFormat
	}
	serializer_ := serializer.NewSerializer(ParseTypeFromString(c.GetHeader("Content-type")))
	err = serializer_.DeserializeAndMerge(string(serializedData), &e, groups...)
	if err != nil {
		return errors.ErrBadFormat
	}
	return nil
}
