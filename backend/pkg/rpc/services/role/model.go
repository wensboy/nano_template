package role

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model         // Embeds ID, CreatedAt, UpdatedAt, DeletedAt for database meta fields
	Name        string `json:"name" gorm:"column:name;not null" binding:"required"`
	Level       int    `json:"level" gorm:"column:level;not null" binding:"required"`
	State       uint   `json:"state" gorm:"column:state;not null" binding:"required"`
	Description string `json:"description" gorm:"column:description"`
}
