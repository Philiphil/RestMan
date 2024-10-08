package restman_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/philiphil/restman"
	"github.com/philiphil/restman/method"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/repository"
)

func TestApiRouter_put(t *testing.T) {
	getDB().AutoMigrate(&Test{})
	getDB().Exec("DELETE FROM tests")
	r := SetupRouter()

	repo := orm.NewORM[Test](repository.NewRepository[Test, Test](getDB()))
	test_ := NewApiRouter[Test](
		*repo,
		method.DefaultApiMethods(),
	)
	test_.AllowRoutes(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", "/api/test/2", bytes.NewBuffer([]byte(`{"name":"test"}`)))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Error("should be no content")
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/test/2", bytes.NewBuffer([]byte(`{"name":"test2"}`)))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var medium Test
	getDB().First(&medium, 2)
	if medium.Name != "test2" {
		t.Error("should be test2")
	}
	//faillures
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/test/2", bytes.NewBuffer([]byte(`{"name":"test2"}`)))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "afesfs")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Error("should be bad request")
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/test/2", bytes.NewBuffer([]byte(`{"name":"test2"`)))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Error("should be bad request")
	}

}
