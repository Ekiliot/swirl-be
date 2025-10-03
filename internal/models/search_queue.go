package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SearchQueue struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	Username  string    `json:"username" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Связи
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// BeforeCreate хук для GORM
func (sq *SearchQueue) BeforeCreate(tx *gorm.DB) error {
	if sq.ID == uuid.Nil {
		sq.ID = uuid.New()
	}
	return nil
}
