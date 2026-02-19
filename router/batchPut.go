package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/route"
)

// BatchPut handles PUT requests for multiple entities, fully replacing existing entities.
func (r *ApiRouter[T]) BatchPut(c *gin.Context) {
	var entities []*T

	batchLimit, err := r.GetBatchLimit(c, route.BatchPut)
	if err != nil {
		abortWithApiError(c, errors.Internal)
		return
	}

	if err := UnserializeBodyAndMerge_A(c, &entities); err != nil {
		//unserializable
		abortWithApiError(c, errors.BadFormat)
		return
	} else if len(entities) == 0 || len(entities) > batchLimit {
		abortWithApiError(c, errors.BadRequest)
		return
	}

	//I must check the id's first
	//All of them must have an id and if it's already in use, I must check write permissions

	ids := make([]entity.ID, len(entities))
	for i, e := range entities {
		id := (*e).GetId()
		if id == entity.NullId {
			abortWithApiError(c, errors.BadFormat)
			return
		}
		ids[i] = id
	}
	//try a batch get
	preexistingEntities, err := r.Orm.FindByIDs(c.Request.Context(), ids)
	if err != nil {
		//check only for database issue, non existing entities are not a problem
		if err != errors.NotAllItemFound {
			abortWithApiError(c, errors.DatabaseIssue)
			return
		}
	}
	//check if preexisting entities are writable
	if len(preexistingEntities) > 0 {
		for _, e := range preexistingEntities {
			if err := r.WritingCheck(c, e); err != nil {
				abortWithError(c, err)
				return
			}
		}
	}

	if err := r.Orm.Update(c.Request.Context(), entities...); err != nil {
		abortWithApiError(c, errors.DatabaseIssue)
		return
	}

	r.Render(c, 200, entities, route.BatchPut)

}
