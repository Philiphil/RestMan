package example_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository"
	"github.com/philiphil/restman/route"
	"github.com/philiphil/restman/router"
)

//If needed we can separate an Entity from a Model
// Entity is the struct that will be used to interact with the database
// Model is the struct that will be used to interact with the API

// Entity needs to implement the entity.Entity interface
type Entity struct {
	ID   int
	Name string
}

func (e Entity) GetId() entity.ID {
	return entity.CastId(e.ID)
}

func (e Entity) SetId(id any) entity.Entity {
	e.ID = int(entity.CastId(id))
	return e
}

// Model must implement the Model interface
type Model struct {
	Entity
}

func (t Model) ToEntity() Entity {
	e := Entity{
		ID:   t.ID,
		Name: t.Name,
	}
	return e
}

func (t Model) FromEntity(entity Entity) any {
	t.ID = entity.ID
	t.Name = entity.Name
	return t
}

func testEntityModel(m *testing.M) {
	api := router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Model](getDB())),
		route.DefaultApiRoutes(),
	)

	r := gin.New()
	api.AllowRoutes(r)

}
