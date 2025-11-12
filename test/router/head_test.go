package router_test

import (
	"context"
	"strconv"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/philiphil/restman/format"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository"
	"github.com/philiphil/restman/route"
	. "github.com/philiphil/restman/router"
	"github.com/philiphil/restman/serializer"
)

func TestApiRouter_Head(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	getDB().Exec("DELETE FROM tests")
	r := SetupRouter()

	repo := orm.NewORM(gormrepository.NewRepository[Test](getDB()))
	test_ := NewApiRouter(
		*repo,
		route.DefaultApiRoutes(),
	)
	test_.AllowRoutes(r)
	// Removed direct Head handler invocation; use HTTP flow only
	ctx := context.Background()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("HEAD", "/api/test/1", nil)
	req.Header.Set("Accept", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Error("should be not found")
	}

	entity := Test{entity.BaseEntity{Id: 1}}
	repo.Create(ctx, &entity)
	req, _ = http.NewRequest("HEAD", "/api/test/1", nil)
	r.ServeHTTP(w, req)
	s := serializer.NewSerializer(format.JSON)
	json, _ := s.Serialize(entity, serializer.SerializeOptions{})
	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Content-Type should be application/json")
	}
	if w.Header().Get("Content-Length") != strconv.Itoa(len(json)) {
		t.Error("Content-Length should be equal to the length of the json")
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("HEAD", "/api/test/1", nil)
	req.Header.Set("Accept", "cac/esfes")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotAcceptable {
		t.Error("should be not acceptable")
	}

}
