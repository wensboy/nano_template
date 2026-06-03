package user

import (
	"gorm.io/gorm"
)

// User represents a user in the system.
type User struct {
	gorm.Model // Embeds ID, CreatedAt, UpdatedAt, DeletedAt for database meta fields
	// username 只是标识用户唯一性的字段, 就像该表的主键一样, 通常可以是: 手机, 身份证件, 邮箱, 按照特殊规则
	// 生成的用户名
	Username string `json:"username" gorm:"column:username;not null" binding:"required"`
	// password 在数据库中永远不能以明文的形式存储
	Password string `json:"password" gorm:"column:password;not null" binding:"required"`
}

// UserProfile represents a user's profile in the system.
type UserProfile struct {
	gorm.Model      // Embeds ID, CreatedAt, UpdatedAt, DeletedAt for database meta fields
	UserID     uint `json:"user_id" gorm:"column:user_id;not null" binding:"required"`
	// avatar 一定是诸如 uri 样式的资源定位路径
	Avatar string `json:"avatar" gorm:"column:avatar"`
	// nickname 不要求唯一
	Nickname string `json:"nickname" gorm:"column:nickname"`
	// email 不要求唯一
	Email string `json:"email" gorm:"column:email"`
	// Phone 不要求唯一
	Phone     string `json:"phone" gorm:"column:phone"`
	Signature string `json:"signature" gorm:"column:signature"`
}

type UserRole struct {
	gorm.Model
	UserID      uint   `json:"user_id" gorm:"column:user_id;not null"`
	RoleID      uint   `json:"role_id" gorm:"column:role_id;not null"`
	Description string `json:"description" gorm:"column:description"`
	State       uint   `json:"state" gorm:"column:state;not null"`
}
