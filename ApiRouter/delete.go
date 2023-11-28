package ApiRouter

import (
	"github.com/gin-gonic/gin"
)

func (r *ApiRouter[T]) Delete(c *gin.Context) {
	id := c.Param("id")
	_, err := r.Orm.GetByID(id)
	if err != nil {
		c.AbortWithStatusJSON(404, "Resource not found")
		return
	}
	err = r.Orm.Delete(id)
	if err != nil {
		c.AbortWithStatusJSON(500, "Database issue")
		return
	}

	c.JSON(204, nil)
}
