package router_test

import (
	"fmt"
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

func TestApiRouter_Get(t *testing.T) {
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
	test_.Get(context)

	entity := Test{entity.BaseEntity{Id: 1}}
	repo.Create(&entity)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/test/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Println(w.Body.String())
		t.Error("Failed to start server")
	}
	serializer := serializer.NewSerializer(format.JSON)
	medium := Test{}
	serializer.Deserialize(w.Body.String(), &medium)
	if medium.Id != entity.Id {
		t.Error("Expected entity with id " + fmt.Sprint(entity.Id) + " but got " + fmt.Sprint(medium.Id))
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/test/1", nil)
	req.Header.Add("Accept", "Wrong")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotAcceptable {
		t.Error("should not accept")
	}
}
