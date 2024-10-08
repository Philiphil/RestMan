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

func TestApiRouter_delete(t *testing.T) {
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
	test_.Delete(context)

	entity := Test{entity.Entity{Id: 1}}
	repo.Create(&entity)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/test/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("Failed to start server")
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/test/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		fmt.Println(w.Body.String())
		fmt.Println(w.Code)
		t.Error("Failed to delete")
	}

	if object, err := repo.GetByID(1); object != nil && err != nil {
		t.Error("Failed to delete")
	}

	w = httptest.NewRecorder() //not sure why the recorder has to be reset sometimes
	req, _ = http.NewRequest("DELETE", "/api/test/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Error("Failed to delete")
	}

}
