package orm

import (
	"context"
)

type IRepository[M GormModel[E], E any] interface {
	Insert(ctx context.Context, entity *E) error
	Delete(ctx context.Context, entity *E) error
	DeleteByID(ctx context.Context, id any) error
	Update(ctx context.Context, entity *E) error
	FindByID(ctx context.Context, id any) (E, error)
	Find(ctx context.Context, specifications ...Specification) ([]E, error)
	Count(ctx context.Context, specifications ...Specification) (int64, error)
	FindWithLimit(ctx context.Context, limit int, offset int, specifications ...Specification) ([]E, error)
	FindAll(ctx context.Context) ([]E, error)
	NewEntity() E
}
