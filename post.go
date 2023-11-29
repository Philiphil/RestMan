package ApiMan

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/serializer/format"
)

func (r *ApiRouter[T]) Post(c *gin.Context) {
	entity := r.Orm.NewEntity()

	if !gin.UnserializeBodyAndMerge(c, &entity) {
		return
	}

	err := r.Orm.Create(&entity)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "Database issue"})
		return
	}

	c.Render(201, gin.SerializerRenderer{
		Data:   &entity,
		format: format.JSON,
		Groups: []string{},
	})
}
