package security

import (
	"github.com/philiphil/apiman/orm/entity"
)

func NoAuthorizationRequired(user IUser, object entity.IEntity) bool {
	return true
}

type AuthorizationFunction func(IUser, entity.IEntity) bool
