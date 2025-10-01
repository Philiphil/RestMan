package example_test

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
	"github.com/philiphil/restman/security"
)

// protected entity is an example of entity that requires authentication
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

// Reading and Writing rights are the interfaces used
// they require GetReadingRights and GetWritingRights functions
// here we use a default implementation that requires authentication and nothing else
// in order to demonstrate the use of the firewall
func (e ProtectedEntity) GetReadingRights() security.AuthorizationFunction {
	return security.AuthenticationRequired
}

// So we need to define a firewall, which is a struct that implements the security.Firewall interface
// the firewall will be used to get the user from the request using Headers, Cookies, etc
type TestFirewall struct {
}

// The GetUser function is used to get the user from the request
// it returns the user and an error
// Here we use the Authorization header to get the user id, not to be used in production obviously
func (f TestFirewall) GetUser(c *gin.Context) (security.User, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return nil, errors.ErrUnauthorized
	}
	id := entity.CastId(token)
	return TestUser{}.SetId(id), nil
}

// The user must implement the security.User interface which is essentialy the same as the entity.Entity interface for convenience
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
		t.Error("Should return 401")
	}

}
