package entity

type Model[E IEntity] interface {
	ToEntity() E
	FromEntity(entity E) interface{}
}
