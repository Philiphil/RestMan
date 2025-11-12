package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
)

// Delete handles HTTP DELETE requests to remove a single entity by ID.
func (r *ApiRouter[T]) Delete(c *gin.Context) {
	id := c.Param("id")
	object, err := r.Orm.GetByID(c.Request.Context(), id)
	if err != nil {
		abortWithApiError(c, errors.ErrNotFound)
		return
	}
	if err = r.WritingCheck(c, object); err != nil {
		abortWithError(c, err)
		return
	}
	if err = r.Orm.Delete(c.Request.Context(), object); err != nil {
		abortWithApiError(c, errors.ErrDatabaseIssue)
		return
	}
	c.Status(204)
}
