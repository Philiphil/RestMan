package apiman

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/errors"
	"github.com/philiphil/apiman/router"
	"github.com/philiphil/apiman/serializer/format"
)

func (r *ApiRouter[T]) Post(c *gin.Context) {
	entity := r.Orm.NewEntity()
	if err := r.WritingCheck(c, &entity); err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	if err := router.UnserializeBodyAndMerge(c, &entity); err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	if err := r.Orm.Create(&entity); err != nil {
		c.AbortWithStatusJSON(errors.ErrDatabaseIssue.Code, errors.ErrDatabaseIssue.Message)
		return
	}

	c.Render(201, router.SerializerRenderer{
		Data:   &entity,
		Format: format.JSON,
		Groups: []string{},
	})
}
