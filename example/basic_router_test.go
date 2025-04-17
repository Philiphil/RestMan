package example_test

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository"
	"github.com/philiphil/restman/route"
	"github.com/philiphil/restman/router"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

//This is the most basic example of how to use the restman package

// We have an entity, a struct that implements the entity.Entity interface
// here we chose to use the BaseEntity struct that implements the ID field and random common fields
type Test struct {
	entity.BaseEntity
}

// This function is required by the Entity interface
func (e Test) GetId() entity.ID {
	return e.Id
}

// and this one, note that we do not use pointer here
func (e Test) SetId(id any) entity.Entity {
	e.Id = entity.CastId(id)
	return e
}

// This one is for the Model Interface, since we are using the same struct for the Entity and the Model
// we can just return the struct itself
func (t Test) ToEntity() Test {
	return t
}

// same thing
func (t Test) FromEntity(entity Test) any {
	return entity
}

func getDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:test?mode=memory&cache=shared&_fk=1"), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Silent),
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic(err)
	}
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	})

	db.Session(&gorm.Session{FullSaveAssociations: true})
	return db
}

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	test_ := router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Test](getDB())),
		route.DefaultApiRoutes(),
	)

	getDB().AutoMigrate(&Test{})
	test_.AllowRoutes(r)
	return r
}

func TestMain(m *testing.M) {
	r := SetupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		panic("Failed to start server")
	}

}
