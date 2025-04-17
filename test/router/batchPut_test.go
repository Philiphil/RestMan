package router_test

import (
	"bytes"
	"fmt"
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

func TestApiRouter_BatchPut(t *testing.T) {

	getDB().AutoMigrate(&Test{})
	r := SetupRouter()

	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)
	test_.Routes[route.BatchPut] = route.Route{RouteType: route.BatchPut}
	test_.AllowRoutes(r)
	w := httptest.NewRecorder()
	repo.Create(&Test{entity.BaseEntity{Id: 1, Name: "test"}})
	repo.Create(&Test{entity.BaseEntity{Id: 5, Name: "test"}})

	//try to edit 1 et 5 name and create a 8
	list := []Test{
		{entity.BaseEntity{Id: 1, Name: "test1"}},
		{entity.BaseEntity{Id: 5, Name: "test5"}},
		{entity.BaseEntity{Id: 8, Name: "test8"}},
	}
	serializer := serializer.NewSerializer(format.JSON)
	jsonEntity, _ := serializer.Serialize(list)
	req, _ := http.NewRequest("PUT", "/api/test", bytes.NewReader([]byte(jsonEntity)))
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
	if object, err := repo.GetByID(8); object == nil || err != nil || object.Name != "test8" {
		t.Error("Failed to create")
	}

	//try to create no id
	list = []Test{
		{entity.BaseEntity{Name: "test8"}},
	}
	jsonEntity, _ = serializer.Serialize(list)
	req, _ = http.NewRequest("PUT", "/api/test", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Error(w.Code)
	}

	//not valid body
	req, _ = http.NewRequest("PUT", "/api/test", bytes.NewReader([]byte("whatever")))
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
	req, _ = http.NewRequest("PUT", "/api/test", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		fmt.Print(w.Code)
		t.Error(w.Code)
	}

	//wrong accept
	list = []Test{
		{entity.BaseEntity{Id: 1, Name: "test1"}},
		{entity.BaseEntity{Id: 5, Name: "test5"}},
	}
	jsonEntity, _ = serializer.Serialize(list)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/test", bytes.NewReader([]byte(jsonEntity)))
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
	req, _ = http.NewRequest("PUT", "/api/test", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "applicatifsefson/jsocscscsn")
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Error("Failed request")
	}

}
