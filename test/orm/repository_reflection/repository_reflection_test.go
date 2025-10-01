package gormrepository_reflection_test

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository_reflection"
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
	panic("implement me")
}

func (p Product) GetId() entity.ID {
	return entity.ID(p.ID)
}

// ProductGorm is DTO used to map Product entity to database
type ProductGorm struct {
	ID          uint   `gorm:"primaryKey;column:id"` // Changed orm tag to gorm
	Name        string `gorm:"column:name"`
	Weight      uint   `gorm:"column:weight"`
	IsAvailable bool   `gorm:"column:is_available"`
	Qtt         *int   `gorm:"column:qtt"`
}

// ToEntity respects the orm.GormModel interface
// Creates new Entity from GORM model.
func (g ProductGorm) ToEntity() entity.Entity { // Return entity.Entity
	return Product(g) // Convert ProductGorm to Product
}

// FromEntity respects the orm.GormModel interface
// Creates new GORM model from Entity.
// It should modify the receiver g and return it as any, as per the interface error.
func (g *ProductGorm) FromEntity(e entity.Entity) any { // entity.Entity as input, pointer receiver
	p, ok := e.(*Product) // Changed to *Product
	if !ok {
		// If e is Product (not *Product), then try to get a pointer to it.
		// This might happen if the entity is passed by value somewhere.
		if val, okVal := e.(Product); okVal {
			p = &val
			ok = true
		} else {
			log.Printf("Error: FromEntity received type %T but expected Product or *Product", e)
			// Handle error appropriately, perhaps return nil or panic depending on desired behavior
			return nil
		}
	}
	g.ID = p.ID
	g.Name = p.Name
	g.Weight = p.Weight
	g.IsAvailable = p.IsAvailable
	g.Qtt = p.Qtt
	return g // Return the modified receiver as any
}

func getDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("file:test_reflection?mode=memory&cache=shared&_fk=1"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

var (
	productModelType  = reflect.TypeOf(ProductGorm{})
	productEntityType = reflect.TypeOf(Product{})
)

func TestMain(m *testing.M) {
	db, err := getDB()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&ProductGorm{}) // Pass pointer for AutoMigrate
	if err != nil {
		log.Fatal(err)
	}
	ret := m.Run()
	os.Exit(ret)
}

func TestGormRepository_Insert(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)
	ctx := context.Background()

	product := Product{
		ID:          8,
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}
	// Pass pointer to entity for Insert
	err := repository.Insert(ctx, &product)
	if err != nil {
		t.Errorf("Insert failed: %v", err)
	}
	if product.ID == 0 { // GORM typically populates the ID if it's auto-incrementing and was zero
		t.Errorf("Insert did not populate ID, got %d", product.ID)
	}

	product2 := Product{
		// ID: 0, // Auto-increment
		Name:        "product2",
		Weight:      100,
		IsAvailable: true,
	}
	err = repository.Insert(ctx, &product2)
	if err != nil {
		t.Errorf("Insert failed for product2: %v", err)
	}
	if product2.ID == 0 {
		t.Errorf("Insert did not populate ID for product2, got %d", product2.ID)
	}
	if product2.ID == product.ID {
		t.Errorf("product2 ID should be different from product1 ID, got %d for both", product.ID)
	}

	product3 := Product{
		ID:          product.ID, // Use existing ID
		Name:        "product1_conflict",
		Weight:      100,
		IsAvailable: true,
	}
	err = repository.Insert(ctx, &product3)
	if err == nil {
		t.Error("Insert should fail due to ID conflict, but it succeeded")
	}
}

func TestGormRepository_FindByID(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)
	ctx := context.Background()

	// First, insert a product to find
	productToInsert := Product{
		ID:          100, // Use a distinct ID for this test
		Name:        "findme",
		Weight:      10,
		IsAvailable: true,
	}
	err := repository.Insert(ctx, &productToInsert)
	if err != nil {
		t.Fatalf("Failed to insert product for FindByID test: %v", err)
	}

	foundEntity, err := repository.FindByID(ctx, entity.ID(productToInsert.ID))
	if err != nil {
		t.Errorf("FindByID failed: %v", err)
	}
	if foundEntity == nil {
		t.Fatal("FindByID returned nil entity")
	}

	foundProduct, ok := foundEntity.(Product)
	if !ok {
		t.Fatalf("FindByID returned type %T, expected Product", foundEntity)
	}

	if foundProduct.ID != productToInsert.ID {
		t.Errorf("FindByID: expected ID %d, got %d", productToInsert.ID, foundProduct.ID)
	}
	if foundProduct.Name != productToInsert.Name {
		t.Errorf("FindByID: expected Name %s, got %s", productToInsert.Name, foundProduct.Name)
	}
}

