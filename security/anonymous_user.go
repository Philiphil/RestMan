package security

import (
	"github.com/philiphil/apiman/orm/entity"
)

type AnonymousUser struct {
	ID entity.ID
}

func (t AnonymousUser) SetId(a any) entity.IEntity {
	panic("implement me")
}

func (t AnonymousUser) GetId() entity.ID {
	panic("implement me")
}
