package security_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository"
	"github.com/philiphil/restman/route"
	"github.com/philiphil/restman/router"
	. "github.com/philiphil/restman/security"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

func TestAuthenticationRequired(t *testing.T) {
	if !AuthenticationRequired(nil, nil) {
		t.Error("AuthenticationRequired() should return true")
	}
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

// ProtectedEntity is a test entity that is protected by the firewall
type ProtectedEntity struct {
	entity.BaseEntity
}

func (e ProtectedEntity) GetId() entity.ID {
	return e.Id
}

func (e ProtectedEntity) SetId(id any) entity.Entity {
	e.Id = entity.CastId(id)
	return e
}
func (t ProtectedEntity) ToEntity() ProtectedEntity {
	return t
}

func (t ProtectedEntity) FromEntity(entity ProtectedEntity) any {
	return entity
}
func (e ProtectedEntity) GetReadingRights() AuthorizationFunction {
	return AuthenticationRequired
}

func SetupFireWallRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	test_ := router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[ProtectedEntity](getDB())),
		route.DefaultApiRoutes(),
	)
	test_.AddFirewall(TestFirewall{})
	getDB().AutoMigrate(&ProtectedEntity{})
	test_.AllowRoutes(r)
	return r
}

type TestUser struct {
	entity.BaseEntity
}

func (e TestUser) GetId() entity.ID {
	return e.Id
}
func (e TestUser) SetId(id any) entity.Entity {
	e.Id = entity.CastId(id)
	return e
}

// This is a test firewall that only checks if the token is a valid id
type TestFirewall struct {
}

func (f TestFirewall) GetUser(c *gin.Context) (User, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return nil, errors.ErrUnauthorized
	}
	id := entity.CastId(token)
	if id == 0 { // non blocking error
		return TestUser{}.SetId(id), errors.ApiError{Code: http.StatusUnauthorized, Message: "unauthorized", Blocking: false}
	}
	return TestUser{}.SetId(id), nil
}

func TestMain_FireWall(t *testing.T) {
	r := SetupFireWallRouter()
	w := httptest.NewRecorder()

	repo := orm.NewORM(gormrepository.NewRepository[ProtectedEntity](getDB()))
	repo.Create(&ProtectedEntity{})
	w = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/protected_entity/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Error("Should return 401")
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/protected_entity/1", nil)
	req.Header.Add("Authorization", "1")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("Should return 200")
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/protected_entity/1", nil)
	req.Header.Add("Authorization", "0")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("Should return 200")
	}

}
