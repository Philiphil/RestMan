package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
)

func (r *ApiRouter[T]) Delete(c *gin.Context) {
	id := c.Param("id")
	object, err := r.Orm.GetByID(id)
	if err != nil {
		c.AbortWithStatusJSON(errors.ErrNotFound.Code, errors.ErrNotFound.Message)
		return
	}
	if err = r.WritingCheck(c, object); err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}
	if err = r.Orm.Delete(id); err != nil {
		c.AbortWithStatusJSON(500, "Database issue")
		return
	}
	c.JSON(204, nil)
}
