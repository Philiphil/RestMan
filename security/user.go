package security

import (
	"github.com/philiphil/restman/orm/entity"
)

type IUser interface {
	SetId(any) entity.IEntity
	GetId() entity.ID
}
