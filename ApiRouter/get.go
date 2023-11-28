package ApiRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/Gin"
	"github.com/philiphil/apiman/Serializer/Format"
)

func (r *ApiRouter[T]) Get(c *gin.Context) {
	object, err := r.Orm.GetByID(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(404, "Resource not found")
		return
	}

	c.Render(200, Gin.SerializerRenderer{
		Data:   object,
		Format: Format.JSON,
		Groups: []string{},
	})
}
