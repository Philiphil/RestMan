package apiman

import (
	"github.com/philiphil/apiman/orm/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type Test1 struct {
	entity.Entity
}

func (e Test1) GetId() entity.ID {
	return e.Id
}

func (e Test1) SetId(id any) entity.IEntity {
	e.Id = entity.CastId(id)
	return e
}
func (t Test1) ToEntity() Test1 {
	return t
}

func (t Test1) FromEntity(entity Test1) any {
	return entity
}

func getDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(os.Getenv("DSN")), &gorm.Config{
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

/*
func TestMain(m *testing.M) {
	test_ := NewApiRouter[Test1](
		*orm.NewORM[Test1](repository.NewRepository[Test1, Test1](getDB())),
		method.DefaultApiMethods(),
	)
	_ = test_
}
*/
