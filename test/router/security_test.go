package router_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/format"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository"
	"github.com/philiphil/restman/route"
	"github.com/philiphil/restman/router"
	"github.com/philiphil/restman/security"
	"github.com/philiphil/restman/serializer"
)

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
func (e ProtectedEntity) GetReadingRights() security.AuthorizationFunction {
	return func(i1 security.User, i2 entity.Entity) bool {
		return i1.GetId() == 1
	}
}

func (e ProtectedEntity) GetWritingRights() security.AuthorizationFunction {
	return func(i1 security.User, i2 entity.Entity) bool {
		return i1.GetId() == 1
	}
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

func (f TestFirewall) GetUser(c *gin.Context) (security.User, error) {
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

	req, _ := http.NewRequest("GET", "/api/protected_entity", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		//get List security is not handled by default
		t.Error("Should return 200")
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/protected_entity", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Error("Should return 401")
	}
	//create a new entity and test to get it
	repo := orm.NewORM(gormrepository.NewRepository[ProtectedEntity](getDB()))
	repo.Create(&ProtectedEntity{})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/protected_entity/1", nil)
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
	if w.Code != http.StatusUnauthorized {
		t.Error("Should return 401")
	}

	w = httptest.NewRecorder()
	entity := ProtectedEntity{entity.BaseEntity{Id: 2, Name: "test"}}

	serializer := serializer.NewSerializer(format.JSON)

	jsonEntity, _ := serializer.Serialize(entity)
	req, _ = http.NewRequest("POST", "/api/protected_entity", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Authorization", "1")
	r.ServeHTTP(w, req)
	if w.Code == http.StatusUnauthorized {
		t.Error("Should not 401")
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/protected_entity", nil)
	req.Header.Add("Authorization", "2")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Error("Should return 401")
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/protected_entity/1", nil)
	req.Header.Add("Authorization", "2")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Error("Should return 401")
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/api/protected_entity/1", nil)
	req.Header.Add("Authorization", "2")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Error("Should return 401")
	}

	//test for delete
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/protected_entity/1", nil)
	req.Header.Add("Authorization", "2")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Error("Should return 401")
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/protected_entity/1", nil)
	req.Header.Add("Authorization", "1")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Error("Should return 204")
	}

}
