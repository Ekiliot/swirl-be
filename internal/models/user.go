package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	
	// Профиль пользователя
	Birthday        *time.Time `json:"birthday,omitempty" gorm:"type:date"`
	ProfilePhoto    string     `json:"profile_photo,omitempty"`
	IsOnline        bool       `json:"is_online" gorm:"default:false"`
	LastSeen        *time.Time `json:"last_seen,omitempty"`
	
	// Настройки конфиденциальности
	ShowUsername    bool `json:"show_username" gorm:"default:true"`
	ShowBirthday    bool `json:"show_birthday" gorm:"default:false"`
	ShowOnlineStatus bool `json:"show_online_status" gorm:"default:true"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HashPassword хеширует пароль
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword проверяет пароль
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// BeforeCreate хук для GORM
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// UpdateOnlineStatus обновляет статус онлайн пользователя
func (u *User) UpdateOnlineStatus(isOnline bool) {
	u.IsOnline = isOnline
	now := time.Now()
	u.LastSeen = &now
}

// GetPublicProfile возвращает публичную информацию о пользователе
func (u *User) GetPublicProfile() map[string]interface{} {
	profile := map[string]interface{}{
		"id": u.ID,
	}
	
	if u.ShowUsername {
		profile["username"] = u.Username
	}
	
	if u.ShowBirthday && u.Birthday != nil {
		profile["birthday"] = u.Birthday
	}
	
	if u.ShowOnlineStatus {
		profile["is_online"] = u.IsOnline
		if u.LastSeen != nil {
			profile["last_seen"] = u.LastSeen
		}
	}
	
	if u.ProfilePhoto != "" {
		profile["profile_photo"] = u.ProfilePhoto
	}
	
	return profile
}

// GetOnlineStatusText возвращает текстовое описание статуса онлайн
func (u *User) GetOnlineStatusText() string {
	if u.IsOnline {
		return "В сети"
	}
	
	if u.LastSeen == nil {
		return "Никогда не был в сети"
	}
	
	now := time.Now()
	diff := now.Sub(*u.LastSeen)
	
	if diff < time.Minute {
		return "Только что был в сети"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		return fmt.Sprintf("Был в сети %d мин. назад", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("Был в сети %d ч. назад", hours)
	} else {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "Был в сети вчера"
		} else if days < 7 {
			return fmt.Sprintf("Был в сети %d дн. назад", days)
		} else {
			return u.LastSeen.Format("Был в сети 2 Jan 2006")
		}
	}
}