func TestGormRepository_Count(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)
	ctx := context.Background()

	// Clear table or use a transaction to ensure count is predictable
	db.Exec("DELETE FROM product_gorms")

	initialCount, err := repository.Count()
	if err != nil {
		t.Errorf("Count failed: %v", err)
	}
	if initialCount != 0 {
		t.Errorf("Expected initial count to be 0, got %d", initialCount)
	}

	p1 := Product{Name: "count_test_1"}
	repository.Insert(ctx, &p1)
	p2 := Product{Name: "count_test_2"}
	repository.Insert(ctx, &p2)

	nb, err := repository.Count()
	if err != nil {
		t.Errorf("Count failed: %v", err)
	}
	if nb != 2 {
		t.Errorf("Expected count to be 2, got %d", nb)
	}
}

func TestGormRepository_Update(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)
	ctx := context.Background()

	productToInsert := Product{
		ID:          200, // Distinct ID
		Name:        "updateme",
		Weight:      20,
		IsAvailable: false,
	}
	err := repository.Insert(ctx, &productToInsert)
	if err != nil {
		t.Fatalf("Failed to insert product for Update test: %v", err)
	}

	// productToUpdate := productToInsert // Create a copy
	productToUpdate := Product{
		ID:          productToInsert.ID,
		Name:        "updated_name",
		Weight:      productToInsert.Weight + 5,
		IsAvailable: !productToInsert.IsAvailable,
	}

	err = repository.Upsert(ctx, &productToUpdate) // Pass pointer
	if err != nil {
		t.Errorf("Upsert failed: %v", err)
	}

	foundAny, err := repository.FindByID(ctx, entity.ID(productToInsert.ID))
	if err != nil {
		t.Fatalf("FindByID after Upsert failed: %v", err)
	}
	foundProduct, ok := foundAny.(Product)
	if !ok {
		t.Fatalf("FindByID after Upsert returned wrong type: %T", foundAny)
	}

	if foundProduct.Name != "updated_name" {
		t.Errorf("Upsert: Name not updated. Expected 'updated_name', got '%s'", foundProduct.Name)
	}
	if foundProduct.Weight != productToInsert.Weight+5 {
		t.Errorf("Upsert: Weight not updated. Expected %d, got %d", productToInsert.Weight+5, foundProduct.Weight)
	}
	if foundProduct.IsAvailable == productToInsert.IsAvailable {
		t.Errorf("Upsert: IsAvailable not updated. Expected %t, got %t", !productToInsert.IsAvailable, foundProduct.IsAvailable)
	}
}

func TestGormRepository_DeleteByID(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)
	ctx := context.Background()

	productToInsert := Product{
		ID:   300, // Distinct ID
		Name: "deleteme",
	}
	err := repository.Insert(ctx, &productToInsert)
	if err != nil {
		t.Fatalf("Failed to insert product for DeleteByID test: %v", err)
	}

	err = repository.DeleteByID(ctx, entity.ID(productToInsert.ID))
	if err != nil {
		t.Errorf("DeleteByID failed: %v", err)
	}

	_, err = repository.FindByID(ctx, entity.ID(productToInsert.ID))
	if err == nil {
		t.Error("FindByID should have failed after DeleteByID, but it succeeded")
	}
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestGormRepository_Find(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)
	ctx := context.Background()
	db.Exec("DELETE FROM product_gorms")

	p1 := Product{ID: 1, Name: "find_p1", Weight: 100, IsAvailable: true}
	repository.Insert(ctx, &p1)
	p2 := Product{ID: 2, Name: "find_p2", Weight: 50, IsAvailable: true}
	repository.Insert(ctx, &p2)
	i := 22
	p3 := Product{ID: 3, Name: "find_p3", Weight: 250, IsAvailable: false, Qtt: &i}
	repository.Insert(ctx, &p3)
	p4 := Product{ID: 4, Name: "find_p4", Weight: 100, IsAvailable: true} // Another one with weight 100
	repository.Insert(ctx, &p4)

	resultsAny, err := repository.Find(ctx, gormrepository_reflection.GreaterOrEqual("weight", 100))
	if err != nil {
		t.Errorf("Find failed: %v", err)
	}
	if len(resultsAny) != 3 { // p1, p3, p4
		t.Errorf("Find with GreaterOrEqual('weight', 100): expected 3 results, got %d", len(resultsAny))
		for _, r := range resultsAny {
			t.Logf("%+v", r.(Product))
		}
	}

	resultsAny, err = repository.Find(ctx, gormrepository_reflection.And(
		gormrepository_reflection.GreaterOrEqual("weight", 90),
		gormrepository_reflection.Equal("is_available", true)),
	)
	if err != nil {
		t.Errorf("Find with And failed: %v", err)
	}
	if len(resultsAny) != 2 { // p1, p4
		t.Errorf("Find with And: expected 2 results, got %d", len(resultsAny))
		for _, r := range resultsAny {
			t.Logf("%+v", r.(Product))
		}
	}
}

