package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	method_type "github.com/philiphil/restman/method/MethodType"
	"github.com/philiphil/restman/serializer"
)

func (r *ApiRouter[T]) Head(c *gin.Context) {
	object, err := r.Orm.GetByID(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(errors.ErrNotFound.Code, errors.ErrNotFound.Message)
		return
	}
	responseFormat, err := ParseAcceptHeader(c.GetHeader("Accept"))
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}
	s := serializer.NewSerializer(responseFormat)

	str, err := s.Serialize(object, r.GetMethodConfiguration(method_type.Get).SerializationGroups...)
	if err != nil {
		c.AbortWithStatusJSON(errors.ErrInternal.Code, errors.ErrInternal.Message)
		return
	}
	c.Header("Content-Type", string(responseFormat))
	c.Header("Content-Length", fmt.Sprint(len(str)))
	c.Status(200)
}
