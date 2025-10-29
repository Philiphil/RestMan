package example_test

import (
	"bytes"
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

//Serialization groups are an easy way to control which fields are included in API responses and requests based on context.
// In this example, we define "read" and "write" groups to manage field visibility during GET and POST/PUT operations respectively.

// With SerializationGroups configuration , as with any configuration we can set router-wide defaults that can be overridden on a operation level
// we also provide client control over group overwriting via query parameters.

type SerialProduct struct {
	entity.BaseEntity
	Name        string  `json:"name" groups:"read,write"`
	Price       float64 `json:"price" groups:"write"`
	InternalSKU string  `json:"internal_sku" groups:"read"`        //should be never editable
	CostPrice   float64 `json:"cost_price" groups:"private_stuff"` //should be never visible via api
}

func (p SerialProduct) GetId() entity.ID { return p.Id }
func (p SerialProduct) SetId(id any) entity.Entity {
	p.Id = entity.CastId(id)
	return p
}
func (p SerialProduct) ToEntity() SerialProduct        { return p }
func (p SerialProduct) FromEntity(e SerialProduct) any { return e }

func getSerialProductDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:serial_product_test?mode=memory&cache=shared&_fk=1"), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Silent),
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic(err)
	}
	return db
}

func TestSerializationGroupsReadWrite(t *testing.T) {
	db := getSerialProductDB()
	db.AutoMigrate(&SerialProduct{})

	r := gin.New()
	r.Use(gin.Recovery())

	routes := route.DefaultApiRoutes()
	routes[route.Post].Configuration[configuration.InputSerializationGroupsType] = configuration.InputSerializationGroups("write")

	productRouter := router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[SerialProduct](db)),
		routes,
		configuration.OutputSerializationGroups("read"), //default router-wide serialization group should be "read"
	)

	productRouter.AllowRoutes(r)

	postData := SerialProduct{
		Name:        "Laptop",
		Price:       999.99,
		CostPrice:   700.00,
		InternalSKU: "SKU-IGNORED",
	}
	jsonData, _ := json.Marshal(postData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/serial_product", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var createdProduct SerialProduct
	db.First(&createdProduct, 1)

	if createdProduct.Name != "Laptop" {
		t.Errorf("Name should have been written, got: %s", createdProduct.Name)
	}

	if createdProduct.Price != 999.99 {
		t.Errorf("Price should have been written to DB, got: %f", createdProduct.Price)
	}

	if createdProduct.InternalSKU == "SKU-IGNORED" {
		t.Error("InternalSKU should not have been written (not in write group)")
	}

	if createdProduct.CostPrice == 700.00 {
		t.Error("CostPrice should not have been written (not in write group)")
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/serial_product/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response SerialProduct
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Name != "Laptop" {
		t.Error("Name should be visible in GET response (read group)")
	}

	if response.Price != 0 {
		t.Errorf("Price should NOT be visible in GET response (write group only), got: %f", response.Price)
	}

	if response.InternalSKU != "" {
		t.Errorf("InternalSKU should be empty in GET (read group but never written), got: %s", response.InternalSKU)
	}

	if response.CostPrice != 0 {
		t.Errorf("CostPrice should be zero in GET (private_stuff group not included), got: %f", response.CostPrice)
	}
}
