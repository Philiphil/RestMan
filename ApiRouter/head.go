package ApiRouter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/Serializer"
	"github.com/philiphil/apiman/Serializer/Format"
)

func (r *ApiRouter[T]) Head(c *gin.Context) {
	object, err := r.Orm.GetByID(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(404, "Resource not found")
		return
	}

	s := serializer.NewSerializer(Format.JSON)

	str, err := s.Serialize(object, "")
	c.Header("Content-Type", "application/json")
	c.Header("Content-Length", fmt.Sprint(len(str)))
	c.Status(200)
}
