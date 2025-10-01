package gormrepository_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository"
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
	Qtt         *int
}

func (p Product) SetId(a any) entity.Entity {
	//TODO implement me
	panic("implement me")
}

func (p Product) GetId() entity.ID {
	return entity.ID(p.ID)
}

// ProductGorm is DTO used to map Product entity to database
type ProductGorm struct {
	ID          uint   `orm:"primaryKey;column:id"`
	Name        string `orm:"column:name"`
	Weight      uint   `orm:"column:weight"`
	IsAvailable bool   `orm:"column:is_available"`
	Qtt         *int   `orm:"column:qtt"`
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
func (g ProductGorm) FromEntity(product Product) any {
	return ProductGorm(product)
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
	repository := gormrepository.NewRepository[ProductGorm, Product](db)
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
	repository := gormrepository.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	_, err := repository.FindByID(ctx, 8)

	if err != nil {
		t.Error(err)
	}
}

func TestGormRepository_Count(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository.NewRepository[ProductGorm, Product](db)

	nb, err := repository.Count()

	if err != nil {
		t.Error(err)
	}
	if nb != 2 {
		t.Error("not good count")
	}
}

func TestGormRepository_Update(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	product, err := repository.FindByID(ctx, 8)
	if err != nil {
		t.Error(err)
	}
	product.Name = "product2"
	err = repository.Upsert(ctx, &product)
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
	repository := gormrepository.NewRepository[ProductGorm, Product](db)
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
	newRepository := gormrepository.NewRepository[ProductGorm, Product](db)
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
	many, err := newRepository.Find(ctx, gormrepository.GreaterOrEqual("weight", 50))
	if err != nil {
		t.Error(err)
	}
	if len(many) != 3 {
		t.Error(len(many))
	}

	i := 22
	newRepository.Insert(ctx, &Product{
		ID:          3,
		Name:        "product3",
		Weight:      250,
		IsAvailable: false,
		Qtt:         &i,
	})

	many, err = newRepository.Find(ctx, gormrepository.GreaterOrEqual("weight", 90))
	if err != nil {
		t.Error(err)
	}
	if len(many) != 3 {
		t.Error("should be 3")
	}

	many, err = newRepository.Find(ctx, gormrepository.And(
		gormrepository.GreaterOrEqual("weight", 90),
		gormrepository.Equal("is_available", true)),
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
	newRepository := gormrepository.NewRepository[ProductGorm, Product](db)

	entity := Product{}
	if newRepository.NewEntity() != entity {
		t.Error("should be empty entity")
	}
}

func TestGormRepository_Delete(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository.NewRepository[ProductGorm, Product](db)
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

func TestGormRepository_FindAll(t *testing.T) {
	db, _ := getDB()
	newRepository := gormrepository.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	many, err := newRepository.FindAll(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 4 {
		t.Error(len(many))
	}
}

func TestGormRepository_GetDB(t *testing.T) {
	db, _ := getDB()
	newRepository := gormrepository.NewRepository[ProductGorm, Product](db)
	if newRepository.GetDB() != db {
		t.Error("should be equal")
	}
}

// until now we test only the GreatOrEqual and Equal
// we will test the other conditions
func TestGormRepository_Conditions(t *testing.T) {
	db, _ := getDB()
	newRepository := gormrepository.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	many, err := newRepository.Find(ctx, gormrepository.And(
		gormrepository.GreaterThan("weight", 100),
		gormrepository.Equal("is_available", false)),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Error(len(many))
	}

	many, err = newRepository.Find(ctx,
		gormrepository.Not(gormrepository.Equal("is_available", true)),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Error(len(many))
	}

	many, err = newRepository.Find(ctx,
		gormrepository.LessThan("weight", 51),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Error(len(many))
	}

	many, err = newRepository.Find(ctx,
		gormrepository.LessOrEqual("weight", 51),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Error(many)
	}

	many, err = newRepository.Find(ctx,
		gormrepository.In("weight", []uint{50, 100}),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 3 {
		t.Error(many)
	}
	many, err = newRepository.Find(ctx,
		gormrepository.Or(
			gormrepository.Equal("weight", 50),
			gormrepository.Equal("weight", 100),
		),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 3 {
		t.Error(many)
	}

	many, err = newRepository.Find(ctx,
		gormrepository.Or(
			gormrepository.Equal("weight", 50),
			gormrepository.Equal("weight", 100),
		),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 3 {
		t.Error(many)
	}

	many, err = newRepository.Find(ctx,
		gormrepository.IsNull("qtt"),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 3 {
		t.Error(many)
	}
	many, err = newRepository.Find(ctx,
		gormrepository.IsNotNull("qtt"),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Error(many)
	}

	many, err = newRepository.Find(ctx,
		gormrepository.Not(gormrepository.IsNull("qtt")),
	)
	if err != nil {
		t.Error(err)
	}
	if len(many) != 1 {
		t.Error(many)
	}
}

func TestGormRepository_Order(t *testing.T) {
	db, _ := getDB()
	newRepository := gormrepository.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	many, err := newRepository.FindAll(ctx, gormrepository.OrderBy("weight", "DESC"))
	if err != nil {
		t.Error(err)
	}
	if len(many) == 0 {
		t.Error(many)
	}
	fmt.Println(many)
	if many[0].Weight != 250 {
		t.Error(many[0].Weight)
	}
}

func TestGormRepository_PreloadAssociation(t *testing.T) {
	db, _ := getDB()
	newRepository := gormrepository.NewRepository[ProductGorm, Product](db)
	newRepository.EnablePreloadAssociations()
	newRepository.DisablePreloadAssociations()
	newRepository.SetPreloadAssociations(true)
	newRepository.SetPreloadAssociations(false)

}

func TestGormRepository_BatchInsert(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	products := []*Product{
		{
			ID:          10,
			Name:        "product1",
			Weight:      100,
			IsAvailable: true,
		},
		{
			ID:          11,
			Name:        "product2",
			Weight:      50,
			IsAvailable: true,
		},
	}
	err := repository.BatchInsert(ctx, products)
	if err != nil {
		t.Error(err)
	}
}

// BatchUpdate
func TestGormRepository_BatchUpdate(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	products := []*Product{
		{
			ID:          11,
			Name:        "product1",
			Weight:      100,
			IsAvailable: true,
		},
		{
			ID:          12,
			Name:        "product2",
			Weight:      50,
			IsAvailable: true,
		},
	}
	err := repository.BatchUpdate(ctx, products)
	if err != nil {
		t.Error(err)
	}
}

// BatchDelete
func TestGormRepository_BatchDelete(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	products := []*Product{
		{
			ID:          11,
			Name:        "product1",
			Weight:      100,
			IsAvailable: true,
		},
		{
			ID:          12,
			Name:        "product2",
			Weight:      50,
			IsAvailable: true,
		},
	}
	err := repository.BatchDelete(ctx, products)
	if err != nil {
		t.Error(err)
	}
}
