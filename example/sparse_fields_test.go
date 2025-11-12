package example_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSparseFields(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/test?fields=name", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		panic("Failed to start server")
	}
	var responseData []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseData)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
		return
	}

	for _, item := range responseData {
		if len(item) != 1 {
			t.Errorf("Expected only 1 field, got %d fields", len(item))
		}
		if _, exists := item["name"]; !exists {
			t.Errorf("Expected field 'name' to be present")
		}
	}

}
