package ApiRouter

import (
	"ApiMan/Gin"
	"ApiMan/Serializer/Format"
	"github.com/gin-gonic/gin"
)

func (r *ApiRouter[T]) Post(c *gin.Context) {
	entity := r.Orm.NewEntity()

	if !Gin.UnserializeBodyAndMerge(c, &entity) {
		return
	}

	err := r.Orm.Create(&entity)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "Database issue"})
		return
	}

	c.Render(201, Gin.SerializerRenderer{
		Data:   &entity,
		Format: Format.JSON,
		Groups: []string{},
	})
}
