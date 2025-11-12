package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/route"
)

// Get handles HTTP GET requests to retrieve a single entity by ID.
func (r *ApiRouter[T]) Get(c *gin.Context) {
	object, err := r.Orm.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		abortWithApiError(c, errors.ErrNotFound)
		return
	}

	err = r.ReadingCheck(c, object)
	if err != nil {
		abortWithError(c, err)
		return
	}

	r.Render(c, 200, object, route.Get)

}
