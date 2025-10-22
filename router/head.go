package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/route"
	"github.com/philiphil/restman/serializer"
)

// Head handles HTTP HEAD requests to retrieve entity metadata without the response body.
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

	groups, err := r.GetConfiguration(configuration.SerializationGroupsType, route.Get)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}
	str, err := s.Serialize(object, groups.Values...)
	if err != nil {
		c.AbortWithStatusJSON(errors.ErrInternal.Code, errors.ErrInternal.Message)
		return
	}
	c.Header("Content-Type", string(responseFormat))
	c.Header("Content-Length", fmt.Sprint(len(str)))
	c.Status(200)
}
