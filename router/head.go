package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/route"
	"github.com/philiphil/restman/serializer"
)

// Head handles HTTP HEAD requests to retrieve entity metadata without the response body.
func (r *ApiRouter[T]) Head(c *gin.Context) {
	object, err := r.Orm.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		abortWithApiError(c, errors.NotFound)
		return
	}
	responseFormat, err := ParseAcceptHeader(c.GetHeader("Accept"))
	if err != nil {
		abortWithError(c, err)
		return
	}
	s := serializer.NewSerializer(responseFormat)

	groups, err := r.GetEffectiveOutputSerializationGroups(c, route.Get)
	if err != nil {
		abortWithError(c, err)
		return
	}
	str, err := s.Serialize(object, serializer.SerializeOptions{
		Groups: groups,
	})
	if err != nil {
		abortWithApiError(c, errors.Internal)
		return
	}
	c.Header("Content-Type", string(responseFormat))
	c.Header("Content-Length", fmt.Sprint(len(str)))
	c.Status(200)
}
