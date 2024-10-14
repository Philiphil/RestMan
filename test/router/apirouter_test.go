package router_test

import (
	"fmt"
	"testing"

	"github.com/philiphil/restman/method"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/repository"
	. "github.com/philiphil/restman/router"
	"github.com/philiphil/restman/security"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Test struct {
	entity.Entity
}

func (e Test) GetId() entity.ID {
	return e.Id
}

func (e Test) SetId(id any) entity.IEntity {
	e.Id = entity.CastId(id)
	return e
}
func (t Test) ToEntity() Test {
	return t
}

func (t Test) FromEntity(entity Test) any {
	return entity
}

func getDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:test?mode=memory&cache=shared&_fk=1"), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Silent),
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic(err)
	}
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	})

	db.Session(&gorm.Session{FullSaveAssociations: true})
	return db
}

func TestNewApiRouter(t *testing.T) {
	test_ := NewApiRouter[Test](
		*orm.NewORM[Test](repository.NewRepository[Test, Test](getDB())),
		method.DefaultApiMethods(),
	)
	if test_ == nil {
		t.Error("Expected not nil")
	}
}

// Test the route generation of the router
// it should start with a / and end without a /
func TestNewApiRouter2(t *testing.T) {
	test_ := NewApiRouter[Test](
		*orm.NewORM[Test](repository.NewRepository[Test, Test](getDB())),
		method.DefaultApiMethods(),
		"v2/test",
	)

	if test_.Route != "/v2/test" {
		t.Error("Expected v2/test")
	}

	test_ = NewApiRouter[Test](
		*orm.NewORM[Test](repository.NewRepository[Test, Test](getDB())),
		method.DefaultApiMethods(),
		"v2/test/",
	)
	if test_.Route != "/v2/test" {
		fmt.Println(test_.Route)
		t.Error("Expected v2/test")
	}
	test_ = NewApiRouter[Test](
		*orm.NewORM[Test](repository.NewRepository[Test, Test](getDB())),
		method.DefaultApiMethods(),
		"/v2/test",
	)

	if test_.Route != "/v2/test" {
		t.Error("Expected v2/test")
	}
}

func TestApiRouter_AddFirewall(t *testing.T) {
	test_ := NewApiRouter[Test](
		*orm.NewORM[Test](repository.NewRepository[Test, Test](getDB())),
		method.DefaultApiMethods(),
	)
	test_.AddFirewall(security.AnonymousFirewall{})
	if len(test_.Firewalls) != 1 {
		t.Error("Expected 1")
	}
	test_.AddFirewall(security.AnonymousFirewall{})
	if len(test_.Firewalls) != 2 {
		t.Error("Expected 2")
	}
}
