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
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	groups, err := r.GetConfiguration(configuration.InputSerializationGroupsType, route.Post)
	if err != nil {
		c.AbortWithStatusJSON(errors.ErrInternal.Code, errors.ErrInternal.Message)
		return
	}

	if err := UnserializeBodyAndMerge(c, &entity, groups.Values...); err != nil {
		//if unserializable, might be array
		if _, ok := r.Routes[route.BatchPost]; ok {
			if err := UnserializeBodyAndMerge_A(c, &entities, groups.Values...); err != nil {
				//its still unserializable as an array
				c.AbortWithStatusJSON(errors.ErrBadFormat.Code, errors.ErrBadFormat.Message)
				return
			} else {
				single = false
			}
		} else {
			//batch is not allowed, so array or not it does not mater
			c.AbortWithStatusJSON(errors.ErrBadFormat.Code, errors.ErrBadFormat.Message)
			return
		}
	} else {
		entities = append(entities, &entity)
	}

	if err := r.Orm.Create(entities...); err != nil {
		c.AbortWithStatusJSON(errors.ErrDatabaseIssue.Code, errors.ErrDatabaseIssue.Message)
		return
	}
	responseFormat, errParse := ParseAcceptHeader(c.GetHeader("Accept"))
	if errParse != nil {
		c.AbortWithStatusJSON(errParse.(errors.ApiError).Code, errParse.(errors.ApiError).Message)
		return
	}

	//what is sent back should use the "get" serialization groups
	outputGroups, err := r.GetEffectiveOutputSerializationGroups(c, route.Get)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	if single {
		c.Render(201, SerializerRenderer{
			Data:   &entity,
			Format: responseFormat,
			Groups: outputGroups,
		})
	} else {
		//batch
		c.Render(201, SerializerRenderer{
			Data:   &entities,
			Format: responseFormat,
			Groups: outputGroups,
		})
	}

}
