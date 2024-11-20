package router_test

import (
	"slices"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"

	"maps"

	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository"
	"github.com/philiphil/restman/route"
	"github.com/philiphil/restman/router"
)

func countOccurrences(slice []string, value string) int {
	count := 0
	for _, v := range slice {
		if v == value {
			count++
		}
	}
	return count
}

func TestApiRouter_Options(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	getDB().Exec("DELETE FROM tests")
	r := SetupRouter()

	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := router.NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)
	test_.AllowRoutes(r)
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	test_.Head(context)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("HEAD", "/api/test/1", nil)
	req.Header.Set("Accept", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Error("should be not found")
	}

	entity := Test{entity.BaseEntity{Id: 1}}
	repo.Create(&entity)

	req, _ = http.NewRequest("OPTIONS", "/api/test/1", nil)
	r.ServeHTTP(w, req)
	if w.Header().Get("Content-Length") != "0" {
		t.Error("Content-Length should be 0")
	}
	allowed := strings.Split(w.Header().Get("Allow"), ",")
	expected := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	for _, method := range expected {
		if !slices.Contains(allowed, method) {
			t.Errorf("Method %s should be allowed", method)
		}
	}

}

func TestApiRouter_Options_Batch(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	getDB().Exec("DELETE FROM tests")
	r := SetupRouter()

	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	//merge batchoperation with normal
	allRoutes := route.DefaultApiRoutes()
	maps.Copy(allRoutes, route.BatchOperations())
	test_ := router.NewApiRouter(
		*repo,
		allRoutes,
	)
	test_.AllowRoutes(r)
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	test_.Head(context)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("HEAD", "/api/test/1", nil)
	req.Header.Set("Accept", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Error("should be not found")
	}
	entity := Test{entity.BaseEntity{Id: 1}}
	repo.Create(&entity)
	req, _ = http.NewRequest("OPTIONS", "/api/test/1", nil)
	r.ServeHTTP(w, req)
	allowed := strings.Split(w.Header().Get("Allow"), ",")
	expected := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	for _, method := range expected {
		if !slices.Contains(allowed, method) {
			t.Errorf("Method %s should be allowed", method)
		}
		if countOccurrences(allowed, method) > 1 {
			t.Errorf("Method %s should be only once", method)
		}
	}

}
