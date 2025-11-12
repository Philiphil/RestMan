package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/route"
)

// Post handles HTTP POST requests to create one or more new entities.
func (r *ApiRouter[T]) Post(c *gin.Context) {
	single := true
	var entities []*T
	entity := r.Orm.NewEntity()

	if err := r.WritingCheck(c, &entity); err != nil {
		abortWithError(c, err)
		return
	}

	groups, err := r.GetConfiguration(configuration.InputSerializationGroupsType, route.Post)
	if err != nil {
		abortWithApiError(c, errors.ErrInternal)
		return
	}

	if err := UnserializeBodyAndMerge(c, &entity, groups.Values...); err != nil {
		//if unserializable, might be array
		if _, ok := r.Routes[route.BatchPost]; ok {
			if err := UnserializeBodyAndMerge_A(c, &entities, groups.Values...); err != nil {
				//its still unserializable as an array
				abortWithApiError(c, errors.ErrBadFormat)
				return
			} else {
				single = false
			}
		} else {
			//batch is not allowed, so array or not it does not mater
			abortWithApiError(c, errors.ErrBadFormat)
			return
		}
	} else {
		entities = append(entities, &entity)
	}

	if err := r.Orm.Create(c.Request.Context(), entities...); err != nil {
		abortWithApiError(c, errors.ErrDatabaseIssue)
		return
	}

	if single {
		r.Render(c, 201, &entity, route.Post)
	} else {
		//batch
		r.Render(c, 201, &entities, route.BatchPost)
	}

}
