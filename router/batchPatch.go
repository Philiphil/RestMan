package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/route"
)

// BatchPatch handles PATCH requests for multiple entities, partially updating existing entities.
func (r *ApiRouter[T]) BatchPatch(c *gin.Context) {
	var entities []*T

	batchLimit, err := r.GetBatchLimit(c, route.BatchPatch)
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
		} else {
			//not all items found
			abortWithApiError(c, errors.NotFound)
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
	} else {
		//no preexisting entities
		abortWithApiError(c, errors.NotFound)
		return
	}

	if err := UnserializeBodyAndMerge_A(c, &preexistingEntities); err != nil {
		//unserializable
		abortWithApiError(c, errors.BadFormat)
		return
	}

	if err := r.Orm.Update(c.Request.Context(), preexistingEntities...); err != nil {
		abortWithApiError(c, errors.DatabaseIssue)
		return
	}

	r.Render(c, 200, entities, route.BatchPatch)
}
