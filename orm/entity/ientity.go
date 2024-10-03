package entity

func (e Entity) GetId() ID {
	return e.Id
}

func (e Entity) SetId(id any) IEntity {
	e.Id = CastId(id)
	return e
}

type IEntity interface {
	SetId(any) IEntity
	GetId() ID
}