func TestGormRepository_NewInstanceEntity(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)

	newEntityAny := repository.NewEntity()
	_, ok := newEntityAny.(*Product) // NewEntity returns a pointer to the zero value of the entity type
	if !ok {
		t.Errorf("NewEntity: expected *Product, got %T", newEntityAny)
	}
}

func TestGormRepository_Delete(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)
	ctx := context.Background()
	db.Exec("DELETE FROM product_gorms WHERE id = 8")

	product := Product{
		ID:          8, // Make sure this ID is unique for this test or clean up
		Name:        "product_to_delete_via_slice",
		Weight:      100,
		IsAvailable: true,
	}
	err := repository.Insert(ctx, &product)
	if err != nil {
		t.Fatalf("Insert failed before Delete: %v", err)
	}

	// Delete expects a slice of entity pointers
	err = repository.Delete([]any{&product})
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	_, err = repository.FindByID(ctx, entity.ID(product.ID))
	if err == nil {
		t.Errorf("Product with ID %d should have been deleted but was found", product.ID)
	}
}

func TestGormRepository_FindAll(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)
	ctx := context.Background()
	db.Exec("DELETE FROM product_gorms")

	p1 := Product{ID: 401, Name: "findall_1"}
	repository.Insert(ctx, &p1)
	p2 := Product{ID: 402, Name: "findall_2"}
	repository.Insert(ctx, &p2)

	resultsAny, err := repository.FindAll(ctx)
	if err != nil {
		t.Errorf("FindAll failed: %v", err)
	}
	if len(resultsAny) != 2 {
		t.Errorf("FindAll: expected 2 results, got %d", len(resultsAny))
	}
}

func TestGormRepository_GetDB(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)
	if repository.GetDB() != db {
		t.Error("GetDB should return the same DB instance")
	}
}

func TestGormRepository_Conditions(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)
	ctx := context.Background()
	db.Exec("DELETE FROM product_gorms")

	// Populate data
	p1 := Product{ID: 501, Name: "cond_p1", Weight: 150, IsAvailable: false} // Matched by GT 100, is_available false
	repository.Insert(ctx, &p1)
	p2 := Product{ID: 502, Name: "cond_p2", Weight: 50, IsAvailable: true} // Matched by LT 51, LTE 51, In(50,100)
	repository.Insert(ctx, &p2)
	p3 := Product{ID: 503, Name: "cond_p3", Weight: 100, IsAvailable: true} // Matched by In(50,100)
	repository.Insert(ctx, &p3)
	i := 30
	p4 := Product{ID: 504, Name: "cond_p4", Weight: 75, IsAvailable: false, Qtt: &i} // Matched by Not(is_available true), IsNotNull(qtt)
	repository.Insert(ctx, &p4)
	p5 := Product{ID: 505, Name: "cond_p5", Weight: 50, IsAvailable: false} // Matched by Not(is_available true), LT 51, LTE 51, In(50,100)
	repository.Insert(ctx, &p5)

	// Test cases
	testCases := []struct {
		name          string
		spec          gormrepository_reflection.Specification
		expectedCount int
	}{
		{
			"GreaterThan And Equal",
			gormrepository_reflection.And(
				gormrepository_reflection.GreaterThan("weight", 100),
				gormrepository_reflection.Equal("is_available", false),
			),
			1, // p1
		},
		{
			"Not Equal",
			gormrepository_reflection.Not(gormrepository_reflection.Equal("is_available", true)),
			3, // p1, p4, p5
		},
		{
			"LessThan",
			gormrepository_reflection.LessThan("weight", 51), // p2, p5
			2,
		},
		{
			"LessOrEqual",
			gormrepository_reflection.LessOrEqual("weight", 50), // p2, p5
			2,
		},
		{
			"In",
			gormrepository_reflection.In("weight", []any{uint(50), uint(100)}), // p2, p3, p5 (uint for GORM)
			3,
		},
		{
			"Or",
			gormrepository_reflection.Or(
				gormrepository_reflection.Equal("weight", 50),
				gormrepository_reflection.Equal("weight", 100),
			),
			3, // p2, p3, p5
		},
		{
			"IsNull",
			gormrepository_reflection.IsNull("qtt"), // p1, p2, p3, p5
			4,
		},
		{
			"IsNotNull",
			gormrepository_reflection.IsNotNull("qtt"), // p4
			1,
		},
		{
			"Not IsNull",
			gormrepository_reflection.Not(gormrepository_reflection.IsNull("qtt")), // p4
			1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultsAny, err := repository.Find(ctx, tc.spec)
			if err != nil {
				t.Errorf("Find with spec '%s' failed: %v", tc.name, err)
				return
			}
			if len(resultsAny) != tc.expectedCount {
				t.Errorf("Find with spec '%s': expected %d results, got %d", tc.name, tc.expectedCount, len(resultsAny))
				for idx, r := range resultsAny {
					t.Logf("  Result %d: %+v", idx, r.(Product))
				}
			}
		})
	}
}

