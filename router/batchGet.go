package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/route"
)

// BatchGet handles GET requests for multiple entities by their IDs.
func (r *ApiRouter[T]) BatchGet(c *gin.Context) {
	idsValues := r.GetIds(c)

	batchLimit, err := r.GetBatchLimit(c, route.BatchGet)
	if err != nil {
		abortWithApiError(c, errors.Internal)
		return
	}
	if len(idsValues) == 0 || len(idsValues) > batchLimit {
		abortWithApiError(c, errors.BadRequest)
		return
	}

	formatedId := make([]entity.ID, len(idsValues))
	for i, v := range idsValues {
		formatedId[i] = entity.CastId(v)
	}
	objects, err := r.Orm.FindByIDs(c.Request.Context(), formatedId)
	if err != nil {
		abortWithApiError(c, errors.NotFound)
		return
	}
	for _, object := range objects {
		err = r.ReadingCheck(c, object)
		if err != nil {
			abortWithError(c, err)
			return
		}
	}

	r.Render(c, 200, objects, route.BatchGet)
}
