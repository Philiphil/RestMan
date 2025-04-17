package entity

import (
	"time"

	"gorm.io/gorm"
)

// BaseEntity is a default implementation of Entity interface
type BaseEntity struct {
	Id          ID `json:"id" orm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `orm:"index"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
}

func (e BaseEntity) GetId() ID {
	return e.Id
}

func (e BaseEntity) SetId(id any) Entity {
	e.Id = CastId(id)
	return e
}
