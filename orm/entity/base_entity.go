package entity

import (
	"time"

	"gorm.io/gorm"
)

// BaseEntity is a default implementation of Entity interface
type BaseEntity struct {
	Id          ID `orm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `orm:"index"`
	Name        string
	Description string
}

// GetId returns the entity's ID.
func (e BaseEntity) GetId() ID {
	return e.Id
}

// SetId sets the entity's ID and returns the entity.
func (e BaseEntity) SetId(id any) Entity {
	e.Id = CastId(id)
	return e
}
