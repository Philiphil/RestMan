package apiman

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/router"
	"github.com/philiphil/apiman/serializer/format"
)

func (r *ApiRouter[T]) Post(c *gin.Context) {
	entity := r.Orm.NewEntity()

	if !router.UnserializeBodyAndMerge(c, &entity) {
		return
	}

	err := r.Orm.Create(&entity)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "Database issue"})
		return
	}

	c.Render(201, router.SerializerRenderer{
		Data:   &entity,
		Format: format.JSON,
		Groups: []string{},
	})
}
