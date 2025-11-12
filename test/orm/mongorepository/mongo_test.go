package mongorepository_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/philiphil/restman/orm/entity"
	. "github.com/philiphil/restman/orm/mongorepository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	ID          uint
	Name        string
	Weight      uint
	IsAvailable bool
	Qtt         *int
}

func (p Product) SetId(a any) entity.Entity {
	switch v := a.(type) {
	case uint:
		p.ID = v
	case int:
		p.ID = uint(v)
	case int64:
		p.ID = uint(v)
	case entity.ID:
		p.ID = uint(v)
	}
	return p
}

func (p Product) GetId() entity.ID {
	return entity.ID(p.ID)
}

type ProductMongo struct {
	ID          uint   `bson:"_id"`
	Name        string `bson:"name"`
	Weight      uint   `bson:"weight"`
	IsAvailable bool   `bson:"is_available"`
	Qtt         *int   `bson:"qtt,omitempty"`
}

func (m ProductMongo) ToEntity() Product {
	return Product{
		ID:          m.ID,
		Name:        m.Name,
		Weight:      m.Weight,
		IsAvailable: m.IsAvailable,
		Qtt:         m.Qtt,
	}
}

func (m ProductMongo) FromEntity(product Product) any {
	return ProductMongo{
		ID:          product.ID,
		Name:        product.Name,
		Weight:      product.Weight,
		IsAvailable: product.IsAvailable,
		Qtt:         product.Qtt,
	}
}

var testClient *mongo.Client
var testCollection *mongo.Collection

func getCollection() *mongo.Collection {
	if testCollection != nil {
		return testCollection
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB not available: ", err)
	}

	testClient = client
	testCollection = client.Database("restman_test").Collection("products")
	return testCollection
}

func TestMain(m *testing.M) {
	collection := getCollection()
	collection.Drop(context.Background())

	ret := m.Run()

	if testClient != nil {
		testClient.Disconnect(context.Background())
	}
	os.Exit(ret)
}

func TestMongoRepository_Insert(t *testing.T) {
	collection := getCollection()
	collection.Drop(context.Background())
	repository := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	product := Product{
		ID:          8,
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}
	err := repository.Insert(ctx, &product)
	if err != nil {
		t.Error(err)
	}

	product2 := Product{
		ID:          9,
		Name:        "product2",
		Weight:      100,
		IsAvailable: true,
	}
	err = repository.Insert(ctx, &product2)
	if err != nil {
		t.Error(err)
	}

	product3 := Product{
		ID:          8,
		Name:        "product3",
		Weight:      100,
		IsAvailable: true,
	}
	err = repository.Insert(ctx, &product3)
	if err == nil {
		t.Error("should not insert duplicate ID")
	}
}

func TestMongoRepository_FindByID(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	product, err := repository.FindByID(ctx, 8)
	if err != nil {
		t.Error(err)
	}
	if product.Name != "product1" {
		t.Error("wrong product")
	}
}

func TestMongoRepository_Count(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)

	nb, err := repository.Count(context.Background())
	if err != nil {
		t.Error(err)
	}
	if nb != 2 {
		t.Errorf("expected 2, got %d", nb)
	}
}

func TestMongoRepository_Upsert(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	product, err := repository.FindByID(ctx, 8)
	if err != nil {
		t.Error(err)
	}
	product.Name = "product_updated"
	err = repository.Upsert(ctx, &product)
	if err != nil {
		t.Error(err)
	}

	product2, err := repository.FindByID(ctx, 8)
	if err != nil {
		t.Error(err)
	}
	if product2.Name != "product_updated" {
		t.Error("not updated")
	}
}

func TestMongoRepository_DeleteByID(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	err := repository.DeleteByID(ctx, 8)
	if err != nil {
		t.Error(err)
	}

	_, err = repository.FindByID(ctx, 8)
	if err == nil {
		t.Error("supposed to be deleted")
	}
}

