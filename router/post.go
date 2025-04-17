package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/route"
)

func (r *ApiRouter[T]) Post(c *gin.Context) {
	single := true
	var entities []*T
	entity := r.Orm.NewEntity()

	if err := r.WritingCheck(c, &entity); err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	if err := UnserializeBodyAndMerge(c, &entity); err != nil {
		//if unserializable, might be array
		if _, ok := r.Routes[route.BatchPost]; ok {
			if err := UnserializeBodyAndMerge_A(c, &entities); err != nil {
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
	responseFormat, err := ParseAcceptHeader(c.GetHeader("Accept"))
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	groups, err := r.GetConfiguration(configuration.SerializationGroupsType, route.Post)
	if err != nil {
		c.AbortWithStatusJSON(errors.ErrInternal.Code, errors.ErrInternal.Message)
		return
	}

	if single {
		c.Render(201, SerializerRenderer{
			Data:   &entity,
			Format: responseFormat,
			Groups: groups.Values, //what is sent shall be compliant to get
		})
	} else {
		//batch
		c.Render(201, SerializerRenderer{
			Data:   &entities,
			Format: responseFormat,
			Groups: groups.Values, //what is sent shall be compliant to get
		})
	}

}
