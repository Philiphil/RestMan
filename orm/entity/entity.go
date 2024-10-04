package entity

import (
	"time"

	"gorm.io/gorm"
)

type Entity struct {
	Id          ID `json:"id" orm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `orm:"index"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
}
