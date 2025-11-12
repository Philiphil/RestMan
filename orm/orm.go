package orm

import (
	"context"

	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/orm/entity"
)

// ORM is the struct used by RestMan to interact with the repository
// The repository is an interface impleemting the CRUD operations
// You can implement your own repository or use the default one
type ORM[T entity.Entity] struct {
	Repo RestRepository[entity.DatabaseModel[T], T]
}

// NewORM creates a new ORM instance with the provided repository.
func NewORM[T entity.Entity](repo RestRepository[entity.DatabaseModel[T], T]) *ORM[T] {
	return &ORM[T]{
		Repo: repo,
	}
}

// GetAll retrieves all entities from the repository with optional sorting.
func (r *ORM[T]) GetAll(ctx context.Context, sort map[string]string) ([]T, error) {
	return r.Repo.List(ctx, -1, -1, sort)
}

// GetByID retrieves a single entity by its ID.
func (r *ORM[T]) GetByID(ctx context.Context, id any) (*T, error) {
	elem, err := r.Repo.Read(ctx, []entity.ID{entity.CastId(id)})

	if err != nil {
		return nil, err
	}
	if len(elem) == 0 {
		return nil, errors.ItemNotFound
	}
	return elem[0], nil
}

// GetPaginatedList retrieves a paginated list of entities with optional sorting.
func (r *ORM[T]) GetPaginatedList(ctx context.Context, itemPerPage int, page int, sort map[string]string) ([]T, error) {
	return r.Repo.List(ctx, itemPerPage, itemPerPage*page, sort)
}

// Count returns the total number of entities in the repository.
func (r *ORM[T]) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx)
}

// Create persists one or more new entities to the repository.
func (r *ORM[T]) Create(ctx context.Context, item ...*T) error {
	return r.Repo.Create(ctx, item)
}

// Update modifies one or more existing entities in the repository.
func (r *ORM[T]) Update(ctx context.Context, item ...*T) error {
	return r.Repo.Update(ctx, item)
}

// Delete removes one or more entities from the repository.
func (r *ORM[T]) Delete(ctx context.Context, item ...*T) error {
	return r.Repo.Delete(ctx, item)
}

// FindByIDs retrieves multiple entities by their IDs, returning an error if not all are found.
func (r *ORM[T]) FindByIDs(ctx context.Context, ids []entity.ID) ([]*T, error) {
	list, err := r.Repo.Read(ctx, ids)
	if err != nil {
		return nil, err
	}
	if len(list) != len(ids) {
		return nil, errors.NotAllItemFound
	}
	return list, nil
}

// NewEntity creates a new empty entity instance.
func (r *ORM[T]) NewEntity() T {
	return r.Repo.New()
}
