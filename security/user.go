package security

import (
	"github.com/philiphil/restman/orm/entity"
)

// User is an interface for the user
// its compatible with entity.Entity
type User interface {
	SetId(any) entity.Entity
	GetId() entity.ID
}