func TestMongoRepository_Find(t *testing.T) {
	collection := getCollection()
	collection.Drop(context.Background())
	repository := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	product := Product{
		ID:          1,
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}
	repository.Insert(ctx, &product)

	product2 := Product{
		ID:          2,
		Name:        "product2",
		Weight:      50,
		IsAvailable: true,
	}
	repository.Insert(ctx, &product2)

	many, err := repository.Find(ctx, GreaterOrEqual("weight", 50))
	if err != nil {
		t.Error(err)
	}
	if len(many) != 2 {
		t.Errorf("expected 2, got %d", len(many))
	}

	i := 22
	repository.Insert(ctx, &Product{
		ID:          3,
		Name:        "product3",
		Weight:      250,
		IsAvailable: false,
		Qtt:         &i,
	})

	many, err = repository.Find(ctx, GreaterOrEqual("weight", 90))
	if err != nil {
		t.Error(err)
	}
	if len(many) != 2 {
		t.Errorf("expected 2, got %d", len(many))
	}

	many, err = repository.Find(ctx, And(
		GreaterOrEqual("weight", 90),
		Equal("is_available", true)),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Errorf("expected 1, got %d", len(many))
	}
}

func TestMongoRepository_NewInstanceEntity(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)

	entity := Product{}
	if repository.NewEntity() != entity {
		t.Error("should be empty entity")
	}
}

func TestMongoRepository_Delete(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	product := Product{
		ID:          8,
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}
	err := repository.Insert(ctx, &product)
	if err != nil {
		t.Error(err)
	}

	err = repository.Delete(ctx, []*Product{&product})
	if err != nil {
		t.Error(err)
	}
}

func TestMongoRepository_FindAll(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	many, err := repository.FindAll(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 3 {
		t.Errorf("expected 3, got %d", len(many))
	}
}

func TestMongoRepository_GetCollection(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)
	if repository.GetCollection() != collection {
		t.Error("should be equal")
	}
}

