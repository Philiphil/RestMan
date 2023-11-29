package user

import "github.com/philiphil/apiman/orm"

type TestUser struct {
	ID orm.ID
}

func (t TestUser) HasReadingRight(entity orm.IEntity) bool {
	//TODO implement me
	panic("implement me")
}

func (t TestUser) HasWritingRight(entity orm.IEntity) bool {
	//TODO implement me
	panic("implement me")
}

func (t TestUser) SetId(a any) orm.IEntity {
	//TODO implement me
	panic("implement me")
}

func (t TestUser) GetId() any {
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
