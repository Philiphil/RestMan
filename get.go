package apiman

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/errors"
	"github.com/philiphil/apiman/format"
	"github.com/philiphil/apiman/method/MethodType"
	"github.com/philiphil/apiman/router"
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

	c.Render(200, router.SerializerRenderer{
		Data:   object,
		Format: format.JSON,
		Groups: config.SerializationGroups,
	})
}
