package orm

import (
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/orm/entity"
)

// ORM is the struct used by RestMan to interact with the repository
// The repository is an interface impleemting the CRUD operations
// You can implement your own repository or use the default one
type ORM[T entity.Entity] struct {
	Repo RestRepository[entity.DatabaseModel[T], T]
}

func NewORM[T entity.Entity](repo RestRepository[entity.DatabaseModel[T], T]) *ORM[T] {
	return &ORM[T]{
		Repo: repo,
	}
}

func (r *ORM[T]) GetAll(sort map[string]string) ([]T, error) {
	return r.Repo.List(-1, -1, sort)
}

func (r *ORM[T]) GetByID(id any) (*T, error) {
	elem, err := r.Repo.Read([]entity.ID{entity.CastId(id)})

	if err != nil {
		return nil, err
	}
	if len(elem) == 0 {
		return nil, errors.ItemNotFound
	}
	return elem[0], nil
}

func (r *ORM[T]) GetPaginatedList(itemPerPage int, page int, sort map[string]string) ([]T, error) {
	return r.Repo.List(itemPerPage, itemPerPage*page, sort)
}

func (r *ORM[T]) Count() (int64, error) {
	return r.Repo.Count()
}

func (r *ORM[T]) Create(item ...*T) error {
	return r.Repo.Create(item)
}

func (r *ORM[T]) Update(item ...*T) error {
	return r.Repo.Update(item)
}

func (r *ORM[T]) Delete(item ...*T) error {
	return r.Repo.Delete(item)
}

func (r *ORM[T]) FindByIDs(ids []entity.ID) ([]*T, error) {
	list, err := r.Repo.Read(ids)
	if err != nil {
		return nil, err
	}
	if len(list) != len(ids) {
		return nil, errors.NotAllItemFound
	}
	return list, nil
}

func (r *ORM[T]) NewEntity() T {
	return r.Repo.New()
}
