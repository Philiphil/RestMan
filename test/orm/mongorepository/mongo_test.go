package mongorepository_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/mongorepository"
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
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
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
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
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
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
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
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)

	nb, err := repository.Count()
	if err != nil {
		t.Error(err)
	}
	if nb != 2 {
		t.Errorf("expected 2, got %d", nb)
	}
}

func TestMongoRepository_Update(t *testing.T) {
	collection := getCollection()
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
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
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
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
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
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

	many, err := repository.Find(ctx, mongorepository.GreaterOrEqual("weight", 50))
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

	many, err = repository.Find(ctx, mongorepository.GreaterOrEqual("weight", 90))
	if err != nil {
		t.Error(err)
	}
	if len(many) != 2 {
		t.Errorf("expected 2, got %d", len(many))
	}

	many, err = repository.Find(ctx, mongorepository.And(
		mongorepository.GreaterOrEqual("weight", 90),
		mongorepository.Equal("is_available", true)),
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
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)

	entity := Product{}
	if repository.NewEntity() != entity {
		t.Error("should be empty entity")
	}
}

func TestMongoRepository_Delete(t *testing.T) {
	collection := getCollection()
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
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

	err = repository.Delete([]*Product{&product})
	if err != nil {
		t.Error(err)
	}
}

func TestMongoRepository_FindAll(t *testing.T) {
	collection := getCollection()
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
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
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
	if repository.GetCollection() != collection {
		t.Error("should be equal")
	}
}

func TestMongoRepository_Conditions(t *testing.T) {
	collection := getCollection()
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	many, err := repository.Find(ctx, mongorepository.And(
		mongorepository.GreaterThan("weight", 100),
		mongorepository.Equal("is_available", false)),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Errorf("expected 1, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		mongorepository.Not(mongorepository.Equal("is_available", true)),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Errorf("expected 1, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		mongorepository.LessThan("weight", 51),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Errorf("expected 1, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		mongorepository.LessOrEqual("weight", 50),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Errorf("expected 1, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		mongorepository.In("weight", []uint{50, 100}),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 2 {
		t.Errorf("expected 2, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		mongorepository.Or(
			mongorepository.Equal("weight", 50),
			mongorepository.Equal("weight", 100),
		),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 2 {
		t.Errorf("expected 2, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		mongorepository.IsNull("qtt"),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 2 {
		t.Errorf("expected 2, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		mongorepository.IsNotNull("qtt"),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Errorf("expected 1, got %d", len(many))
	}

	many, err = repository.Find(ctx,
		mongorepository.Not(mongorepository.IsNull("qtt")),
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
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	many, err := repository.FindAll(ctx, mongorepository.OrderBy("weight", "DESC"))
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
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
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
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
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
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
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
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
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
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	many, err := repository.Find(ctx, mongorepository.Like("name", "product[12]"))
	if err != nil {
		t.Error(err)
	}
	if len(many) < 2 {
		t.Errorf("expected at least 2, got %d", len(many))
	}
}

func TestMongoRepository_Ilike(t *testing.T) {
	collection := getCollection()
	repository := mongorepository.NewRepository[ProductMongo, Product](collection)
	ctx := context.Background()

	many, err := repository.Find(ctx, mongorepository.Ilike("name", "PRODUCT[12]"))
	if err != nil {
		t.Error(err)
	}
	if len(many) < 2 {
		t.Errorf("expected at least 2, got %d", len(many))
	}
}
