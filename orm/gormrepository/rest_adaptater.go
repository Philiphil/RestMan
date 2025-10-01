package gormrepository

import (
	"context"

	"github.com/philiphil/restman/orm/entity"
)

func (r *GormRepository[M, E]) Create(entities []*E) error {
	return r.BatchInsert(context.Background(), entities)
}

func (r *GormRepository[M, E]) Update(entities []*E) error {
	return r.BatchUpdate(context.Background(), entities)
}

func (r *GormRepository[M, E]) Read(ids []entity.ID) ([]*E, error) {
	return r.FindByIDs(context.Background(), ids)
}

func (r *GormRepository[M, E]) Delete(entities []*E) error {
	return r.BatchDelete(context.Background(), entities)
}

func (r *GormRepository[M, E]) List(limit int, offset int, order map[string]string) ([]E, error) {
	orderSpecification := make([]Specification, 0, len(order))
	for k, v := range order {
		orderSpecification = append(orderSpecification, OrderBy(k, v))
	}
	return r.FindWithLimit(context.Background(), limit, offset, orderSpecification...)
}

func (r *GormRepository[M, E]) New() E {
	return r.NewEntity()
}

func (r *GormRepository[M, E]) Count() (i int64, err error) {
	model := new(M)
	err = r.getPreWarmDbForSelect(context.TODO()).Model(model).Count(&i).Error
	return
}
