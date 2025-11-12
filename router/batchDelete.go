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
	objects, err := r.Orm.FindByIDs(c.Request.Context(), formatedId)
	if err != nil {
		abortWithApiError(c, errors.ErrNotFound)
		return
	}
	for _, object := range objects {
		if err = r.WritingCheck(c, object); err != nil {
			abortWithError(c, err)
			return
		}
	}
	err = r.Orm.Delete(c.Request.Context(), objects...)
	if err != nil {
		abortWithApiError(c, errors.ErrDatabaseIssue)
		return
	}

	c.Status(204)
}
