package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository"
	"github.com/philiphil/restman/route"
	. "github.com/philiphil/restman/router"
)

func TestApiRouter_BatchDelete(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	r := SetupRouter()

	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)
	test_.Routes[route.BatchDelete] = route.Route{RouteType: route.BatchDelete}
	w := httptest.NewRecorder()
	repo.Create(&Test{entity.BaseEntity{Id: 1}})
	repo.Create(&Test{entity.BaseEntity{Id: 5}})
	repo.Create(&Test{entity.BaseEntity{Id: 8}})
	repo.Create(&Test{entity.BaseEntity{Id: 99}})

	test_.AllowRoutes(r)

	req, _ := http.NewRequest("DELETE", "/api/test?ids=1&ids=99", nil)
	r.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	w = httptest.NewRecorder() //not sure why the recorder has to be reset sometimes
	req, _ = http.NewRequest("DELETE", "/api/test/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Error("Failed to delete")
	}

	w = httptest.NewRecorder() //not sure why the recorder has to be reset sometimes
	req, _ = http.NewRequest("DELETE", "/api/test/99", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Error("Failed to delete")
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/test?ids=929", nil)
	r.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
