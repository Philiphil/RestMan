package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/orm/entity"
)

func (r *ApiRouter[T]) batchDelete(c *gin.Context) {
	ids := r.GetIds(c)
	formatedId := make([]entity.ID, len(ids))
	for i, v := range ids {
		formatedId[i] = entity.CastId(v)
	}
	objects, err := r.Orm.FindByIDs(formatedId)
	if err != nil {
		c.AbortWithStatusJSON(errors.ErrNotFound.Code, errors.ErrNotFound.Message)
		return
	}
	for _, object := range objects {
		if err = r.WritingCheck(c, object); err != nil {
			c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
			return
		}
	}
	err = r.Orm.Delete(objects...)
	if err != nil {
		c.AbortWithStatusJSON(500, "Database issue")
		return
	}

	c.JSON(204, nil)
}
