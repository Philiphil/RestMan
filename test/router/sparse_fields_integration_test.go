package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	. "github.com/philiphil/restman/router"
)

type SparseUser struct {
	ID       int    `json:"id" groups:"test"`
	Name     string `json:"name" groups:"test"`
	Email    string `json:"email" groups:"test"`
	Password string `json:"password" groups:"test"`
}

func TestSparseFields_QueryParam_Simple(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/users/1?fields=id,name", nil)
	c.Request = req

	router := createTestRouterWithSparseFields[Test](t, true)

	sparseFieldsRequest, err := router.GetSparseFields(c, 0)
	if err != nil {
		t.Fatalf("GetSparseFields failed: %v", err)
	}

	if sparseFieldsRequest == nil {
		t.Fatal("Expected sparse fields request, got nil")
	}

	rootFields := sparseFieldsRequest.ResourceFields[""]
	if len(rootFields) != 2 {
		t.Errorf("Expected 2 root fields, got %d: %v", len(rootFields), rootFields)
	}

	expectedFields := map[string]bool{"id": true, "name": true}
	for _, field := range rootFields {
		if !expectedFields[field] {
			t.Errorf("Unexpected field: %s", field)
		}
	}
}

func TestSparseFields_QueryParam_JsonApi(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/users?fields[users]=id,name&fields[posts]=title", nil)
	c.Request = req

	router := createTestRouterWithSparseFields[Test](t, true)

	sparseFieldsRequest, err := router.GetSparseFields(c, 0)
	if err != nil {
		t.Fatalf("GetSparseFields failed: %v", err)
	}

	if sparseFieldsRequest == nil {
		t.Fatal("Expected sparse fields request, got nil")
	}

	usersFields := sparseFieldsRequest.ResourceFields["users"]
	if len(usersFields) != 2 {
		t.Errorf("Expected 2 users fields, got %d: %v", len(usersFields), usersFields)
	}

	postsFields := sparseFieldsRequest.ResourceFields["posts"]
	if len(postsFields) != 1 || postsFields[0] != "title" {
		t.Errorf("Expected posts field 'title', got %v", postsFields)
	}
}

func TestSparseFields_QueryParam_DotNotation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/users/1?fields=id,name,posts.title,posts.content", nil)
	c.Request = req

	router := createTestRouterWithSparseFields[Test](t, true)

	sparseFieldsRequest, err := router.GetSparseFields(c, 0)
	if err != nil {
		t.Fatalf("GetSparseFields failed: %v", err)
	}

	if sparseFieldsRequest == nil {
		t.Fatal("Expected sparse fields request, got nil")
	}

	rootFields := sparseFieldsRequest.ResourceFields[""]
	if len(rootFields) != 2 {
		t.Errorf("Expected 2 root fields, got %d: %v", len(rootFields), rootFields)
	}

	postsFields := sparseFieldsRequest.ResourceFields["posts"]
	if len(postsFields) != 2 {
		t.Errorf("Expected 2 posts fields, got %d: %v", len(postsFields), postsFields)
	}

	expectedPostsFields := map[string]bool{"title": true, "content": true}
	for _, field := range postsFields {
		if !expectedPostsFields[field] {
			t.Errorf("Unexpected posts field: %s", field)
		}
	}
}

func TestSparseFields_QueryParam_Disabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/users/1?fields=id,name", nil)
	c.Request = req

	router := createTestRouterWithSparseFields[Test](t, false)

	sparseFieldsRequest, err := router.GetSparseFields(c, 0)
	if err != nil {
		t.Fatalf("GetSparseFields failed: %v", err)
	}

	if sparseFieldsRequest != nil {
		t.Errorf("Expected nil when sparse fields disabled, got %v", sparseFieldsRequest)
	}
}

func TestSparseFields_QueryParam_EmptyValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/users/1?fields=", nil)
	c.Request = req

	router := createTestRouterWithSparseFields[Test](t, true)

	sparseFieldsRequest, err := router.GetSparseFields(c, 0)
	if err != nil {
		t.Fatalf("GetSparseFields failed: %v", err)
	}

	if sparseFieldsRequest != nil {
		t.Errorf("Expected nil for empty fields param, got %v", sparseFieldsRequest)
	}
}

func TestSparseFields_QueryParam_NoParam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	c.Request = req

	router := createTestRouterWithSparseFields[Test](t, true)

	sparseFieldsRequest, err := router.GetSparseFields(c, 0)
	if err != nil {
		t.Fatalf("GetSparseFields failed: %v", err)
	}

	if sparseFieldsRequest != nil {
		t.Errorf("Expected nil when no fields param, got %v", sparseFieldsRequest)
	}
}

func TestSparseFields_EndToEnd_JSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := SparseUser{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret123",
	}

	router := gin.New()
	router.GET("/users/:id", func(c *gin.Context) {
		apiRouter := createTestRouterWithSparseFields[Test](t, true)

		sparseFieldsRequest, _ := apiRouter.GetSparseFields(c, 0)

		c.Header("Content-Type", "application/json")

		if sparseFieldsRequest != nil {
			c.JSON(http.StatusOK, gin.H{
				"sparse_fields": sparseFieldsRequest.ResourceFields,
				"user":          user,
			})
		} else {
			c.JSON(http.StatusOK, user)
		}
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1?fields=id,name", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if sparseFields, ok := response["sparse_fields"].(map[string]interface{}); ok {
		rootFields := sparseFields[""].([]interface{})
		if len(rootFields) != 2 {
			t.Errorf("Expected 2 fields, got %d", len(rootFields))
		}
	} else {
		t.Error("Expected sparse_fields in response")
	}
}

func createTestRouterWithSparseFields[T any](t *testing.T, enabled bool) *ApiRouter[Test] {
	t.Helper()

	router := &ApiRouter[Test]{
		Configuration: map[configuration.ConfigurationType]configuration.Configuration{
			configuration.SparseFieldsEnabledType: {
				Values: []string{boolToString(enabled)},
			},
			configuration.SparseFieldsParameterNameType: {
				Values: []string{"fields"},
			},
		},
	}

	return router
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
