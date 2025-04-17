package entity

// Model is an interface that should be implemented by the Model struct, it is used to convert the Model to an Entity and vice versa
type Model[E Entity] interface {
	ToEntity() E
	FromEntity(entity E) any
}
