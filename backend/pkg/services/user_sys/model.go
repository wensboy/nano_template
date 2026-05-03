package userSys

import (
	"gorm.io/gorm"
)

// User represents a user in the system.
type User struct {
	gorm.Model        // Embeds ID, CreatedAt, UpdatedAt, DeletedAt for database meta fields
	Username   string `json:"username" gorm:"column:username;not null" binding:"required"`
	Password   string `json:"password" gorm:"column:password;not null" binding:"required"`
}

// UserProfile represents a user's profile in the system.
type UserProfile struct {
	gorm.Model        // Embeds ID, CreatedAt, UpdatedAt, DeletedAt for database meta fields
	UserID     uint   `json:"user_id" gorm:"column:user_id;not null" binding:"required"`
	Avatar     string `json:"avatar" gorm:"column:avatar"`
	Nickname   string `json:"nickname" gorm:"column:nickname"`
	Gender     string `json:"gender" gorm:"column:gender"`
	Email      string `json:"email" gorm:"column:email;unique"`
	Signature  string `json:"signature" gorm:"column:signature"`
}

type UserWithProfile struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Signature string `json:"signature"`
}
