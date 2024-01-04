package apiman

import (
	"github.com/gin-gonic/gin"
)

func (r *ApiRouter[T]) Delete(c *gin.Context) {
	id := c.Param("id")
	object, err := r.Orm.GetByID(id)
	if err != nil {
		c.AbortWithStatusJSON(ErrNotFound.Code, ErrNotFound.Message)
		return
	}
	if err = r.WritingCheck(c, object); err != nil {
		c.AbortWithStatusJSON(err.(ApiError).Code, err.(ApiError).Message)
		return
	}
	if err = r.Orm.Delete(id); err != nil {
		c.AbortWithStatusJSON(500, "Database issue")
		return
	}
	c.JSON(204, nil)
}
