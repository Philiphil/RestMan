package example_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository"
	"github.com/philiphil/restman/route"
	"github.com/philiphil/restman/router"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Author struct {
	entity.BaseEntity
	Name  string `json:"name" groups:"read,write"`
	Email string `json:"email" groups:"read,write"`
	Books []Book `json:"books,omitempty" groups:"read" gorm:"foreignKey:AuthorID"`
}

func (a Author) GetId() entity.ID { return a.Id }
func (a Author) SetId(id any) entity.Entity {
	a.Id = entity.CastId(id)
	return a
}
func (a Author) ToEntity() Author       { return a }
func (a Author) FromEntity(e Author) any { return e }

type Book struct {
	entity.BaseEntity
	Title    string `json:"title" groups:"read,write"`
	ISBN     string `json:"isbn" groups:"read,write"`
	AuthorID uint   `json:"author_id" groups:"write"`
}

func (b Book) GetId() entity.ID { return b.Id }
func (b Book) SetId(id any) entity.Entity {
	b.Id = entity.CastId(id)
	return b
}
func (b Book) ToEntity() Book       { return b }
func (b Book) FromEntity(e Book) any { return e }

func getSubresourceDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:subresource_test?mode=memory&cache=shared&_fk=1"), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Silent),
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic(err)
	}
	return db
}

func TestSubresources(t *testing.T) {
	db := getSubresourceDB()
	db.AutoMigrate(&Author{}, &Book{})

	r := gin.New()
	r.Use(gin.Recovery())

	authorRouter := router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Author](db)),
		route.DefaultApiRoutes(),
	)

	bookRouter := router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Book](db)),
		route.DefaultApiRoutes(),
	)

	authorRouter.AddSubresource(bookRouter)
	authorRouter.AllowRoutes(r)

	author := Author{Name: "J.K. Rowling", Email: "jk@example.com"}
	db.Create(&author)

	book := Book{Title: "Harry Potter", ISBN: "123-456", AuthorID: uint(author.Id)}
	db.Create(&book)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/author/1/book/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
