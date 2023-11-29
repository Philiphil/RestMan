package router

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/serializer"
	"github.com/philiphil/apiman/serializer/format"
	"io"
)

func UnserializeBody[T any](c *gin.Context, e *T) bool {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Could not parse JSON"})
		return false
	}
	serializer_ := serializer.NewSerializer(format.JSON)
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
	serializer_ := serializer.NewSerializer(format.JSON)
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
