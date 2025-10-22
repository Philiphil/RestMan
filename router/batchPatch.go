package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/route"
)

// BatchPatch handles PATCH requests for multiple entities, partially updating existing entities.
func (r *ApiRouter[T]) BatchPatch(c *gin.Context) {
	var entities []*T
	if err := UnserializeBodyAndMerge_A(c, &entities); err != nil {
		//unserializable
		c.AbortWithStatusJSON(errors.ErrBadFormat.Code, errors.ErrBadFormat.Message)
		return
	} else if len(entities) == 0 {
		//empty body
		c.AbortWithStatusJSON(errors.ErrBadFormat.Code, errors.ErrBadFormat.Message)
		return
	}
	var ids []entity.ID
	var preexistingEntities []*T
	for _, e := range entities {
		if (*e).GetId() != entity.NullId {
			ids = append(ids, (*e).GetId())
		} else {
			//null id is not
			c.AbortWithStatusJSON(errors.ErrBadFormat.Code, errors.ErrBadFormat.Message)
			return
		}
	}
	//try a batch get
	preexistingEntities, err := r.Orm.FindByIDs(ids)
	if err != nil {
		//check only for database issue, non existing entities are not a problem
		if err != errors.NotAllItemFound {
			c.AbortWithStatusJSON(errors.ErrDatabaseIssue.Code, errors.ErrDatabaseIssue.Message)
			return
		} else {
			//not all items found
			c.AbortWithStatusJSON(errors.ErrNotFound.Code, errors.ErrNotFound.Message)
		}
	}
	//check if preexisting entities are writable
	if len(preexistingEntities) > 0 || len(preexistingEntities) != len(entities) {
		for _, e := range preexistingEntities {
			if err := r.WritingCheck(c, e); err != nil {
				c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
				return
			}
		}
	} else {
		//no preexisting entities
		c.AbortWithStatusJSON(errors.ErrNotFound.Code, errors.ErrNotFound.Message)
		return
	}

	if err := UnserializeBodyAndMerge_A(c, &preexistingEntities); err != nil {
		//unserializable
		c.AbortWithStatusJSON(errors.ErrBadFormat.Code, errors.ErrBadFormat.Message)
		return
	}

	if err := r.Orm.Update(preexistingEntities...); err != nil {
		c.AbortWithStatusJSON(errors.ErrDatabaseIssue.Code, errors.ErrDatabaseIssue.Message)
		return
	}
	responseFormat, err := ParseAcceptHeader(c.GetHeader("Accept"))
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	groups, err := r.GetConfiguration(configuration.SerializationGroupsType, route.Post)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	c.Render(204, SerializerRenderer{
		Data:   &entities,
		Format: responseFormat,
		Groups: groups.Values, //what is sent shall be compliant to get
	})
}
