package filter_test

import (
	"encoding/json"
	"testing"

	"github.com/philiphil/restman/format"
	"github.com/philiphil/restman/serializer"
)

type User struct {
	ID       int    `json:"id" groups:"test"`
	Name     string `json:"name" groups:"test"`
	Email    string `json:"email" groups:"test"`
	Password string `json:"password" groups:"test"`
}

type Post struct {
	ID      int    `json:"id" groups:"test"`
	Title   string `json:"title" groups:"test"`
	Content string `json:"content" groups:"test"`
	Author  User   `json:"author" groups:"test"`
}

func TestSparseFields_RootOnly(t *testing.T) {
	user := User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret",
	}

	s := serializer.NewSerializer(format.JSON)
	sparseFields := map[string][]string{
		"": {"id", "name"},
	}

	jsonStr, err := s.Serialize(user, serializer.SerializeOptions{Groups: []string{"test"}, SparseFields: sparseFields})
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 fields in result, got %d: %v", len(result), result)
	}

	if _, ok := result["id"]; !ok {
		t.Error("Expected 'id' field in result")
	}

	if _, ok := result["name"]; !ok {
		t.Error("Expected 'name' field in result")
	}

	if _, ok := result["email"]; ok {
		t.Error("Did not expect 'email' field in result")
	}

	if _, ok := result["password"]; ok {
		t.Error("Did not expect 'password' field in result")
	}
}

func TestSparseFields_Nested(t *testing.T) {
	post := Post{
		ID:      1,
		Title:   "Test Post",
		Content: "This is test content",
		Author: User{
			ID:       1,
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "secret",
		},
	}

	s := serializer.NewSerializer(format.JSON)
	sparseFields := map[string][]string{
		"":       {"id", "title", "author"},
		"Author": {"id", "name"},
	}

	jsonStr, err := s.Serialize(post, serializer.SerializeOptions{Groups: []string{"test"}, SparseFields: sparseFields})
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 fields in root result, got %d: %v", len(result), result)
	}

	if _, ok := result["id"]; !ok {
		t.Error("Expected 'id' field in result")
	}

	if _, ok := result["title"]; !ok {
		t.Error("Expected 'title' field in result")
	}

	if _, ok := result["content"]; ok {
		t.Error("Did not expect 'content' field in result")
	}

	author, ok := result["author"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'author' to be an object")
	}

	if len(author) != 2 {
		t.Errorf("Expected 2 fields in author, got %d: %v", len(author), author)
	}

	if _, ok := author["id"]; !ok {
		t.Error("Expected 'id' field in author")
	}

	if _, ok := author["name"]; !ok {
		t.Error("Expected 'name' field in author")
	}

	if _, ok := author["email"]; ok {
		t.Error("Did not expect 'email' field in author")
	}

	if _, ok := author["password"]; ok {
		t.Error("Did not expect 'password' field in author")
	}
}

func TestSparseFields_Slice(t *testing.T) {
	users := []User{
		{ID: 1, Name: "John", Email: "john@example.com", Password: "secret1"},
		{ID: 2, Name: "Jane", Email: "jane@example.com", Password: "secret2"},
	}

	s := serializer.NewSerializer(format.JSON)
	sparseFields := map[string][]string{
		"": {"id", "name"},
	}

	jsonStr, err := s.Serialize(users, serializer.SerializeOptions{Groups: []string{"test"}, SparseFields: sparseFields})
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 users in result, got %d", len(result))
	}

	for i, user := range result {
		if len(user) != 2 {
			t.Errorf("Expected 2 fields in user %d, got %d: %v", i, len(user), user)
		}

		if _, ok := user["id"]; !ok {
			t.Errorf("Expected 'id' field in user %d", i)
		}

		if _, ok := user["name"]; !ok {
			t.Errorf("Expected 'name' field in user %d", i)
		}

		if _, ok := user["email"]; ok {
			t.Errorf("Did not expect 'email' field in user %d", i)
		}
	}
}

func TestSparseFields_None(t *testing.T) {
	user := User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret",
	}

	s := serializer.NewSerializer(format.JSON)

	jsonStr, err := s.Serialize(user, serializer.SerializeOptions{Groups: []string{"test"}})
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(result) != 4 {
		t.Errorf("Expected 4 fields when no sparse fields provided, got %d: %v", len(result), result)
	}
}
