package router_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/philiphil/restman/format"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository"
	"github.com/philiphil/restman/route"
	. "github.com/philiphil/restman/router"
	"github.com/philiphil/restman/serializer"
)

func TestApiRouter_BatchPatch(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	r := SetupRouter()

	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)
	test_.Routes[route.BatchPatch] = route.Route{RouteType: route.BatchPatch}
	test_.AllowRoutes(r)
	w := httptest.NewRecorder()
	repo.Create(&Test{entity.BaseEntity{Id: 1, Name: "test"}})
	repo.Create(&Test{entity.BaseEntity{Id: 5, Name: "test"}})

	//try to edit 1 and 5
	list := []Test{
		{entity.BaseEntity{Id: 1, Name: "test1"}},
		{entity.BaseEntity{Id: 5, Name: "test5"}},
	}
	serializer := serializer.NewSerializer(format.JSON)
	jsonEntity, _ := serializer.Serialize(list)
	req, _ := http.NewRequest("PATCH", "/api/test", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Error("Failed request")
	}
	//check if the entities have been updated
	if object, err := repo.GetByID(1); object == nil || err != nil || object.Name != "test1" {
		t.Error("Failed to update")
	}
	if object, err := repo.GetByID(5); object == nil || err != nil || object.Name != "test5" {
		t.Error("Failed to update")
	}

	//try to create a new one
	list = []Test{
		{entity.BaseEntity{Id: 4545, Name: "test8"}},
	}
	repo.Delete(&Test{entity.BaseEntity{Id: 4545}})
	jsonEntity, _ = serializer.Serialize(list)
	req, _ = http.NewRequest("PATCH", "/api/test", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Error(w.Code)
	}

	//try to create no id
	list = []Test{
		{entity.BaseEntity{Name: "test8"}},
	}
	jsonEntity, _ = serializer.Serialize(list)
	req, _ = http.NewRequest("PATCH", "/api/test", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Error(w.Code)
	}

	//not valid body
	req, _ = http.NewRequest("PATCH", "/api/test", bytes.NewReader([]byte("whatever")))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Error(w.Code)
	}

	//empty list
	list = []Test{}
	jsonEntity, _ = serializer.Serialize(list)
	req, _ = http.NewRequest("PATCH", "/api/test", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Error(w.Code)
	}

	//wrong accept
	list = []Test{
		{entity.BaseEntity{Id: 1, Name: "test1"}},
		{entity.BaseEntity{Id: 5, Name: "test5"}},
	}
	jsonEntity, _ = serializer.Serialize(list)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/api/test", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Accept", "application/jsfeson")
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotAcceptable {
		t.Error("Failed request")
	}

	//wrong content type

	list = []Test{
		{entity.BaseEntity{Id: 1, Name: "test1"}},
		{entity.BaseEntity{Id: 5, Name: "test5"}},
	}
	jsonEntity, _ = serializer.Serialize(list)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/api/test", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "applicatifsefson/jsocscscsn")
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Error("Failed request")
	}

}
