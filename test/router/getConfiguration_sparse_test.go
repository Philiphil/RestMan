package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/philiphil/restman/router"
)

func TestParseDotNotationFields_Simple(t *testing.T) {
	router := &ApiRouter[Test]{}
	request := &SparseFieldRequest{
		ResourceFields: make(map[string][]string),
	}

	router.ParseDotNotationFields("id,name,email", request)

	// Check root fields
	rootFields := request.ResourceFields[""]
	if len(rootFields) != 3 {
		t.Errorf("Expected 3 root fields, got %d", len(rootFields))
	}

	expected := []string{"id", "name", "email"}
	for _, exp := range expected {
		found := false
		for _, field := range rootFields {
			if field == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected field '%s' not found in root fields", exp)
		}
	}
}

func TestParseDotNotationFields_Nested(t *testing.T) {
	router := &ApiRouter[Test]{}
	request := &SparseFieldRequest{
		ResourceFields: make(map[string][]string),
	}

	router.ParseDotNotationFields("id,name,posts.title,posts.content", request)

	// Check root fields
	rootFields := request.ResourceFields[""]
	if len(rootFields) != 2 {
		t.Errorf("Expected 2 root fields, got %d", len(rootFields))
	}

	// Check nested posts fields
	postsFields := request.ResourceFields["posts"]
	if len(postsFields) != 2 {
		t.Errorf("Expected 2 posts fields, got %d", len(postsFields))
	}

	expectedPosts := []string{"title", "content"}
	for _, exp := range expectedPosts {
		found := false
		for _, field := range postsFields {
			if field == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected posts field '%s' not found", exp)
		}
	}
}

func TestParseDotNotationFields_DeepNested(t *testing.T) {
	router := &ApiRouter[Test]{}
	request := &SparseFieldRequest{
		ResourceFields: make(map[string][]string),
	}

	router.ParseDotNotationFields("id,posts.category.name,posts.title", request)

	// Check root fields
	rootFields := request.ResourceFields[""]
	if len(rootFields) != 1 || rootFields[0] != "id" {
		t.Errorf("Expected root field 'id', got %v", rootFields)
	}

	// Check posts fields
	postsFields := request.ResourceFields["posts"]
	if len(postsFields) != 1 || postsFields[0] != "title" {
		t.Errorf("Expected posts field 'title', got %v", postsFields)
	}

	// Check posts.category fields
	categoryFields := request.ResourceFields["posts.category"]
	if len(categoryFields) != 1 || categoryFields[0] != "name" {
		t.Errorf("Expected posts.category field 'name', got %v", categoryFields)
	}
}

func TestParseDotNotationFields_EmptyAndWhitespace(t *testing.T) {
	router := &ApiRouter[Test]{}
	request := &SparseFieldRequest{
		ResourceFields: make(map[string][]string),
	}

	router.ParseDotNotationFields("id, name , , posts.title ", request)

	// Check root fields (should trim whitespace and skip empty)
	rootFields := request.ResourceFields[""]
	if len(rootFields) != 2 {
		t.Errorf("Expected 2 root fields after trimming, got %d: %v", len(rootFields), rootFields)
	}

	// Check posts fields
	postsFields := request.ResourceFields["posts"]
	if len(postsFields) != 1 || postsFields[0] != "title" {
		t.Errorf("Expected posts field 'title' after trimming, got %v", postsFields)
	}
}

func TestParseSeparateParamFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate query params: ?fields=id&fields.posts=title,content&fields.category=name
	req, _ := http.NewRequest("GET", "/test?fields=id&fields.posts=title,content&fields.category=name", nil)
	c.Request = req

	router := &ApiRouter[Test]{}
	request := &SparseFieldRequest{
		ResourceFields: make(map[string][]string),
	}

	router.ParseSeparateParamFields(c, "fields", request)

	// Check posts fields
	postsFields := request.ResourceFields["posts"]
	if len(postsFields) != 2 {
		t.Errorf("Expected 2 posts fields, got %d: %v", len(postsFields), postsFields)
	}

	// Check category fields
	categoryFields := request.ResourceFields["category"]
	if len(categoryFields) != 1 || categoryFields[0] != "name" {
		t.Errorf("Expected category field 'name', got %v", categoryFields)
	}
}

func TestParseJsonApiFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate query params: ?fields[users]=id,name&fields[posts]=title
	req, _ := http.NewRequest("GET", "/test?fields[users]=id,name&fields[posts]=title", nil)
	c.Request = req

	router := &ApiRouter[Test]{}
	request := &SparseFieldRequest{
		ResourceFields: make(map[string][]string),
	}

	router.ParseJsonApiFields(c, request)

	// Check users fields
	usersFields := request.ResourceFields["users"]
	if len(usersFields) != 2 {
		t.Errorf("Expected 2 users fields, got %d: %v", len(usersFields), usersFields)
	}

	expectedUsers := []string{"id", "name"}
	for _, exp := range expectedUsers {
		found := false
		for _, field := range usersFields {
			if field == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected users field '%s' not found", exp)
		}
	}

	// Check posts fields
	postsFields := request.ResourceFields["posts"]
	if len(postsFields) != 1 || postsFields[0] != "title" {
		t.Errorf("Expected posts field 'title', got %v", postsFields)
	}
}

func TestParseJsonApiFields_EmptyResource(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate query params with empty resource: ?fields[]=id
	req, _ := http.NewRequest("GET", "/test?fields[]=id", nil)
	c.Request = req

	router := &ApiRouter[Test]{}
	request := &SparseFieldRequest{
		ResourceFields: make(map[string][]string),
	}

	router.ParseJsonApiFields(c, request)

	// Should be ignored (empty resource name)
	if len(request.ResourceFields) != 0 {
		t.Errorf("Expected empty ResourceFields for empty resource name, got %v", request.ResourceFields)
	}
}

func TestParseDotNotationFields_MultipleResourcesSameField(t *testing.T) {
	router := &ApiRouter[Test]{}
	request := &SparseFieldRequest{
		ResourceFields: make(map[string][]string),
	}

	router.ParseDotNotationFields("posts.title,comments.title", request)

	// Check posts fields
	postsFields := request.ResourceFields["posts"]
	if len(postsFields) != 1 || postsFields[0] != "title" {
		t.Errorf("Expected posts field 'title', got %v", postsFields)
	}

	// Check comments fields
	commentsFields := request.ResourceFields["comments"]
	if len(commentsFields) != 1 || commentsFields[0] != "title" {
		t.Errorf("Expected comments field 'title', got %v", commentsFields)
	}
}

func TestParseDotNotationFields_DuplicateFields(t *testing.T) {
	router := &ApiRouter[Test]{}
	request := &SparseFieldRequest{
		ResourceFields: make(map[string][]string),
	}

	router.ParseDotNotationFields("id,name,id,posts.title,posts.title", request)

	// Check root fields (implementation does not deduplicate)
	rootFields := request.ResourceFields[""]
	if len(rootFields) != 3 { // id, name, id - but wait, it should be 3 because "id" appears twice
		t.Errorf("Expected 3 root fields, got %d", len(rootFields))
	}

	// Check posts fields
	postsFields := request.ResourceFields["posts"]
	if len(postsFields) != 2 { // title appears twice
		t.Errorf("Expected 2 posts fields, got %d", len(postsFields))
	}
}
