package orm

import (
	"context"
	"github.com/philiphil/apiman/orm/entity"
)

type ORM[T entity.IEntity] struct {
	repo IRepository[entity.Model[T], T]
}

func NewORM[T entity.IEntity](repo IRepository[entity.Model[T], T]) *ORM[T] {
	return &ORM[T]{
		repo: repo,
	}
}

func (r *ORM[T]) GetAll() ([]T, error) {
	ctx := context.Background()
	items, err := r.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ORM[T]) GetByID(id any) (*T, error) {
	ctx := context.Background()
	item, err := r.repo.FindByID(ctx, entity.CastId(id))
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ORM[T]) GetPaginatedList(itemPerPage int, page int) ([]T, error) {
	ctx := context.Background()
	many, err := r.repo.FindWithLimit(ctx, itemPerPage, itemPerPage*page)
	return many, err
}

func (r *ORM[T]) Count() (int64, error) {
	ctx := context.Background()
	i, err := r.repo.Count(ctx)
	return i, err
}

func (r *ORM[T]) Create(item *T) error {
	ctx := context.Background()
	err := r.repo.Insert(ctx, item)
	if err != nil {
		return err
	}

	return nil
}

func (r *ORM[T]) Update(item *T) error {
	ctx := context.Background()
	err := r.repo.Update(ctx, item)
	if err != nil {
		return err
	}
	return nil
}

func (r *ORM[T]) Delete(id any) error {
	ctx := context.Background()
	err := r.repo.DeleteByID(ctx, entity.CastId(id))
	if err != nil {
		return err
	}
	return nil
}

func (r *ORM[T]) NewEntity() T {
	return r.repo.NewEntity()
}
