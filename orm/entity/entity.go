package entity

// Entity is an interface that represents a database entity
type Entity interface {
	SetId(any) Entity
	GetId() ID
}
