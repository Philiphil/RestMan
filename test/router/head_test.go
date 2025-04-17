package router_test

import (
	"strconv"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"

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

	repo := orm.NewORM[Test](gormrepository.NewRepository[Test, Test](getDB()))
	test_ := NewApiRouter[Test](
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
	req, _ = http.NewRequest("HEAD", "/api/test/1", nil)
	r.ServeHTTP(w, req)
	serializer := serializer.NewSerializer(format.JSON)
	json, _ := serializer.Serialize(entity)
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
