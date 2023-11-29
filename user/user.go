package user

import "github.com/philiphil/apiman/orm"

type IUser interface {
	HasReadingRight(entity orm.IEntity) bool
	HasWritingRight(entity orm.IEntity) bool
	SetId(any) orm.IEntity
	GetId() any
}

type Auth interface {
}

type UserRepository interface {
	GetUser(Auth) (IUser, error)
}

type SecurityManager interface {
	GetUserRepository() UserRepository
}
