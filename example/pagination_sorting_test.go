package example_test

import (
	"encoding/json"
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

type Article struct {
	entity.BaseEntity
	Title       string `json:"title" groups:"read,write"`
	Content     string `json:"content" groups:"read,write"`
	PublishedAt string `json:"published_at" groups:"read,write"`
	Views       int    `json:"views" groups:"read"`
}

func (a Article) GetId() entity.ID { return a.Id }
func (a Article) SetId(id any) entity.Entity {
	a.Id = entity.CastId(id)
	return a
}
func (a Article) ToEntity() Article       { return a }
func (a Article) FromEntity(e Article) any { return e }

func getPaginationDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:pagination_test?mode=memory&cache=shared&_fk=1"), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Silent),
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic(err)
	}
	return db
}

func TestPaginationAndSorting(t *testing.T) {
	db := getPaginationDB()
	db.AutoMigrate(&Article{})

	r := gin.New()
	r.Use(gin.Recovery())

	articleRouter := router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Article](db)),
		route.DefaultApiRoutes(),
		configuration.ItemPerPage(5),
		configuration.MaxItemPerPage(20),
		configuration.PaginationClientControl(true),
		configuration.SortingClientControl(true),
		configuration.Sorting(map[string]string{"published_at": "desc"}),
	)

	articleRouter.AllowRoutes(r)

	for i := 1; i <= 15; i++ {
		article := Article{
			Title:       "Article " + string(rune(i+'0')),
			Content:     "Content " + string(rune(i+'0')),
			PublishedAt: "2024-01-" + string(rune(i+'0')),
			Views:       i * 100,
		}
		db.Create(&article)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/article?page=1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var results []Article
	json.Unmarshal(w.Body.Bytes(), &results)

	if len(results) != 5 {
		t.Errorf("Expected 5 articles per page, got %d", len(results))
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/article?page=1&itemsPerPage=10", nil)
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &results)

	if len(results) != 10 {
		t.Errorf("Expected 10 articles with custom page size, got %d", len(results))
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/article?order[title]=asc", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Sorting request failed: expected 200, got %d", w.Code)
	}
}
