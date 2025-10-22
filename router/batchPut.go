package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/route"
)

// BatchPut handles PUT requests for multiple entities, fully replacing existing entities.
func (r *ApiRouter[T]) BatchPut(c *gin.Context) {
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
	//I must check the id's first
	//All of them must have an id and if it's already in use, I must check write permissions

	var ids []entity.ID
	var preexistingEntities []*T
	for _, e := range entities {
		if (*e).GetId() != entity.NullId {
			ids = append(ids, (*e).GetId())
		} else {
			//null id is not allowed
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
		}
	}
	//check if preexisting entities are writable
	if len(preexistingEntities) > 0 {
		for _, e := range preexistingEntities {
			if err := r.WritingCheck(c, e); err != nil {
				c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
				return
			}
		}
	}

	if err := r.Orm.Update(entities...); err != nil {
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
