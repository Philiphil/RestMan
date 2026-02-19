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
	obj, err := r.Orm.GetByID(c.Request.Context(), id)
	if err != nil {
		bfr := r.Orm.NewEntity()
		obj = &bfr
	}

	if err = r.WritingCheck(c, obj); err != nil {
		abortWithError(c, err)
		return
	}

	groups, errGroups := r.GetConfiguration(configuration.InputSerializationGroupsType, route.Put)
	if errGroups != nil {
		abortWithApiError(c, errors.Internal)
		return
	}

	if err = UnserializeBodyAndMerge(c, obj, groups.Values...); err != nil {
		abortWithError(c, err)
		return
	}

	var cast entity.Entity
	cast = *obj
	cast = cast.SetId(id)

	convertedEntity, ok := cast.(T)
	if !ok {
		abortWithApiError(c, errors.Internal)
		return
	}
	err = r.Orm.Update(c.Request.Context(), &convertedEntity)
	if err != nil {
		abortWithApiError(c, errors.DatabaseIssue)
		return
	}

	r.Render(c, 200, &convertedEntity, route.Put)

}
