package router_test

import (
	"fmt"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"

	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/format"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository"
	"github.com/philiphil/restman/route"
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

	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
		configuration.Pagination(false),
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

	entity := Test{entity.BaseEntity{Id: 1}}
	repo.Create(&entity)
	w.Body.Reset()
	req, _ = http.NewRequest("GET", "/api/test", nil)
	r.ServeHTTP(w, req)

	serializer.Deserialize(w.Body.String(), &entities)

	if len(entities) != 1 {
		fmt.Println(entities)
		t.Error("Expected 1 entity")
	}
	if entities[0].Id != entity.Id {
		t.Error("Expected entity with id " + fmt.Sprint(entity.Id) + " but got " + fmt.Sprint(entities[0].Id))
	}
}

func TestApiRouter_GetListPaginated(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	r := SetupRouter()

	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
		configuration.Pagination(true),
		configuration.PaginationClientControl(true),
		configuration.ItemPerPage(5),
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
	req, _ = http.NewRequest("GET", "/api/test?pagination=true&page=2", nil)
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
	req, _ = http.NewRequest("GET", "/api/test?page=3", nil)
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

	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
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

	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
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

	test_.Configuration[configuration.MaxItemPerPageType] = configuration.MaxItemPerPage(5)
	test_.Configuration[configuration.ItemPerPageType] = configuration.ItemPerPage(6)

	req, _ = http.NewRequest("GET", "/api/test?pagination=true", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotAcceptable {
		t.Error("should not accept")
	}

}

func TestApiRouter_IsPaginationEnabled(t *testing.T) {
	r := SetupRouter()
	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)

	c := gin.Context{}
	if _, err := test_.IsPaginationEnabled(&c); err != nil {
		t.Error("Expected no errors")
	}
	delete(test_.Configuration, configuration.PaginationType)
	if _, err := test_.IsPaginationEnabled(&c); err == nil {
		t.Error("Expected errors")
	}

	test_ = NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)
	delete(test_.Configuration, configuration.PaginationClientControlType)
	if _, err := test_.IsPaginationEnabled(&c); err == nil {
		t.Error("Expected errors")
	}

	test_ = NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)
	test_.Configuration[configuration.PaginationClientControlType] = configuration.PaginationClientControl(true)
	delete(test_.Configuration, configuration.PaginationParameterNameType)
	if _, err := test_.IsPaginationEnabled(&c); err == nil {
		t.Error("Expected errors")
	}

	test_ = NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
		configuration.Configuration{Type: configuration.PaginationClientControlType, Values: []string{"test"}},
	)
	if _, err := test_.IsPaginationEnabled(&c); err == nil {
		t.Error("Expected errors")
	}

	test_ = NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
		configuration.Configuration{Type: configuration.PaginationType, Values: []string{"test"}},
	)
	if _, err := test_.IsPaginationEnabled(&c); err == nil {
		t.Error("Expected errors")
	}

	w := httptest.NewRecorder()
	test_ = NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)
	test_.AllowRoutes(r)
	req, _ := http.NewRequest("GET", "/api/test?pagination=true&itemsPerPage=5&page=1", nil)
	req.Header.Add("Accept", "Wrong")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotAcceptable {
		t.Error("should not accept")
	}
}

func TestApiRouter_GetItemPerPage(t *testing.T) {
	r := SetupRouter()
	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)

	c := gin.Context{}
	if _, err := test_.GetItemPerPage(&c); err != nil {
		t.Error("Expected no errors")
	}
	delete(test_.Configuration, configuration.ItemPerPageType)
	if _, err := test_.GetItemPerPage(&c); err == nil {
		t.Error("Expected errors")
	}
	test_ = NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)
	delete(test_.Configuration, configuration.MaxItemPerPageType)
	if _, err := test_.GetItemPerPage(&c); err == nil {
		t.Error("Expected errors")
	}

	test_ = NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)
	delete(test_.Configuration, configuration.ItemPerPageParameterNameType)
	if _, err := test_.GetItemPerPage(&c); err == nil {
		t.Error("Expected errors")
	}

	test_ = NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
		configuration.Configuration{Type: configuration.MaxItemPerPageType, Values: []string{"test"}},
	)
	if _, err := test_.GetItemPerPage(&c); err == nil {
		t.Error("Expected errors")
	}
	test_ = NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)
	w := httptest.NewRecorder()
	test_.AllowRoutes(r)
	req, _ := http.NewRequest("GET", "/api/test?pagination=true&itemsPerPage=5&page=1", nil)
	req.Header.Add("Accept", "Wrong")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotAcceptable {
		fmt.Println(w.Code)
		t.Error("should not accept")
	}
}

func TestApiRouter_SortOrder(t *testing.T) {
	r := SetupRouter()
	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)

	c := gin.Context{}
	if _, err := test_.GetSortOrder(&c); err != nil {
		t.Error("Expected no errors")
	}
	delete(test_.Configuration, configuration.SortingClientControlType)
	if _, err := test_.GetSortOrder(&c); err == nil {
		t.Error("Expected errors")
	}
	test_ = NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)
	delete(test_.Configuration, configuration.SortingType)
	if _, err := test_.GetSortOrder(&c); err == nil {
		t.Error("Expected errors")
	}

	test_ = NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)
	delete(test_.Configuration, configuration.SortingParameterNameType)
	if _, err := test_.GetSortOrder(&c); err == nil {
		t.Error("Expected errors")
	}

	test_ = NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)
	w := httptest.NewRecorder()
	test_.AllowRoutes(r)
	req, _ := http.NewRequest("GET", "/api/test?pagination=true&itemsPerPage=5&page=1&sort[id]=asc", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("Failed to start server")
	}
	serializer := serializer.NewSerializer(format.JSON)
	entities := make([]Test, 0)
	serializer.Deserialize(w.Body.String(), &entities)
	for i, e := range entities {
		if e.GetId() != entity.ID(i+1) {
			t.Error("Expected entity with id " + fmt.Sprint(i) + " but got " + fmt.Sprint(e.Id))
		}
	}
	w.Body.Reset()
	req, _ = http.NewRequest("GET", "/api/test?pagination=true&itemsPerPage=5&page=1&sort[id]=desc", nil)
	r.ServeHTTP(w, req)
	entities = entities[:0]
	serializer.Deserialize(w.Body.String(), &entities)
	for i, e := range entities {
		expected := entity.ID(11 - i)
		if e.GetId() != expected {
			t.Error("Expected entity with id " + fmt.Sprint(expected) + " but got " + fmt.Sprint(e.Id))
		}
	}

}
