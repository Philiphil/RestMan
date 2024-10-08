package restman_test

import (
	"fmt"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"

	. "github.com/philiphil/restman"
	"github.com/philiphil/restman/method"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/repository"
)

func TestApiRouter_Options(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	getDB().Exec("DELETE FROM tests")
	r := SetupRouter()

	repo := orm.NewORM[Test](repository.NewRepository[Test, Test](getDB()))
	test_ := NewApiRouter[Test](
		*repo,
		method.DefaultApiMethods(),
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

	entity := Test{entity.Entity{Id: 1}}
	repo.Create(&entity)
	//serializer mon truc en json et compare avec le json de la reponse
	req, _ = http.NewRequest("OPTIONS", "/api/test/1", nil)
	r.ServeHTTP(w, req)

	if w.Header().Get("Allow") != "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS" {
		fmt.Println(w.Header().Get("Allow"))
		t.Error("Allow should be GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS")
	}
	if w.Header().Get("Content-Length") != "0" {
		t.Error("Content-Length should be 0")
	}

}
