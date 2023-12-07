package user

import (
	"github.com/philiphil/apiman/orm/entity"
)

type TestUser struct {
	ID entity.ID
}

func (t TestUser) HasReadingRight(entity entity.IEntity) bool {
	//TODO implement me
	panic("implement me")
}

func (t TestUser) HasWritingRight(entity entity.IEntity) bool {
	//TODO implement me
	panic("implement me")
}

func (t TestUser) SetId(a any) entity.IEntity {
	//TODO implement me
	panic("implement me")
}

func (t TestUser) GetId() entity.ID {
	//TODO implement me
	panic("implement me")
}

type TestUserRepository struct {
}

func (t TestUserRepository) GetUser(auth Auth) (IUser, error) {
	return TestUser{ID: 1}, nil
}

type TestSecurity struct {
	TestUserRepository TestUserRepository
}

func (t TestSecurity) GetUserRepository() UserRepository {
	return t.TestUserRepository
}