func TestGormRepository_Order(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)
	ctx := context.Background()
	db.Exec("DELETE FROM product_gorms")

	p1 := Product{ID: 601, Name: "order_p1", Weight: 250}
	repository.Insert(ctx, &p1)
	p2 := Product{ID: 602, Name: "order_p2", Weight: 50}
	repository.Insert(ctx, &p2)
	p3 := Product{ID: 603, Name: "order_p3", Weight: 100}
	repository.Insert(ctx, &p3)

	resultsAny, err := repository.FindAll(ctx, gormrepository_reflection.OrderBy("weight", "DESC"))
	if err != nil {
		t.Errorf("FindAll with OrderBy DESC failed: %v", err)
	}
	if len(resultsAny) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(resultsAny))
	}

	firstProduct, ok := resultsAny[0].(Product)
	if !ok {
		t.Fatalf("Result is not Product: %T", resultsAny[0])
	}
	if firstProduct.Weight != 250 {
		t.Errorf("Expected first product weight to be 250 (DESC order), got %d", firstProduct.Weight)
	}

	resultsAnyAsc, err := repository.FindAll(ctx, gormrepository_reflection.OrderBy("weight", "ASC"))
	if err != nil {
		t.Errorf("FindAll with OrderBy ASC failed: %v", err)
	}
	if len(resultsAnyAsc) != 3 {
		t.Fatalf("Expected 3 results for ASC, got %d", len(resultsAnyAsc))
	}
	firstProductAsc, ok := resultsAnyAsc[0].(Product)
	if !ok {
		t.Fatalf("Result (ASC) is not Product: %T", resultsAnyAsc[0])
	}
	if firstProductAsc.Weight != 50 {
		t.Errorf("Expected first product weight to be 50 (ASC order), got %d", firstProductAsc.Weight)
	}
}

func TestGormRepository_PreloadAssociation(t *testing.T) {
	db, _ := getDB()
	// This test mainly checks if the methods can be called without panic.
	// Actual preload testing would require models with defined associations.
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)
	repository.EnablePreloadAssociations()
	repository.DisablePreloadAssociations()
	repository.SetPreloadAssociations(true)
	repository.SetPreloadAssociations(false)
	// No assertions needed, just checking for panics.
}

func TestGormRepository_BatchInsert(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)
	ctx := context.Background()
	db.Exec("DELETE FROM product_gorms WHERE id >= 700 AND id < 800")

	products := []any{ // Slice of any, elements are *Product
		&Product{
			ID:          700,
			Name:        "batch_insert_1",
			Weight:      100,
			IsAvailable: true,
		},
		&Product{
			ID:          701,
			Name:        "batch_insert_2",
			Weight:      50,
			IsAvailable: true,
		},
	}
	err := repository.BatchInsert(ctx, products)
	if err != nil {
		t.Errorf("BatchInsert failed: %v", err)
	}
	/*
	   	TODO

	   the count(condition) does not exists but is prettu good

	   	//count, err := repository.Count(gormrepository_reflection.In("id", []any{uint(700), uint(701)}))
	   	if err != nil {
	   		t.Errorf("Count after BatchInsert failed: %v", err)
	   	}
	   	if count != 2 {
	   		t.Errorf("Expected count of 2 after BatchInsert, got %d", count)
	   	}
	*/
}

