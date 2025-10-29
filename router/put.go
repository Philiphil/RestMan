package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/route"
)

// Put handles HTTP PUT requests to replace or create an entity at a specific ID.
func (r *ApiRouter[T]) Put(c *gin.Context) {
	id := c.Param("id")
	obj, err := r.Orm.GetByID(id)
	if err != nil {
		bfr := r.Orm.NewEntity()
		obj = &bfr
	}

	if err = r.WritingCheck(c, obj); err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	groups, errGroups := r.GetConfiguration(configuration.InputSerializationGroupsType, route.Put)
	if errGroups != nil {
		c.AbortWithStatusJSON(errors.ErrInternal.Code, errors.ErrInternal.Message)
		return
	}

	if err = UnserializeBodyAndMerge(c, obj, groups.Values...); err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	var cast entity.Entity
	cast = *obj
	cast = cast.SetId(id)

	convertedEntity, _ := cast.(T)
	err = r.Orm.Update(&convertedEntity)
	if err != nil {
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

	c.Render(200, SerializerRenderer{
		Data:   obj,
		Format: responseFormat,
		Groups: outputGroups,
	})
}
