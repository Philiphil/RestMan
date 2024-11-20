package router_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestApiRouter_IsBatchGetOrGetList(t *testing.T) {
	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
		configuration.Pagination(false),
	)
	w := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(w)

	if test_.IsBatchGetOrGetList(context) != route.GetList {
		//batchGet is not allowed
		t.Error("Expected GetList")
	}
	delete(test_.Routes, route.GetList)
	test_.Routes[route.BatchGet] = route.Route{RouteType: route.BatchGet} //this is a hack

	if test_.IsBatchGetOrGetList(context) != route.BatchGet {
		//getlist is not allowed
		t.Error("Expected BatchGet")
	}
	test_.Routes[route.GetList] = route.Route{RouteType: route.GetList} //this is a hack
	context.Request = httptest.NewRequest("GET", "/test", nil)
	if test_.IsBatchGetOrGetList(context) != route.GetList {
		//both are allowed but there's no id parameters
		t.Error("Expected GetList")
	}
	//cannot test the case where it is a batchGet because it requires a query parameter
	//and the context.Request is not a real request
}

func TestApiRouter_GetListOrBatchGet(t *testing.T) {
	getDB().Migrator().DropTable(&Test{})
	getDB().AutoMigrate(&Test{})
	r := SetupRouter()

	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
		configuration.Pagination(false),
	)
	test_.Routes[route.BatchGet] = route.Route{RouteType: route.BatchGet}
	w := httptest.NewRecorder()
	repo.Create(&Test{entity.BaseEntity{Id: 1}})
	repo.Create(&Test{entity.BaseEntity{Id: 5}})
	repo.Create(&Test{entity.BaseEntity{Id: 8}})
	repo.Create(&Test{entity.BaseEntity{Id: 99}})
	test_.AllowRoutes(r)
	serializer := serializer.NewSerializer(format.JSON)

	w.Body.Reset()
	req, _ := http.NewRequest("GET", "/api/test?ids=1&ids=99", nil)
	r.ServeHTTP(w, req)
	entities := make([]Test, 0)
	serializer.Deserialize(w.Body.String(), &entities)
	if len(entities) != 2 {
		fmt.Println(entities)
		t.Error("Expected 2 entities")
	}

	w.Body.Reset()
	req, _ = http.NewRequest("GET", "/api/test?ids[]=1&ids[]=99", nil)
	r.ServeHTTP(w, req)
	entities = entities[:0]
	serializer.Deserialize(w.Body.String(), &entities)
	if len(entities) != 2 {
		fmt.Println(entities)
		t.Error("Expected 2 entities")
	}

	w.Body.Reset()
	req, _ = http.NewRequest("GET", "/api/test?ids=1,99", nil)
	r.ServeHTTP(w, req)
	entities = entities[:0]
	serializer.Deserialize(w.Body.String(), &entities)
	if len(entities) != 2 {
		fmt.Println(entities)
		t.Error("Expected 2 entities")
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/test?ids=1,979", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Error("Expected 404")
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/test?ids=", nil)
	r.ServeHTTP(w, req)
	entities = entities[:0]
	serializer.Deserialize(w.Body.String(), &entities)
	if len(entities) != 0 {
		t.Error("Expected 0 entities")
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/test?ids=1,99", nil)
	req.Header.Add("Accept", "Wrong")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotAcceptable {
		fmt.Println(w.Code)
		t.Error("should not accept")
	}

}
