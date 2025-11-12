package gormrepository

import (
	"context"

	"github.com/philiphil/restman/orm/entity"
)

// Create implements RestRepository.Create by inserting entities.
func (r *GormRepository[M, E]) Create(ctx context.Context, entities []*E) error {
	return r.BatchInsert(ctx, entities)
}

// Update implements RestRepository.Update by updating entities.
func (r *GormRepository[M, E]) Update(ctx context.Context, entities []*E) error {
	return r.BatchUpdate(ctx, entities)
}

// Read implements RestRepository.Read by finding entities by IDs.
func (r *GormRepository[M, E]) Read(ctx context.Context, ids []entity.ID) ([]*E, error) {
	return r.FindByIDs(ctx, ids)
}

// Delete implements RestRepository.Delete by deleting entities.
func (r *GormRepository[M, E]) Delete(ctx context.Context, entities []*E) error {
	return r.BatchDelete(ctx, entities)
}

// List implements RestRepository.List by finding entities with pagination and sorting.
func (r *GormRepository[M, E]) List(ctx context.Context, limit int, offset int, order map[string]string) ([]E, error) {
	orderSpecification := make([]Specification, 0, len(order))
	for k, v := range order {
		orderSpecification = append(orderSpecification, OrderBy(k, v))
	}
	return r.FindWithLimit(ctx, limit, offset, orderSpecification...)
}

// New implements RestRepository.New by creating a new entity instance.
func (r *GormRepository[M, E]) New() E {
	return r.NewEntity()
}

// Count implements RestRepository.Count by returning the total number of entities.
func (r *GormRepository[M, E]) Count(ctx context.Context) (i int64, err error) {
	model := new(M)
	err = r.getPreWarmDbForSelect(ctx).Model(model).Count(&i).Error
	return
}
