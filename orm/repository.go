// The package is the layer that interacts with the database.
package orm

import (
	"context"

	"github.com/philiphil/restman/orm/entity"
)

// RestRepository defines the interface for database operations on entities.
type RestRepository[M entity.DatabaseModel[E], E entity.Entity] interface {
	Create(ctx context.Context, entities []*E) error
	Read(ctx context.Context, ids []entity.ID) ([]*E, error)
	Update(ctx context.Context, entities []*E) error
	Delete(ctx context.Context, entities []*E) error
	List(ctx context.Context, limit int, offset int, order map[string]string) ([]E, error)
	Count(ctx context.Context) (int64, error)

	New() E
}
