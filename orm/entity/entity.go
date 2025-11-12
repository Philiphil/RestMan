package entity

// This interface allows an object to have an ID set and retrieved in a generic way.
// this mechanism allows structs to be used generically in the ORM layer.
type Entity interface {
	SetId(any) Entity
	GetId() ID
}
