package security

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/orm/entity"
)

type AnonymousUser struct {
	ID entity.ID
}

func (t AnonymousUser) SetId(a any) entity.IEntity {
	t.ID = entity.CastId(a)
	return t
}

func (t AnonymousUser) GetId() entity.ID {
	return t.ID
}

type AnonymousFirewall struct{}

func (a AnonymousFirewall) GetUser(c *gin.Context) (IUser, error) {
	return AnonymousUser{entity.ID(1)}, nil
}
