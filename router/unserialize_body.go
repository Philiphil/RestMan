package router

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/serializer"
)

func UnserializeBodyAndMerge[T any](c *gin.Context, e *T) error {
	serializedData, err := io.ReadAll(c.Request.Body)
	bodyReader := bytes.NewReader(serializedData)
	c.Request.Body = io.NopCloser(bodyReader)
	if err != nil {
		return errors.ErrBadFormat
	}
	serializer_ := serializer.NewSerializer(ParseTypeFromString(c.GetHeader("Content-type")))
	err = serializer_.DeserializeAndMerge(string(serializedData), e)
	if err != nil {
		return errors.ErrBadFormat
	}
	return nil
}

func UnserializeBodyAndMerge_A[T any](c *gin.Context, e *[]*T) error {
	serializedData, err := io.ReadAll(c.Request.Body)
	bodyReader := bytes.NewReader(serializedData)
	c.Request.Body = io.NopCloser(bodyReader)
	if err != nil {
		return errors.ErrBadFormat
	}
	serializer_ := serializer.NewSerializer(ParseTypeFromString(c.GetHeader("Content-type")))
	err = serializer_.DeserializeAndMerge(string(serializedData), &e)
	if err != nil {
		return errors.ErrBadFormat
	}
	return nil
}
