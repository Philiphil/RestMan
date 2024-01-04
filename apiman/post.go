package apiman

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/router"
	"github.com/philiphil/apiman/serializer/format"
)

func (r *ApiRouter[T]) Post(c *gin.Context) {
	entity := r.Orm.NewEntity()
	if err := r.WritingCheck(c, &entity); err != nil {
		c.AbortWithStatusJSON(err.(ApiError).Code, err.(ApiError).Message)
		return
	}

	if err := router.UnserializeBodyAndMerge(c, &entity); err != nil {
		c.AbortWithStatusJSON(err.(ApiError).Code, err.(ApiError).Message)
		return
	}

	if err := r.Orm.Create(&entity); err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "Database issue"})
		return
	}

	c.Render(201, router.SerializerRenderer{
		Data:   &entity,
		Format: format.JSON,
		Groups: []string{},
	})
}
