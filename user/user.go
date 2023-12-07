package user

import (
	"github.com/philiphil/apiman/orm/entity"
)

type IUser interface {
	SetId(any) entity.IEntity
	GetId() entity.ID
}

type Auth interface {
}

type UserRepository interface {
	GetUser(Auth) (IUser, error)
}

type SecurityManager interface {
	GetUserRepository() UserRepository
}
