package restman_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	. "github.com/philiphil/restman"
	"github.com/philiphil/restman/format"
	"github.com/philiphil/restman/method"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/repository"
	"github.com/philiphil/restman/serializer"
)

func TestApiRouter_post(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	getDB().Exec("DELETE FROM tests")
	r := SetupRouter()

	repo := orm.NewORM[Test](repository.NewRepository[Test, Test](getDB()))
	test_ := NewApiRouter[Test](
		*repo,
		method.DefaultApiMethods(),
	)
	test_.AllowRoutes(r)

	entity := Test{entity.Entity{Id: 2, Name: "test"}}
	w := httptest.NewRecorder()

	serializer := serializer.NewSerializer(format.JSON)

	jsonEntity, _ := serializer.Serialize(entity)
	req, _ := http.NewRequest("POST", "/api/test", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Error("Failed to create")
	}
	medium := Test{}
	serializer.Deserialize(w.Body.String(), &medium)
	if medium.Id != 2 {
		fmt.Println(w.Body.String())
		t.Error("Failed to create")
	}
	if object, err := repo.GetByID(2); object == nil || err != nil {
		t.Error("Failed to create")
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/test", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code == http.StatusCreated {
		t.Error("wrong status code")
	}

}

func TestApiRouter_postFailure(t *testing.T) {
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

	entity := Test{entity.Entity{Id: 2, Name: "test"}}
	w := httptest.NewRecorder()

	serializer := serializer.NewSerializer(format.JSON)

	jsonEntity, _ := serializer.Serialize(entity)
	req, _ := http.NewRequest("POST", "/api/test", bytes.NewReader([]byte(jsonEntity+"salt")))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("should have failed")
	}

	req, _ = http.NewRequest("POST", "/api/test", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "bad")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("should have failed")
	}

	req, _ = http.NewRequest("POST", "/api/test", bytes.NewReader([]byte(jsonEntity)))
	req.Header.Add("Accept", "bad")
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("should have failed")
	}

}
