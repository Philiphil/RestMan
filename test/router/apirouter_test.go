package router_test

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository"
	"github.com/philiphil/restman/route"
	. "github.com/philiphil/restman/router"
	"github.com/philiphil/restman/security"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Test struct {
	entity.BaseEntity
}

func (e Test) GetId() entity.ID {
	return e.Id
}

func (e Test) SetId(id any) entity.Entity {
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
		*orm.NewORM[Test](gormrepository.NewRepository[Test, Test](getDB())),
		route.DefaultApiRoutes(),
	)
	if test_ == nil {
		t.Error("Expected not nil")
	}
}

// Test the route generation of the router
// it should start with a / and end without a /
func TestNewApiRouter2(t *testing.T) {
	test_ := NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Test](getDB())),
		route.DefaultApiRoutes(),
		configuration.RoutePrefix("v2"),
		configuration.RouteName("test"),
	)

	if test_.Route() != "/v2/test" {
		t.Error("Expected v2/test")
	}

	test_ = NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Test](getDB())),
		route.DefaultApiRoutes(),
		configuration.RoutePrefix("/v2/"),
		configuration.RouteName("/test/"),
	)
	if test_.Route() != "/v2/test" {
		fmt.Println(test_.Route())
		t.Error("Expected v2/test")
	}
}

type AnonymousUser struct {
}

func (u AnonymousUser) GetId() entity.ID {
	return 0
}
func (u AnonymousUser) SetId(id any) entity.Entity {
	return u
}

type AnonymousFirewall struct {
}

func (f AnonymousFirewall) GetUser(*gin.Context) (security.User, error) {
	return nil, nil
}

func TestApiRouter_AddFirewall(t *testing.T) {
	test_ := NewApiRouter[Test](
		*orm.NewORM[Test](gormrepository.NewRepository[Test, Test](getDB())),
		route.DefaultApiRoutes(),
	)
	test_.AddFirewall(AnonymousFirewall{})
	if len(test_.Firewalls) != 1 {
		t.Error("Expected 1")
	}
	test_.AddFirewall(AnonymousFirewall{})
	if len(test_.Firewalls) != 2 {
		t.Error("Expected 2")
	}
}

func TestTrimSlash(t *testing.T) {
	if TrimSlash("/test/") != "test" {
		t.Error("Expected test")
	}
	if TrimSlash("test/") != "test" {
		t.Error("Expected test")
	}
	if TrimSlash("/test") != "test" {
		t.Error("Expected test")
	}
	if TrimSlash("test") != "test" {
		t.Error("Expected test")
	}
}
