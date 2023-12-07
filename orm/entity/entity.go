package entity

import (
	"gorm.io/gorm"
	"time"
)

type Entity struct {
	Id          ID `json:"id" orm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `orm:"index"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
}

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
