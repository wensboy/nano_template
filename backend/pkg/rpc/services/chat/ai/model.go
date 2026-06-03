package ai

import "gorm.io/gorm"

type AISession struct {
	gorm.Model        // Embeds ID, CreatedAt, UpdatedAt, DeletedAt for database meta fields
	UserID     uint   `json:"user_id" gorm:"column:user_id;not null" binding:"required"`
	SessionID  []byte `json:"session_id" gorm:"column:session_id;type:binary(16);not null" binding:"required"`
	Title      string `json:"title" gorm:"column:title;not null" binding:"required"`
}

func (AISession) TableName() string {
	return "ai_sessions"
}

type AIMessage struct {
	gorm.Model         // Embeds ID, CreatedAt, UpdatedAt, DeletedAt for database meta fields
	SessionID   uint   `json:"session_id" gorm:"column:session_id;not null" binding:"required"`
	UserMessage string `json:"user_message" gorm:"column:user_message;not null" binding:"required"`
	AIMessage   string `json:"ai_message" gorm:"column:ai_message;not null" binding:"required"`
	AIModel     string `json:"ai_model" gorm:"column:ai_model;not null" binding:"required"`
}

func (AIMessage) TableName() string {
	return "ai_messages"
}

type ChatProfile struct {
	gorm.Model          // Embeds ID, CreatedAt, UpdatedAt, DeletedAt for database meta fields
	UserID      uint    `json:"user_id" gorm:"column:user_id;not null" binding:"required"`
	Stream      bool    `json:"stream" gorm:"column:stream;not null"`
	Think       bool    `json:"think" gorm:"column:think;not null"`
	Temperature float64 `json:"temperature" gorm:"column:temperature;type:decimal(10,2);not null" binding:"required"`
	Render      bool    `json:"render" gorm:"column:render;not null"`
	RenderRate  int     `json:"render_rate" gorm:"column:render_rate;not null" binding:"required"`
	RenderJit   bool    `json:"render_jit" gorm:"column:render_jit;not null"`
	MaxToken    int64   `json:"max_token" gorm:"column:max_token;not null" binding:"required"`
}

func (ChatProfile) TableName() string {
	return "chat_profiles"
}

type ChatAppendix struct {
	gorm.Model        // Embeds ID, CreatedAt, UpdatedAt, DeletedAt for database meta fields
	MessageID  uint   `json:"message_id" gorm:"column:message_id;not null" binding:"required"`
	Mime       string `json:"mime" gorm:"column:mime;not null" binding:"required"`
	Name       string `json:"name" gorm:"column:name;not null" binding:"required"`
	Link       string `json:"link" gorm:"column:link;not null" binding:"required"`
	Size       int64  `json:"size" gorm:"column:size;not null" binding:"required"`
	Source     int    `json:"source" gorm:"column:source;not null" binding:"required"`
}

func (ChatAppendix) TableName() string {
	return "chat_appendixs"
}
