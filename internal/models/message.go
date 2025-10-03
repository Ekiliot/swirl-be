package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeSticker  MessageType = "sticker"
	MessageTypeGif      MessageType = "gif"
	MessageTypeVoice    MessageType = "voice"
	MessageTypeVideo    MessageType = "video"
	MessageTypeImage    MessageType = "image"
	MessageTypeReply    MessageType = "reply"
)

type MessageStatus string

const (
	MessageStatusSent     MessageStatus = "sent"     // Отправлено
	MessageStatusDelivered MessageStatus = "delivered" // Доставлено
	MessageStatusRead     MessageStatus = "read"     // Прочитано
	MessageStatusEdited   MessageStatus = "edited"   // Изменено
	MessageStatusDeleted  MessageStatus = "deleted"  // Удалено
)

type Message struct {
	ID        uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ChatID    uuid.UUID     `json:"chat_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID     `json:"user_id" gorm:"type:uuid;not null"`
	Type      MessageType   `json:"type" gorm:"not null"`
	Content   string        `json:"content" gorm:"not null"`
	MediaURL  string        `json:"media_url,omitempty"`
	ReplyToID *uuid.UUID    `json:"reply_to_id,omitempty" gorm:"type:uuid"`
	
	// Статус сообщения
	Status      MessageStatus `json:"status" gorm:"default:'sent'"`
	IsEdited    bool          `json:"is_edited" gorm:"default:false"`
	EditedAt    *time.Time    `json:"edited_at,omitempty"`
	ReadAt      *time.Time    `json:"read_at,omitempty"`
	ReadBy      []uuid.UUID   `json:"read_by,omitempty" gorm:"type:uuid[]"` // Кто прочитал сообщение
	
	// Лайки
	LikesCount int         `json:"likes_count" gorm:"default:0"` // Количество лайков
	LikedBy    []uuid.UUID `json:"liked_by,omitempty" gorm:"type:uuid[]"` // Кто лайкнул
	LikedAt    *time.Time  `json:"liked_at,omitempty"` // Время последнего лайка
	
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`

	// Связи
	Chat    Chat     `json:"chat" gorm:"foreignKey:ChatID"`
	User    User     `json:"user" gorm:"foreignKey:UserID"`
	ReplyTo *Message `json:"reply_to,omitempty" gorm:"foreignKey:ReplyToID"`
}

// BeforeCreate хук для GORM
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	// Устанавливаем статус "отправлено" по умолчанию
	if m.Status == "" {
		m.Status = MessageStatusSent
	}
	return nil
}

// MarkAsDelivered отмечает сообщение как доставленное
func (m *Message) MarkAsDelivered() {
	m.Status = MessageStatusDelivered
}

// MarkAsRead отмечает сообщение как прочитанное пользователем
func (m *Message) MarkAsRead(userID uuid.UUID) {
	// Проверяем, не прочитал ли уже этот пользователь
	for _, id := range m.ReadBy {
		if id == userID {
			return // Уже прочитано этим пользователем
		}
	}
	
	// Добавляем пользователя в список прочитавших
	m.ReadBy = append(m.ReadBy, userID)
	
	// Если это первое прочтение, обновляем время
	if m.ReadAt == nil {
		now := time.Now()
		m.ReadAt = &now
	}
	
	// Если все участники чата прочитали, меняем статус
	// TODO: Нужно получить список участников чата для полной проверки
	m.Status = MessageStatusRead
}

// MarkAsEdited отмечает сообщение как измененное
func (m *Message) MarkAsEdited() {
	m.Status = MessageStatusEdited
	m.IsEdited = true
	now := time.Now()
	m.EditedAt = &now
}

// MarkAsDeleted отмечает сообщение как удаленное
func (m *Message) MarkAsDeleted() {
	m.Status = MessageStatusDeleted
}

// GetStatusText возвращает текстовое описание статуса
func (m *Message) GetStatusText() string {
	switch m.Status {
	case MessageStatusSent:
		return "Отправлено"
	case MessageStatusDelivered:
		return "Доставлено"
	case MessageStatusRead:
		return "Прочитано"
	case MessageStatusEdited:
		return "Изменено"
	case MessageStatusDeleted:
		return "Удалено"
	default:
		return "Неизвестно"
	}
}

// IsReadBy проверяет, прочитал ли сообщение конкретный пользователь
func (m *Message) IsReadBy(userID uuid.UUID) bool {
	for _, id := range m.ReadBy {
		if id == userID {
			return true
		}
	}
	return false
}

// LikeMessage добавляет лайк от пользователя
func (m *Message) LikeMessage(userID uuid.UUID) bool {
	// Проверяем, не лайкнул ли уже этот пользователь
	for _, id := range m.LikedBy {
		if id == userID {
			return false // Уже лайкнул
		}
	}
	
	// Добавляем лайк
	m.LikedBy = append(m.LikedBy, userID)
	m.LikesCount++
	now := time.Now()
	m.LikedAt = &now
	
	return true
}

// UnlikeMessage убирает лайк от пользователя
func (m *Message) UnlikeMessage(userID uuid.UUID) bool {
	// Ищем пользователя в списке лайкнувших
	for i, id := range m.LikedBy {
		if id == userID {
			// Удаляем из списка
			m.LikedBy = append(m.LikedBy[:i], m.LikedBy[i+1:]...)
			m.LikesCount--
			
			// Обновляем время последнего лайка
			if len(m.LikedBy) > 0 {
				now := time.Now()
				m.LikedAt = &now
			} else {
				m.LikedAt = nil
			}
			
			return true
		}
	}
	
	return false // Не лайкал
}

// IsLikedBy проверяет, лайкнул ли сообщение конкретный пользователь
func (m *Message) IsLikedBy(userID uuid.UUID) bool {
	for _, id := range m.LikedBy {
		if id == userID {
			return true
		}
	}
	return false
}

// GetLikesInfo возвращает информацию о лайках
func (m *Message) GetLikesInfo() map[string]interface{} {
	return map[string]interface{}{
		"likes_count": m.LikesCount,
		"liked_by":    m.LikedBy,
		"liked_at":    m.LikedAt,
	}
}
