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

func TestApiRouter_GetXML(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	getDB().Exec("DELETE FROM tests")
	r := SetupRouter()

	repo := orm.NewORM[Test](gormrepository.NewRepository[Test, Test](getDB()))
	test_ := NewApiRouter[Test](
		*repo,
		route.DefaultApiRoutes(),
	)
	test_.AllowRoutes(r)

	entity := Test{entity.BaseEntity{Id: 42}}
	repo.Create(&entity)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/test/42", nil)
	req.Header.Add("Accept", "application/xml")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Println(w.Body.String())
		t.Error("Failed to get XML response")
	}

	// Check content type
	if w.Header().Get("Content-Type") != "application/xml; charset=utf-8" {
		t.Errorf("Expected Content-Type application/xml but got %s", w.Header().Get("Content-Type"))
	}

	// Check XML header and content
	body := w.Body.String()
	if len(body) < 6 || body[:6] != "<?xml " {
		t.Error("Expected XML header in response")
	}

	// Just verify XML contains the ID (don't try to deserialize nested struct)
	if !containsString(body, "<Id>42</Id>") && !containsString(body, "<id>42</id>") {
		fmt.Println("Body:", body)
		t.Error("Expected XML to contain ID 42")
	}
}

func TestApiRouter_GetListCSV(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	getDB().Exec("DELETE FROM tests")
	r := SetupRouter()

	repo := orm.NewORM[Test](gormrepository.NewRepository[Test, Test](getDB()))
	test_ := NewApiRouter[Test](
		*repo,
		route.DefaultApiRoutes(),
	)
	test_.AllowRoutes(r)

	// Create test entities
	entity1 := Test{entity.BaseEntity{Id: 1}}
	entity2 := Test{entity.BaseEntity{Id: 2}}
	entity3 := Test{entity.BaseEntity{Id: 3}}
	repo.Create(&entity1)
	repo.Create(&entity2)
	repo.Create(&entity3)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/test", nil)
	req.Header.Add("Accept", "application/csv")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Println(w.Body.String())
		t.Error("Failed to get CSV response")
	}

	// Check content type
	if w.Header().Get("Content-Type") != "text/csv" {
		t.Errorf("Expected Content-Type text/csv but got %s", w.Header().Get("Content-Type"))
	}

	body := w.Body.String()
	fmt.Println("CSV Body:", body)

	// Check CSV format
	if len(body) == 0 {
		t.Error("Expected CSV content but got empty response")
	}

	// Check it contains Id header
	if !containsString(body, "Id") {
		t.Error("Expected CSV to contain 'Id' header")
	}

	// Check it has the IDs in the content
	if !containsString(body, ",1,") && !containsString(body, "1,") {
		t.Error("Expected CSV to contain ID 1")
	}
	if !containsString(body, ",2,") && !containsString(body, "2,") {
		t.Error("Expected CSV to contain ID 2")
	}
	if !containsString(body, ",3,") && !containsString(body, "3,") {
		t.Error("Expected CSV to contain ID 3")
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || containsString(s[1:], substr)))
}
