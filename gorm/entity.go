package gorm

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type ID uint

func (e ID) String() string {
	return strconv.FormatUint(uint64(e), 10)
}

func CastId(id any) ID {
	switch v := id.(type) {
	case ID:
		return v
	case int:
		return ID(v)
	case uint:
		return ID(v)
	case string:
		convertedID, _ := strconv.ParseUint(v, 10, 64)
		return ID(convertedID)
	}

	return ID(0)
}

type Entity struct {
	Id          ID `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
}

func (e Entity) GetId() any {
	return e.Id
}

func (e Entity) SetId(id any) IEntity {
	e.Id = CastId(id)
	return e
}

type IEntity interface {
	SetId(any) IEntity
	GetId() any
}
