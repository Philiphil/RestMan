package mongorepository

import (
	"context"

	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
)

// Create implements RestRepository.Create by inserting entities.
func (r *MongoRepository[M, E]) Create(ctx context.Context, entities []*E) error {
	return r.BatchInsert(ctx, entities)
}

// Update implements RestRepository.Update by updating entities.
func (r *MongoRepository[M, E]) Update(ctx context.Context, entities []*E) error {
	return r.BatchUpdate(ctx, entities)
}

// Read implements RestRepository.Read by finding entities by IDs.
func (r *MongoRepository[M, E]) Read(ctx context.Context, ids []entity.ID) ([]*E, error) {
	return r.FindByIDs(ctx, ids)
}

// Delete implements RestRepository.Delete by deleting entities.
func (r *MongoRepository[M, E]) Delete(ctx context.Context, entities []*E) error {
	return r.BatchDelete(ctx, entities)
}

// List implements RestRepository.List by finding entities with pagination and sorting.
func (r *MongoRepository[M, E]) List(ctx context.Context, limit int, offset int, order map[string]string) ([]E, error) {
	orderSpecifications := make([]Specification, 0, len(order))
	for k, v := range order {
		orderSpecifications = append(orderSpecifications, OrderBy(k, v))
	}
	return r.FindWithLimit(ctx, limit, offset, orderSpecifications...)
}

// New implements RestRepository.New by creating a new entity instance.
func (r *MongoRepository[M, E]) New() E {
	return r.NewEntity()
}

// Count implements RestRepository.Count by returning the total number of documents.
func (r *MongoRepository[M, E]) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, map[string]interface{}{})
}

// LoadAssociations implements RestRepository.LoadAssociations as a no-op, MongoDB has no associations to preload.
func (r *MongoRepository[M, E]) LoadAssociations(associations []string) (orm.RestRepository[entity.DatabaseModel[E], E], error) {
	return r, nil
}
