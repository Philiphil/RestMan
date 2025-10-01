package gormrepository_reflection

import (
	"context"
	"reflect"

	"github.com/philiphil/restman/orm/entity"
)

func (r *GormRepository) Create(entities []any) error {
	return r.BatchInsert(context.Background(), entities)
}

func (r *GormRepository) Update(entities []any) error {
	return r.BatchUpdate(context.Background(), entities)
}

func (r *GormRepository) Read(ids []entity.ID) ([]any, error) {
	return r.FindByIDs(context.Background(), ids)
}

func (r *GormRepository) Delete(entities []any) error {
	return r.BatchDelete(context.Background(), entities)
}

func (r *GormRepository) List(limit int, offset int, order map[string]string) ([]any, error) {
	orderSpecifications := make([]Specification, 0, len(order))
	for k, v := range order {
		orderSpecifications = append(orderSpecifications, OrderBy(k, v))
	}
	return r.FindWithLimit(context.Background(), limit, offset, orderSpecifications...)
}

func (r *GormRepository) New() any {
	return r.NewEntity()
}

func (r *GormRepository) Count() (i int64, err error) {
	// Create a new instance of the model type for counting
	modelInstance := reflect.New(r.modelType).Interface()
	err = r.getPreWarmDbForSelect(context.TODO()).Model(modelInstance).Count(&i).Error
	return
}
