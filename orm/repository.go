// The package is the layer that interacts with the database.
package orm

import (
	"github.com/philiphil/restman/orm/entity"
)

// RestRepository defines the interface for database operations on entities.
type RestRepository[M entity.DatabaseModel[E], E entity.Entity] interface {
	Create(entities []*E) error
	Read(ids []entity.ID) ([]*E, error)
	Update(entities []*E) error
	Delete(entities []*E) error
	List(limit int, offset int, order map[string]string) ([]E, error)
	Count() (int64, error)

	New() E
}
