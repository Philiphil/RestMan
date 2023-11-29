package ApiMan

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/gorm"
	"github.com/philiphil/apiman/router"
)

func (r *ApiRouter[T]) Put(c *gin.Context) {
	id := c.Param("id")
	entity, err := r.Orm.GetByID(id)
	if err != nil {
		bfr := r.Orm.NewEntity()
		entity = &bfr
	}

	if !router.UnserializeBodyAndMerge(c, entity) {
		return
	}

	var cast gorm.IEntity
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
