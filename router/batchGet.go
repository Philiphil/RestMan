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

	formatedId := make([]entity.ID, len(idsValues))
	for i, v := range idsValues {
		formatedId[i] = entity.CastId(v)
	}
	objects, err := r.Orm.FindByIDs(formatedId)
	if err != nil {
		c.AbortWithStatusJSON(errors.ErrNotFound.Code, errors.ErrNotFound.Message)
		return
	}
	for _, object := range objects {
		err = r.ReadingCheck(c, object)
		if err != nil {
			c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
			return
		}
	}

	responseFormat, err := ParseAcceptHeader(c.GetHeader("Accept"))
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	groups, err := r.GetEffectiveOutputSerializationGroups(c, route.BatchGet)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	c.Render(200, SerializerRenderer{
		Data:   objects,
		Format: responseFormat,
		Groups: groups,
	})
}
