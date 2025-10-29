package example_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository"
	"github.com/philiphil/restman/route"
	"github.com/philiphil/restman/router"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	entity.BaseEntity
	Email         string `json:"email" groups:"read,write,public"`
	Password      string `json:"password" groups:"write"`
	Token         string `json:"token" groups:"admin"`
	Name          string `json:"name" groups:"read,write,public"`
	InternalNotes string `json:"internal_notes" groups:"admin"`
}

func (u User) GetId() entity.ID { return u.Id }
func (u User) SetId(id any) entity.Entity {
	u.Id = entity.CastId(id)
	return u
}
func (u User) ToEntity() User        { return u }
func (u User) FromEntity(e User) any { return e }

func getSerializationDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:serialization_test?mode=memory&cache=shared&_fk=1"), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Silent),
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic(err)
	}
	return db
}

func TestCustomSerialization(t *testing.T) {
	db := getSerializationDB()
	db.AutoMigrate(&User{})

	r := gin.New()
	r.Use(gin.Recovery())

	routes := route.DefaultApiRoutes()

	userRouter := router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[User](db)),
		routes,
		configuration.InputSerializationGroups("read", "public"),
	)

	userRouter.AllowRoutes(r)

	user := User{
		Email:         "user@example.com",
		Password:      "secret123",
		Token:         "admin-token-xyz",
		Name:          "John Doe",
		InternalNotes: "This is an internal note",
	}
	db.Create(&user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	response := w.Body.String()

	if contains(response, "password") || contains(response, "secret123") {
		t.Error("Password should not be visible in GET response")
	}

	if contains(response, "token") || contains(response, "admin-token-xyz") {
		t.Error("Token should not be visible with public groups")
	}

	if contains(response, "internal_notes") {
		t.Error("Internal notes should not be visible with public groups")
	}

	if !contains(response, "email") || !contains(response, "name") {
		t.Error("Email and name should be visible with public groups")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && hasSubstring(s, substr))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
