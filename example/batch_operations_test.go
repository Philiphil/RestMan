package example_test

import (
	"bytes"
	"encoding/json"
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

type Product struct {
	entity.BaseEntity
	Name  string  `json:"name" groups:"read,write"`
	Price float64 `json:"price" groups:"read,write"`
	Stock int     `json:"stock" groups:"read,write"`
}

func (p Product) GetId() entity.ID { return p.Id }
func (p Product) SetId(id any) entity.Entity {
	p.Id = entity.CastId(id)
	return p
}
func (p Product) ToEntity() Product       { return p }
func (p Product) FromEntity(e Product) any { return e }

func getBatchDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:batch_test?mode=memory&cache=shared&_fk=1"), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Silent),
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic(err)
	}
	return db
}

func TestBatchOperations(t *testing.T) {
	db := getBatchDB()
	db.AutoMigrate(&Product{})

	r := gin.New()
	r.Use(gin.Recovery())

	productRouter := router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Product](db)),
		route.DefaultApiRoutes(),
	)

	productRouter.AllowRoutes(r)

	products := []Product{
		{Name: "Product 1", Price: 10.99, Stock: 100},
		{Name: "Product 2", Price: 20.99, Stock: 50},
		{Name: "Product 3", Price: 30.99, Stock: 25},
	}

	body, _ := json.Marshal(products)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/product/batch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Batch POST failed: expected 201, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/product/batch?ids=1,2,3", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Batch GET failed: expected 200, got %d", w.Code)
	}

	var results []Product
	json.Unmarshal(w.Body.Bytes(), &results)

	if len(results) != 3 {
		t.Errorf("Expected 3 products, got %d", len(results))
	}
}
