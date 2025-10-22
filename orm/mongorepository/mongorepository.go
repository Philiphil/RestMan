package mongorepository

import (
	"context"

	"github.com/philiphil/restman/orm/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewRepository creates a new MongoRepository instance with the provided MongoDB collection.
func NewRepository[M entity.DatabaseModel[E], E entity.Entity](collection *mongo.Collection) *MongoRepository[M, E] {
	return &MongoRepository[M, E]{
		collection: collection,
	}
}

// MongoRepository is a MongoDB-based implementation of the RestRepository interface.
type MongoRepository[M entity.DatabaseModel[E], E entity.Entity] struct {
	collection *mongo.Collection
}

// Insert creates a new entity in the MongoDB collection.
func (r *MongoRepository[M, E]) Insert(ctx context.Context, entity *E) error {
	var start M
	model := start.FromEntity(*entity).(M)

	result, err := r.collection.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	(*entity).SetId(result.InsertedID)

	return nil
}

// DeleteByID removes an entity from the MongoDB collection by its ID.
func (r *MongoRepository[M, E]) DeleteByID(ctx context.Context, id entity.ID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// Upsert creates or updates an entity in the MongoDB collection.
func (r *MongoRepository[M, E]) Upsert(ctx context.Context, entity *E) error {
	var start M
	model := start.FromEntity(*entity).(M)

	id := (*entity).GetId()
	filter := bson.M{"_id": id}
	opts := options.Replace().SetUpsert(true)

	_, err := r.collection.ReplaceOne(ctx, filter, model, opts)
	if err != nil {
		return err
	}

	return nil
}

// FindByID retrieves an entity by its ID from the MongoDB collection.
func (r *MongoRepository[M, E]) FindByID(ctx context.Context, id entity.ID) (E, error) {
	var model M
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		return *new(E), err
	}

	return model.ToEntity(), nil
}

// Find retrieves entities matching the provided specifications.
func (r *MongoRepository[M, E]) Find(ctx context.Context, specifications ...Specification) ([]E, error) {
	return r.FindWithLimit(ctx, -1, -1, specifications...)
}

func (r *MongoRepository[M, E]) buildFilter(specifications []Specification) bson.M {
	filters := make([]bson.M, 0)

	for _, spec := range specifications {
		if filter := spec.GetFilter(); filter != nil {
			filters = append(filters, filter)
		}
	}

	if len(filters) == 0 {
		return bson.M{}
	}

	if len(filters) == 1 {
		return filters[0]
	}

	return bson.M{"$and": filters}
}

func (r *MongoRepository[M, E]) buildOptions(limit int, offset int, specifications []Specification) *options.FindOptions {
	opts := options.Find()

	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	if offset > 0 {
		opts.SetSkip(int64(offset))
	}

	for _, spec := range specifications {
		if sort := spec.GetSort(); sort != nil {
			opts.SetSort(sort)
		}
	}

	return opts
}

// FindWithLimit retrieves entities matching the provided specifications with pagination.
func (r *MongoRepository[M, E]) FindWithLimit(ctx context.Context, limit int, offset int, specifications ...Specification) ([]E, error) {
	filter := r.buildFilter(specifications)
	opts := r.buildOptions(limit, offset, specifications)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var models []M
	if err := cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	result := make([]E, 0, len(models))
	for _, model := range models {
		result = append(result, model.ToEntity())
	}

	return result, nil
}

// FindAll retrieves all entities matching the provided specifications.
func (r *MongoRepository[M, E]) FindAll(ctx context.Context, specification ...Specification) ([]E, error) {
	return r.FindWithLimit(ctx, -1, -1, specification...)
}

// FindByIDs retrieves multiple entities by their IDs.
func (r *MongoRepository[M, E]) FindByIDs(ctx context.Context, ids []entity.ID) ([]*E, error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var models []M
	if err := cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	result := make([]*E, 0, len(models))
	for _, model := range models {
		entity := model.ToEntity()
		result = append(result, &entity)
	}

	return result, nil
}

// DeleteByIDs removes multiple entities by their IDs.
func (r *MongoRepository[M, E]) DeleteByIDs(ctx context.Context, ids []entity.ID) error {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

// BatchDelete removes multiple entities in a single operation.
func (r *MongoRepository[M, E]) BatchDelete(ctx context.Context, entities []*E) error {
	ids := make([]entity.ID, 0, len(entities))
	for _, entity := range entities {
		ids = append(ids, (*entity).GetId())
	}
	return r.DeleteByIDs(ctx, ids)
}

// BatchUpdate updates multiple entities.
func (r *MongoRepository[M, E]) BatchUpdate(ctx context.Context, entities []*E) error {
	for _, entity := range entities {
		if err := r.Upsert(ctx, entity); err != nil {
			return err
		}
	}
	return nil
}

// BatchInsert creates multiple entities in a single operation.
func (r *MongoRepository[M, E]) BatchInsert(ctx context.Context, entities []*E) error {
	if len(entities) == 0 {
		return nil
	}

	docs := make([]interface{}, 0, len(entities))
	for _, entity := range entities {
		var start M
		model := start.FromEntity(*entity).(M)
		docs = append(docs, model)
	}

	results, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	for i, insertedID := range results.InsertedIDs {
		(*entities[i]).SetId(insertedID)
	}

	return nil
}

// GetCollection returns the underlying MongoDB collection.
func (r *MongoRepository[M, E]) GetCollection() *mongo.Collection {
	return r.collection
}

// NewEntity creates a new empty entity instance.
func (r *MongoRepository[M, E]) NewEntity() E {
	var entity E
	return entity
}

// CountWithSpecifications returns the number of entities matching the provided specifications.
func (r *MongoRepository[M, E]) CountWithSpecifications(ctx context.Context, specifications ...Specification) (int64, error) {
	filter := r.buildFilter(specifications)
	return r.collection.CountDocuments(ctx, filter)
}
