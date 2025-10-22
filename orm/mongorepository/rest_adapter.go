package mongorepository

import (
	"context"

	"github.com/philiphil/restman/orm/entity"
)

func (r *MongoRepository[M, E]) Create(entities []*E) error {
	return r.BatchInsert(context.Background(), entities)
}

func (r *MongoRepository[M, E]) Update(entities []*E) error {
	return r.BatchUpdate(context.Background(), entities)
}

func (r *MongoRepository[M, E]) Read(ids []entity.ID) ([]*E, error) {
	return r.FindByIDs(context.Background(), ids)
}

func (r *MongoRepository[M, E]) Delete(entities []*E) error {
	return r.BatchDelete(context.Background(), entities)
}

func (r *MongoRepository[M, E]) List(limit int, offset int, order map[string]string) ([]E, error) {
	orderSpecifications := make([]Specification, 0, len(order))
	for k, v := range order {
		orderSpecifications = append(orderSpecifications, OrderBy(k, v))
	}
	return r.FindWithLimit(context.Background(), limit, offset, orderSpecifications...)
}

func (r *MongoRepository[M, E]) New() E {
	return r.NewEntity()
}

func (r *MongoRepository[M, E]) Count() (int64, error) {
	return r.collection.CountDocuments(context.Background(), map[string]interface{}{})
}
