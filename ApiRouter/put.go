package ApiRouter

import (
	"ApiMan/Gin"
	"ApiMan/Gorm"
	"github.com/gin-gonic/gin"
)

func (r *ApiRouter[T]) Put(c *gin.Context) {
	id := c.Param("id")
	entity, err := r.Orm.GetByID(id)
	if err != nil {
		bfr := r.Orm.NewEntity()
		entity = &bfr
	}
	//fait peter le deserializer
	if !Gin.UnserializeBodyAndMerge(c, entity) {
		return
	}

	var cast Gorm.IEntity
	cast = *entity
	cast = cast.SetId(id)

	convertedEntity, _ := cast.(T)
	err = r.Orm.Update(&convertedEntity)
	if err != nil {
		c.AbortWithStatusJSON(500, "Database issue")
		return
	}

	c.JSON(204, nil)
}
