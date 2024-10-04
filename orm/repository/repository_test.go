package repository_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/philiphil/apiman/orm/entity"
	"github.com/philiphil/apiman/orm/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Product is a domain entity
type Product struct {
	ID          uint
	Name        string
	Weight      uint
	IsAvailable bool
}

func (p Product) SetId(a any) entity.IEntity {
	//TODO implement me
	panic("implement me")
}

func (p Product) GetId() entity.ID {
	//TODO implement me
	panic("implement me")
}

// ProductGorm is DTO used to map Product entity to database
type ProductGorm struct {
	ID          uint   `orm:"primaryKey;column:id"`
	Name        string `orm:"column:name"`
	Weight      uint   `orm:"column:weight"`
	IsAvailable bool   `orm:"column:is_available"`
}

// ToEntity respects the orm.GormModel interface
// Creates new Entity from GORM model.
func (g ProductGorm) ToEntity() Product {
	return Product{
		ID:          g.ID,
		Name:        g.Name,
		Weight:      g.Weight,
		IsAvailable: g.IsAvailable,
	}
}

// FromEntity respects the orm.GormModel interface
// Creates new GORM model from Entity.
func (g ProductGorm) FromEntity(product Product) interface{} {
	return ProductGorm{
		ID:          product.ID,
		Name:        product.Name,
		Weight:      product.Weight,
		IsAvailable: product.IsAvailable,
	}
}

func getDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("file:test?mode=memory&cache=shared&_fk=1"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}
func TestMain(m *testing.M) {
	db, err := getDB()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(ProductGorm{})
	if err != nil {
		log.Fatal(err)
	}
	ret := m.Run()
	os.Exit(ret)
}
func TestGormRepository_Insert(t *testing.T) {
	db, _ := getDB()
	repository := repository.NewRepository[ProductGorm, Product](db)
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
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}
	err = repository.Insert(ctx, &product2)
	if err != nil {
		t.Error(err)
	}

	product3 := Product{
		ID:          8,
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}
	err = repository.Insert(ctx, &product3)
	if err == nil {
		t.Error("should not insert")
	}

}

func TestGormRepository_FindByID(t *testing.T) {
	db, _ := getDB()
	repository := repository.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	_, err := repository.FindByID(ctx, 8)

	if err != nil {
		t.Error(err)
	}
}

func TestGormRepository_Count(t *testing.T) {
	db, _ := getDB()
	repository := repository.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	nb, err := repository.Count(ctx)

	if err != nil {
		t.Error(err)
	}
	if nb != 2 {
		t.Error("not good count")
	}
}

func TestGormRepository_Update(t *testing.T) {
	db, _ := getDB()
	repository := repository.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	product, err := repository.FindByID(ctx, 8)
	if err != nil {
		t.Error(err)
	}
	product.Name = "product2"
	err = repository.Update(ctx, &product)
	if err != nil {
		t.Error(err)
	}
	product2, err := repository.FindByID(ctx, 8)
	if err != nil {
		t.Error(err)
	}
	if product2.Name != "product2" {
		t.Error("not updated")
	}
}

func TestGormRepository_DeleteByID(t *testing.T) {
	db, _ := getDB()
	repository := repository.NewRepository[ProductGorm, Product](db)
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

func TestGormRepository_Find(t *testing.T) {
	db, _ := getDB()
	newRepository := repository.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	product := Product{
		ID:          1,
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}
	newRepository.Insert(ctx, &product)
	product2 := Product{
		ID:          2,
		Name:        "product2",
		Weight:      50,
		IsAvailable: true,
	}
	newRepository.Insert(ctx, &product2)
	many, err := newRepository.Find(ctx, repository.GreaterOrEqual("weight", 50))
	if err != nil {
		t.Error(err)
	}
	if len(many) != 3 {
		t.Error(len(many))
	}

	newRepository.Insert(ctx, &Product{
		ID:          3,
		Name:        "product3",
		Weight:      250,
		IsAvailable: false,
	})

	many, err = newRepository.Find(ctx, repository.GreaterOrEqual("weight", 90))
	if err != nil {
		t.Error(err)
	}
	if len(many) != 3 {
		t.Error("should be 3")
	}

	many, err = newRepository.Find(ctx, repository.And(
		repository.GreaterOrEqual("weight", 90),
		repository.Equal("is_available", true)),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 2 {
		t.Error("should be 1")
	}
}

func TestGormRepository_NewInstanceEntity(t *testing.T) {
	db, _ := getDB()
	newRepository := repository.NewRepository[ProductGorm, Product](db)

	entity := Product{}
	if newRepository.NewEntity() != entity {
		t.Error("should be empty entity")
	}
}

func TestGormRepository_Delete(t *testing.T) {
	db, _ := getDB()
	repository := repository.NewRepository[ProductGorm, Product](db)
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
	err = repository.Delete(ctx, &product)
	if err != nil {
		t.Error(err)
	}
}

func TestGormRepository_GetDB(t *testing.T) {
	db, _ := getDB()
	newRepository := repository.NewRepository[ProductGorm, Product](db)
	if newRepository.GetDB() != db {
		t.Error("should be equal")
	}
}

/*
TODO
Find (with sql cond)
*/
