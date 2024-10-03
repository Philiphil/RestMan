package apiman_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	. "github.com/philiphil/apiman"
	"github.com/philiphil/apiman/format"
	"github.com/philiphil/apiman/method"
	"github.com/philiphil/apiman/orm"
	"github.com/philiphil/apiman/orm/entity"
	"github.com/philiphil/apiman/orm/repository"
	"github.com/philiphil/apiman/serializer"
)

func TestApiRouter_post(t *testing.T) {
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
	jsonEntity, _ := json.Marshal(entity)
	req, _ := http.NewRequest("POST", "/api/test", bytes.NewReader(jsonEntity))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Error("Failed to create")
	}
	serializer := serializer.NewSerializer(format.JSON)
	medium := Test{}
	serializer.Deserialize(w.Body.String(), &medium)
	if medium.Id != 2 {
		t.Error("Failed to create")
	}
	if object, err := repo.GetByID(2); object == nil || err != nil {
		t.Error("Failed to create")
	}

}
