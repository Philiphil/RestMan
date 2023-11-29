package ApiMan

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/serializer"
	"github.com/philiphil/apiman/serializer/format"
)

func (r *ApiRouter[T]) Head(c *gin.Context) {
	object, err := r.Orm.GetByID(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(404, "Resource not found")
		return
	}

	s := serializer.NewSerializer(format.JSON)

	str, err := s.Serialize(object, "")
	c.Header("Content-Type", "application/json")
	c.Header("Content-Length", fmt.Sprint(len(str)))
	c.Status(200)
}