func TestGormRepository_BatchUpdate(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)
	ctx := context.Background()
	db.Exec("DELETE FROM product_gorms WHERE id >= 800 AND id < 900")

	p1 := Product{ID: 800, Name: "batch_update_orig_1", Weight: 10}
	p2 := Product{ID: 801, Name: "batch_update_orig_2", Weight: 20}
	repository.BatchInsert(ctx, []any{&p1, &p2})

	productsToUpdate := []any{
		&Product{
			ID:          800, // Existing ID
			Name:        "batch_update_MOD_1",
			Weight:      15,
			IsAvailable: true,
		},
		&Product{
			ID:          801, // Existing ID
			Name:        "batch_update_MOD_2",
			Weight:      25,
			IsAvailable: false,
		},
		// &Product{ // Test with a new product - should also work due to Upsert nature of Save
		// 	ID: 802,
		// 	Name: "batch_update_NEW_3",
		// 	Weight: 35,
		// },
	}
	err := repository.BatchUpdate(ctx, productsToUpdate)
	if err != nil {
		t.Errorf("BatchUpdate failed: %v", err)
	}

	updatedP1Any, _ := repository.FindByID(ctx, entity.ID(800))
	updatedP1 := updatedP1Any.(Product)
	if updatedP1.Name != "batch_update_MOD_1" || updatedP1.Weight != 15 {
		t.Errorf("Product 800 not updated correctly: %+v", updatedP1)
	}

	updatedP2Any, _ := repository.FindByID(ctx, entity.ID(801))
	updatedP2 := updatedP2Any.(Product)
	if updatedP2.Name != "batch_update_MOD_2" || updatedP2.Weight != 25 {
		t.Errorf("Product 801 not updated correctly: %+v", updatedP2)
	}

	// newP3Any, errP3 := repository.FindByID(ctx, entity.ID(802))
	// if errP3 != nil {
	// 	t.Errorf("Product 802 (new) should have been created by BatchUpdate(Save): %v", errP3)
	// } else if newP3Any.(Product).Name != "batch_update_NEW_3" {
	// 	t.Errorf("Product 802 (new) has wrong name: %+v", newP3Any.(Product))
	// }
}

func TestGormRepository_BatchDelete(t *testing.T) {
	db, _ := getDB()
	repository := gormrepository_reflection.NewRepository(db, productModelType, productEntityType)
	ctx := context.Background()
	db.Exec("DELETE FROM product_gorms WHERE id >= 900 AND id < 1000")

	p1 := Product{ID: 900, Name: "batch_delete_1"}
	p2 := Product{ID: 901, Name: "batch_delete_2"}
	p3 := Product{ID: 902, Name: "batch_delete_3_keeps"} // This one is not in the delete list
	repository.BatchInsert(ctx, []any{&p1, &p2, &p3})

	productsToDelete := []any{
		&Product{ID: 900}, // Only ID is needed for deletion by GORM
		&Product{ID: 901},
	}
	err := repository.BatchDelete(ctx, productsToDelete)
	if err != nil {
		t.Errorf("BatchDelete failed: %v", err)
	}

	_, errP1 := repository.FindByID(ctx, entity.ID(900))
	if errP1 == nil {
		t.Error("Product 900 should have been deleted by BatchDelete")
	}
	_, errP2 := repository.FindByID(ctx, entity.ID(901))
	if errP2 == nil {
		t.Error("Product 901 should have been deleted by BatchDelete")
	}
	foundP3, errP3 := repository.FindByID(ctx, entity.ID(902))
	if errP3 != nil {
		t.Errorf("Product 902 should NOT have been deleted: %v", errP3)
	}
	if foundP3.(Product).ID != 902 {
		t.Error("Product 902 data is incorrect after batch delete.")
	}
}

// Ensure ProductGorm implements entity.DatabaseModel[entity.Entity]
// This is a compile-time check.
var _ entity.DatabaseModel[entity.Entity] = (*ProductGorm)(nil)

// Ensure Product implements entity.Entity
var _ entity.Entity = Product{} // Product is a struct, methods have value receivers for GetId, ToModel
// SetId is not fully implemented and panics, which is fine for this check if not called directly.

// ToModel method for Product to implement entity.Entity if it requires converting to a DatabaseModel
func (p Product) ToModel() entity.DatabaseModel[entity.Entity] {
	pg := &ProductGorm{}
	pg.FromEntity(p) // pg is a pointer, FromEntity has a pointer receiver and now returns any, but we use it for its side effect here.
	return pg        // pg (*ProductGorm) implements entity.DatabaseModel[entity.Entity]
}
