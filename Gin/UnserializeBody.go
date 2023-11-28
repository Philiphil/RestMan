package Gin

import (
	"ApiMan/Serializer/Format"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
)

func UnserializeBody[T any](c *gin.Context, e *T) bool {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Could not parse JSON"})
		return false
	}
	serializer_ := serializer.NewSerializer(Format.JSON)
	err = serializer_.Deserialize(string(jsonData), e)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Bad schema"})
		return false
	}
	bodyReader := bytes.NewReader(jsonData)
	c.Request.Body = io.NopCloser(bodyReader)
	return true
}

func UnserializeBodyAndMerge[T any](c *gin.Context, e *T) bool {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Could not parse JSON"})
		return false
	}
	serializer_ := serializer.NewSerializer(Format.JSON)
	err = serializer_.DeserializeAndMerge(string(jsonData), e)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"message": "Bad schema"})
		return false
	}
	bodyReader := bytes.NewReader(jsonData)
	c.Request.Body = io.NopCloser(bodyReader)
	return true
}