func TestMongoRepository_Conditions(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	many, err := repository.Find(ctx, And(
		GreaterThan("weight", 100),
		Equal("is_available", false)),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Errorf("expected 1, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		Not(Equal("is_available", true)),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Errorf("expected 1, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		LessThan("weight", 51),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Errorf("expected 1, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		LessOrEqual("weight", 50),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Errorf("expected 1, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		In("weight", []uint{50, 100}),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 2 {
		t.Errorf("expected 2, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		Or(
			Equal("weight", 50),
			Equal("weight", 100),
		),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 2 {
		t.Errorf("expected 2, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		IsNull("qtt"),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 2 {
		t.Errorf("expected 2, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		IsNotNull("qtt"),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Errorf("expected 1, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		Not(IsNull("qtt")),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Errorf("expected 1, got %d", len(many))
	}
}

func TestMongoRepository_Order(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	many, err := repository.FindAll(ctx, OrderBy("weight", "DESC"))
	if err != nil {
		t.Error(err)
	}
	if len(many) == 0 {
		t.Error("no results")
	}
	fmt.Println(many)
	if many[0].Weight != 250 {
		t.Errorf("expected 250, got %d", many[0].Weight)
	}
}

func TestMongoRepository_BatchInsert(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	products := []*Product{
		{
			ID:          10,
			Name:        "product10",
			Weight:      100,
			IsAvailable: true,
		},
		{
			ID:          11,
			Name:        "product11",
			Weight:      50,
			IsAvailable: true,
		},
	}
	err := repository.BatchInsert(ctx, products)
	if err != nil {
		t.Error(err)
	}
}

func TestMongoRepository_BatchUpdate(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	products := []*Product{
		{
			ID:          10,
			Name:        "product10_updated",
			Weight:      150,
			IsAvailable: true,
		},
		{
			ID:          11,
			Name:        "product11_updated",
			Weight:      75,
			IsAvailable: true,
		},
	}
	err := repository.BatchUpdate(ctx, products)
	if err != nil {
		t.Error(err)
	}
}

func TestMongoRepository_BatchDelete(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	products := []*Product{
		{
			ID:          10,
			Name:        "product10",
			Weight:      100,
			IsAvailable: true,
		},
		{
			ID:          11,
			Name:        "product11",
			Weight:      50,
			IsAvailable: true,
		},
	}
	err := repository.BatchDelete(ctx, products)
	if err != nil {
		t.Error(err)
	}
}

func TestMongoRepository_FindByIDs(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	ids := []entity.ID{1, 2, 3}
	products, err := repository.FindByIDs(ctx, ids)
	if err != nil {
		t.Error(err)
	}
	if len(products) != 3 {
		t.Errorf("expected 3, got %d", len(products))
	}
}

func TestMongoRepository_Like(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	many, err := repository.Find(ctx, Like("name", "product[12]"))
	if err != nil {
		t.Error(err)
	}
	if len(many) < 2 {
		t.Errorf("expected at least 2, got %d", len(many))
	}
}

func TestMongoRepository_Ilike(t *testing.T) {
	collection := getCollection()
	repository := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	many, err := repository.Find(ctx, Ilike("name", "PRODUCT[12]"))
	if err != nil {
		t.Error(err)
	}
	if len(many) < 2 {
		t.Errorf("expected at least 2, got %d", len(many))
	}
}

func TestMongoRepository_Create(t *testing.T) {
	collection := getCollection()
	collection.Drop(context.Background())
	repo := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	product := Product{ID: 1, Name: "p1", Weight: 10, IsAvailable: true}
	if err := repo.Create(ctx, []*Product{&product}); err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	retrieved, err := repo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("Failed to retrieve entity: %v", err)
	}
	if retrieved.GetId() != product.GetId() {
		t.Fatalf("Retrieved entity ID %d does not match created entity ID %d", retrieved.GetId(), product.GetId())
	}
}

func TestMongoRepository_Read(t *testing.T) {
	collection := getCollection()
	collection.Drop(context.Background())
	repo := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	p1 := Product{ID: 3, Name: "p3", Weight: 30, IsAvailable: true}
	p2 := Product{ID: 4, Name: "p4", Weight: 40, IsAvailable: true}
	if err := repo.Create(ctx, []*Product{&p1, &p2}); err != nil {
		t.Fatalf("Failed to create entities: %v", err)
	}

	products, err := repo.Read(ctx, []entity.ID{3, 4})
	if err != nil {
		t.Fatalf("Failed to read entities: %v", err)
	}
	if len(products) != 2 {
		t.Fatalf("Expected to read 2 entities, got %d", len(products))
	}
}

func TestMongoRepository_List(t *testing.T) {
	collection := getCollection()
	collection.Drop(context.Background())
	repo := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	productsToCreate := []*Product{
		{ID: 6, Name: "p6", Weight: 60, IsAvailable: true},
		{ID: 7, Name: "p7", Weight: 70, IsAvailable: true},
		{ID: 8, Name: "p8", Weight: 80, IsAvailable: true},
	}
	if err := repo.Create(ctx, productsToCreate); err != nil {
		t.Fatalf("Failed to create entities: %v", err)
	}

	products, err := repo.List(ctx, 10, 0, nil)
	if err != nil {
		t.Fatalf("Failed to list entities: %v", err)
	}
	if len(products) != 3 {
		t.Fatalf("Expected to list 3 entities, got %d", len(products))
	}
}
func TestMongoRepository_CountWithSpecifications(t *testing.T) {
	collection := getCollection()
	collection.Drop(context.Background())
	repo := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	productsToCreate := []*Product{
		{ID: 12, Name: "p12", Weight: 120, IsAvailable: true},
		{ID: 13, Name: "p13", Weight: 130, IsAvailable: false},
		{ID: 14, Name: "p14", Weight: 140, IsAvailable: true},
	}
	if err := repo.Create(ctx, productsToCreate); err != nil {
		t.Fatalf("Failed to create entities: %v", err)
	}

	count, err := repo.CountWithSpecifications(ctx, Equal("is_available", true))
	if err != nil {
		t.Fatalf("Failed to count entities: %v", err)
	}
	if count != 2 {
		t.Fatalf("Expected to count 2 entities, got %d", count)
	}
}
func TestMongoRepository_Update(t *testing.T) {
	collection := getCollection()
	collection.Drop(context.Background())
	repo := NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	product := Product{ID: 20, Name: "p20", Weight: 200, IsAvailable: true}
	if err := repo.Create(ctx, []*Product{&product}); err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	product.Name = "p20_updated"
	if err := repo.Update(ctx, []*Product{&product}); err != nil {
		t.Fatalf("Failed to update entity: %v", err)
	}

	updatedProduct, err := repo.FindByID(ctx, 20)
	if err != nil {
		t.Fatalf("Failed to retrieve updated entity: %v", err)
	}
	if updatedProduct.Name != "p20_updated" {
		t.Fatalf("Expected updated entity name to be 'p20_updated', got '%s'", updatedProduct.Name)
	}
}
