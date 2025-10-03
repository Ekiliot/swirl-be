package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatType string

const (
	ChatTypePrivate ChatType = "private"
	ChatTypeGroup   ChatType = "group"
	ChatTypeSaved   ChatType = "saved" // Сохраненный чат из чатрулетки
)

type Chat struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Type        ChatType  `json:"type" gorm:"not null"`
	Description string    `json:"description"`
	CreatedBy   uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Связи
	CreatedByUser User           `json:"created_by_user" gorm:"foreignKey:CreatedBy"`
	Participants  []ChatUser     `json:"participants" gorm:"foreignKey:ChatID"`
	Messages      []Message      `json:"messages" gorm:"foreignKey:ChatID"`
}

type ChatUser struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ChatID    uuid.UUID `json:"chat_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	JoinedAt  time.Time `json:"joined_at"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`

	// Связи
	Chat Chat `json:"chat" gorm:"foreignKey:ChatID"`
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// BeforeCreate хук для GORM
func (c *Chat) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (cu *ChatUser) BeforeCreate(tx *gorm.DB) error {
	if cu.ID == uuid.Nil {
		cu.ID = uuid.New()
	}
	return nil
}
