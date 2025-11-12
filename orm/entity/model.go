package entity

// This allows mapping between database models and entities
// this is for DDD style separation between database models and domain entities.
// this can be ignored for simple use cases where the database model is the entity.
type DatabaseModel[E Entity] interface {
	ToEntity() E
	FromEntity(entity E) any
}
