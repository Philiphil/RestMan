package router

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/errors"
	"github.com/philiphil/apiman/serializer"
	"github.com/philiphil/apiman/serializer/format"
	"io"
)

func UnserializeBody[T any](c *gin.Context, e *T) error {
	jsonData, err := io.ReadAll(c.Request.Body)
	bodyReader := bytes.NewReader(jsonData)
	c.Request.Body = io.NopCloser(bodyReader)
	if err != nil {
		return errors.ErrBadFormat
	}
	serializer_ := serializer.NewSerializer(format.JSON)
	err = serializer_.Deserialize(string(jsonData), e)
	if err != nil {
		return errors.ErrBadFormat
	}
	return nil
}

func UnserializeBodyAndMerge[T any](c *gin.Context, e *T) error {
	jsonData, err := io.ReadAll(c.Request.Body)
	bodyReader := bytes.NewReader(jsonData)
	c.Request.Body = io.NopCloser(bodyReader)
	if err != nil {
		return errors.ErrBadFormat
	}
	serializer_ := serializer.NewSerializer(format.JSON)
	err = serializer_.DeserializeAndMerge(string(jsonData), e)
	if err != nil {
		return errors.ErrBadFormat
	}
	return nil
}
