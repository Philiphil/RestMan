package router_test

import (
	"fmt"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"

	"github.com/philiphil/restman/format"
	"github.com/philiphil/restman/method"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/repository"
	. "github.com/philiphil/restman/router"
	"github.com/philiphil/restman/serializer"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	r.Use(gin.Recovery())
	return r
}

func TestApiRouter_GetList(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	getDB().Exec("DELETE FROM tests")
	r := SetupRouter()

	repo := orm.NewORM[Test](repository.NewRepository[Test, Test](getDB()))
	test_ := NewApiRouter[Test](
		*repo,
		method.DefaultApiMethods(),
	)
	test_.AllowRoutes(r)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Failed to start server")
	}
	serializer := serializer.NewSerializer(format.JSON)
	entities := make([]Test, 0)
	serializer.Deserialize(w.Body.String(), &entities)
	if len(entities) != 0 {
		t.Error("Expected mo entity")
	}

	entity := Test{entity.Entity{Id: 1}}
	repo.Create(&entity)
	w.Body.Reset()
	req, _ = http.NewRequest("GET", "/api/test", nil)
	r.ServeHTTP(w, req)

	serializer.Deserialize(w.Body.String(), &entities)

	if len(entities) != 1 {
		t.Error("Expected 1 entity")
	}
	if entities[0].Id != entity.Id {
		t.Error("Expected entity with id " + fmt.Sprint(entity.Id) + " but got " + fmt.Sprint(entities[0].Id))
	}
}

func TestApiRouter_GetListPaginated(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	r := SetupRouter()

	repo := orm.NewORM[Test](repository.NewRepository[Test, Test](getDB()))
	test_ := NewApiRouter[Test](
		*repo,
		method.DefaultApiMethods(),
	)
	test_.AllowRoutes(r)
	w := httptest.NewRecorder()
	serializer := serializer.NewSerializer(format.JSON)
	entities := make([]Test, 0)
	e := Test{}
	for i := 0; i < 11; i++ {
		e.Id = entity.ID(i + 1)
		repo.Create(&e)
	}

	req, _ := http.NewRequest("GET", "/api/test?pagination=true&itemsPerPage=5", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Failed to start server")
	}
	entities = entities[:0]
	serializer.Deserialize(w.Body.String(), &entities)
	if len(entities) != 5 {
		fmt.Println(w.Body.String())
		t.Error("Expected 5 entities")
	}
	for i, e := range entities {
		if e.GetId() != entity.ID(i+1) {
			t.Error("Expected entity with id " + fmt.Sprint(i) + " but got " + fmt.Sprint(e.Id))
		}
	}
	w.Body.Reset()
	req, _ = http.NewRequest("GET", "/api/test?pagination=true&itemsPerPage=5&page=2", nil)
	r.ServeHTTP(w, req)
	entities = entities[:0]
	serializer.Deserialize(w.Body.String(), &entities)
	if len(entities) != 5 {
		fmt.Println(w.Body.String())
		t.Error("Expected 5 entities")
	}
	for i, e := range entities {
		if e.GetId() != entity.ID(i+1+5) {
			t.Error("Expected entity with id " + fmt.Sprint(i) + " but got " + fmt.Sprint(e.Id))
		}
	}

	w.Body.Reset()
	req, _ = http.NewRequest("GET", "/api/test?pagination=true&itemsPerPage=5&page=3", nil)
	r.ServeHTTP(w, req)
	entities = entities[:0]
	serializer.Deserialize(w.Body.String(), &entities)
	if len(entities) != 1 {
		fmt.Println(w.Body.String())
		t.Error("Expected 1 entities")
	}
	for i, e := range entities {
		if e.GetId() != entity.ID(i+1+10) {
			t.Error("Expected entity with id " + fmt.Sprint(i) + " but got " + fmt.Sprint(e.Id))
		}
	}

	w.Body.Reset()
	req, _ = http.NewRequest("GET", "/api/test?pagination=true&itemsPerPage=50&page=1", nil)
	r.ServeHTTP(w, req)
	entities = entities[:0]
	serializer.Deserialize(w.Body.String(), &entities)
	if len(entities) != 11 {
		fmt.Println(w.Body.String())
		t.Error("Expected 11 entities")
	}

	w.Body.Reset()
	req, _ = http.NewRequest("GET", "/api/test?pagination=true&itemsPerPage=50&page=2", nil)
	r.ServeHTTP(w, req)
	entities = entities[:0]
	serializer.Deserialize(w.Body.String(), &entities)
	if len(entities) != 0 {
		fmt.Println(w.Body.String())
		t.Error("Expected no entities")
	}
}

func TestApiRouter_GetListJSONLD(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	r := SetupRouter()

	repo := orm.NewORM[Test](repository.NewRepository[Test, Test](getDB()))
	test_ := NewApiRouter[Test](
		*repo,
		method.DefaultApiMethods(),
	)
	test_.AllowRoutes(r)
	e := Test{}
	for i := 0; i < 11; i++ {
		e.Id = entity.ID(i + 1)
		repo.Create(&e)
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/test?pagination=true&itemsPerPage=5&page=1", nil)
	req.Header.Add("Accept", "application/ld+json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Failed to start server")
	}
	serializer := serializer.NewSerializer(format.JSONLD)
	entities := make(map[string]any, 0)
	serializer.Deserialize(w.Body.String(), &entities)
	if entities["hydra:totalItems"].(float64) != 5 {
		fmt.Println(entities["hydra:totalItems"])
		t.Error("Expected 5 entities")
	}
	if entities["hydra:view"].(map[string]any)["hydra:first"] != "/api/test?page=1" {
		t.Error("Expected /api/test?page=1")
	}

	if entities["hydra:view"].(map[string]any)["hydra:next"] != "/api/test?page=2" {
		t.Error("Expected /api/test?page=2")
	}

}

func TestApiRouter_GetListErrors(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	r := SetupRouter()

	repo := orm.NewORM[Test](repository.NewRepository[Test, Test](getDB()))
	test_ := NewApiRouter[Test](
		*repo,
		method.DefaultApiMethods(),
	)
	test_.AllowRoutes(r)
	e := Test{}
	for i := 0; i < 11; i++ {
		e.Id = entity.ID(i + 1)
		repo.Create(&e)
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/test?pagination=true&itemsPerPage=5&page=1", nil)
	req.Header.Add("Accept", "Wrong")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotAcceptable {
		t.Error("should not accept")
	}
}
