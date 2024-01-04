package apiman

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/apiman/method/MethodType"
	"github.com/philiphil/apiman/router"
	"github.com/philiphil/apiman/serializer/format"
)

func (r *ApiRouter[T]) Get(c *gin.Context) {
	object, err := r.Orm.GetByID(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(ErrNotFound.Code, ErrNotFound.Message)
		return
	}
	config := r.GetMethodConfiguration(method_type.Get)

	err = r.ReadingCheck(c, object)
	if err != nil {
		c.AbortWithStatusJSON(err.(ApiError).Code, err.(ApiError).Message)
		return
	}

	c.Render(200, router.SerializerRenderer{
		Data:   object,
		Format: format.JSON,
		Groups: config.SerializationGroups,
	})
}
