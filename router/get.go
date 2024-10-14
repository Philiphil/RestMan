package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	method_type "github.com/philiphil/restman/method/MethodType"
)

func (r *ApiRouter[T]) Get(c *gin.Context) {
	object, err := r.Orm.GetByID(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(errors.ErrNotFound.Code, errors.ErrNotFound.Message)
		return
	}
	config := r.GetMethodConfiguration(method_type.Get)

	err = r.ReadingCheck(c, object)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	responseFormat, err := ParseAcceptHeader(c.GetHeader("Accept"))
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	c.Render(200, SerializerRenderer{
		Data:   object,
		Format: responseFormat,
		Groups: config.SerializationGroups,
	})
}
