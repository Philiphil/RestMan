package apiman_test

import (
	"fmt"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"

	. "github.com/philiphil/apiman"
	"github.com/philiphil/apiman/method"
	"github.com/philiphil/apiman/orm"
	"github.com/philiphil/apiman/orm/entity"
	"github.com/philiphil/apiman/orm/repository"
)

func TestApiRouter_delete(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	r := SetupRouter()

	repo := orm.NewORM[Test](repository.NewRepository[Test, Test](getDB()))
	test_ := NewApiRouter[Test](
		*repo,
		method.DefaultApiMethods(),
	)
	test_.AllowRoutes(r)
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	test_.Get(context)

	entity := Test{entity.Entity{Id: 1}}
	repo.Create(&entity)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/test/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Println(w.Body.String())
		t.Error("Failed to start server")
	}

	req, _ = http.NewRequest("DELETE", "/api/test/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("Failed to delete")
	}

	if object, err := repo.GetByID(1); object != nil && err != nil {
		t.Error("Failed to delete")
	}

	w = httptest.NewRecorder() //not sure why the recorder has to be reset sometimes
	req, _ = http.NewRequest("DELETE", "/api/test/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		fmt.Println(w.Code)
		fmt.Println(w.Body.String())
		t.Error("Failed to delete")
	}

}
