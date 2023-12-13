package security

import (
	"github.com/philiphil/apiman/orm/entity"
)

type IUser interface {
	SetId(any) entity.IEntity
	GetId() entity.ID
}
