package orm

import (
	"context"
	"github.com/philiphil/apiman/orm/entity"
	"github.com/philiphil/apiman/orm/repository"
)

type IRepository[M entity.Model[E], E entity.IEntity] interface {
	Insert(ctx context.Context, entity *E) error
	Delete(ctx context.Context, entity *E) error
	DeleteByID(ctx context.Context, id entity.ID) error
	Update(ctx context.Context, entity *E) error
	FindByID(ctx context.Context, id entity.ID) (E, error)
	Find(ctx context.Context, specifications ...repository.Specification) ([]E, error)
	Count(ctx context.Context, specifications ...repository.Specification) (int64, error)
	FindWithLimit(ctx context.Context, limit int, offset int, specifications ...repository.Specification) ([]E, error)
	FindAll(ctx context.Context) ([]E, error)
	NewEntity() E
}
