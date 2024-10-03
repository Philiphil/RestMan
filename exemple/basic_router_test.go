package exemple_test

import (
	"fmt"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman"
	"github.com/philiphil/apiman/method"
	"github.com/philiphil/apiman/orm"
	"github.com/philiphil/apiman/orm/entity"
	"github.com/philiphil/apiman/orm/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Test struct {
	entity.Entity
}

func (e Test) GetId() entity.ID {
	return e.Id
}

func (e Test) SetId(id any) entity.IEntity {
	e.Id = entity.CastId(id)
	return e
}
func (t Test) ToEntity() Test {
	return t
}

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
	test_ := apiman.NewApiRouter[Test](
		*orm.NewORM[Test](repository.NewRepository[Test, Test](getDB())),
		method.DefaultApiMethods(),
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
		fmt.Println(w.Body.String())
		panic("Failed to start server")
	}

}
