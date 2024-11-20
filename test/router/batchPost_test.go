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

func TestApiRouter_BatchPost(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	getDB().Exec("DELETE FROM tests")
	r := SetupRouter()

	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := NewApiRouter(
		*repo,
		map[route.RouteType]route.Route{route.BatchPost: {RouteType: route.BatchPost}},
	)
	test_.AllowRoutes(r)

	list := []Test{
		{entity.BaseEntity{Id: 2, Name: "test"}},
		{entity.BaseEntity{Id: 3, Name: "test"}},
	}
	w := httptest.NewRecorder()

	serializer := serializer.NewSerializer(format.JSON)

	jsonEntity, _ := serializer.Serialize(list)
	req, _ := http.NewRequest("POST", "/api/test", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Error("Failed to create")
	}
	list = list[:0]
	serializer.Deserialize(w.Body.String(), &list)
	if len(list) != 2 {
		fmt.Println(w.Body.String())
		t.Error("Failed to create")
	}
	if object, err := repo.GetByID(2); object == nil || err != nil {
		t.Error("Failed to create")
	}

	req, _ = http.NewRequest("POST", "/api/test", bytes.NewReader([]byte("salt")))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("should have failed")
	}

}

// fais un batch psot de 1000 entities
func BenchmarkApiRouter_BatchPost(b *testing.B) {
	getDB().AutoMigrate(&Test{})
	getDB().Exec("DELETE FROM tests")
	r := SetupRouter()

	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := NewApiRouter(
		*repo,
		map[route.RouteType]route.Route{route.BatchPost: {RouteType: route.BatchPost}},
	)
	test_.AllowRoutes(r)
	list := []Test{}
	for range 100000 {
		list = append(list, Test{entity.BaseEntity{Name: "test"}})
	}
	w := httptest.NewRecorder()

	serializer := serializer.NewSerializer(format.JSON)

	jsonEntity, _ := serializer.Serialize(list)
	req, _ := http.NewRequest("POST", "/api/test", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		b.Error("Failed to create")
	}
	if objects, err := repo.GetAll(map[string]string{"id": "asc"}); len(objects) != 100000 || err != nil {
		b.Error("Failed to create")
	}

}
