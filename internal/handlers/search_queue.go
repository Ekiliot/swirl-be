package handlers

import (
	"fmt"
	"time"

	"swirl-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SearchQueueHandler struct {
	db *gorm.DB
}

func NewSearchQueueHandler(db *gorm.DB) *SearchQueueHandler {
	return &SearchQueueHandler{db: db}
}

// AddToQueue добавляет пользователя в очередь поиска
func (h *SearchQueueHandler) AddToQueue(userID, username string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	// Проверяем, есть ли пользователь уже в очереди
	var existingUser models.SearchQueue
	if err := h.db.Where("user_id = ?", userUUID).First(&existingUser).Error; err == nil {
		// Пользователь уже в очереди, обновляем время активности
		fmt.Printf("Пользователь %s уже в очереди, обновляем активность\n", username)
		return h.db.Model(&existingUser).Update("updated_at", time.Now()).Error
	}

	// Удаляем пользователя из очереди если он уже там есть (дополнительная защита)
	h.db.Where("user_id = ?", userUUID).Delete(&models.SearchQueue{})

	// Добавляем пользователя в очередь
	searchQueue := models.SearchQueue{
		UserID:   userUUID,
		Username: username,
	}

	fmt.Printf("Добавляем пользователя %s в очередь поиска\n", username)
	return h.db.Create(&searchQueue).Error
}

// RemoveFromQueue удаляет пользователя из очереди поиска
func (h *SearchQueueHandler) RemoveFromQueue(userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	return h.db.Where("user_id = ?", userUUID).Delete(&models.SearchQueue{}).Error
}

// FindRandomUserInQueue находит случайного пользователя в очереди поиска
func (h *SearchQueueHandler) FindRandomUserInQueue(currentUserID string) (*models.SearchQueue, error) {
	currentUserUUID, err := uuid.Parse(currentUserID)
	if err != nil {
		return nil, err
	}

	var searchQueue models.SearchQueue
	err = h.db.Where("user_id != ?", currentUserUUID).Order("RANDOM()").First(&searchQueue).Error
	if err != nil {
		return nil, err
	}

	return &searchQueue, nil
}

// CleanupInactiveUsers удаляет неактивных пользователей из очереди (старше 1 минуты)
func (h *SearchQueueHandler) CleanupInactiveUsers() error {
	oneMinuteAgo := time.Now().Add(-time.Minute)
	
	// Сначала получаем количество пользователей для удаления
	var count int64
	h.db.Model(&models.SearchQueue{}).Where("updated_at < ?", oneMinuteAgo).Count(&count)
	
	if count > 0 {
		fmt.Printf("Очистка очереди: удаляем %d неактивных пользователей\n", count)
	}
	
	return h.db.Where("updated_at < ?", oneMinuteAgo).Delete(&models.SearchQueue{}).Error
}

// UpdateUserActivity обновляет время активности пользователя в очереди
func (h *SearchQueueHandler) UpdateUserActivity(userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	return h.db.Model(&models.SearchQueue{}).Where("user_id = ?", userUUID).Update("updated_at", time.Now()).Error
}

// GetQueueSize возвращает количество пользователей в очереди
func (h *SearchQueueHandler) GetQueueSize() (int64, error) {
	var count int64
	err := h.db.Model(&models.SearchQueue{}).Count(&count).Error
	return count, err
}

// ClearAllQueue очищает всю очередь (для административных целей)
func (h *SearchQueueHandler) ClearAllQueue() error {
	var count int64
	h.db.Model(&models.SearchQueue{}).Count(&count)
	
	if count > 0 {
		fmt.Printf("Полная очистка очереди: удаляем %d пользователей\n", count)
	}
	
	return h.db.Delete(&models.SearchQueue{}, "1 = 1").Error
}

// CheckExistingChat проверяет, есть ли уже сохраненный чат между пользователями
func (h *SearchQueueHandler) CheckExistingChat(userID1, userID2 string) (bool, error) {
	userUUID1, err := uuid.Parse(userID1)
	if err != nil {
		return false, err
	}
	
	userUUID2, err := uuid.Parse(userID2)
	if err != nil {
		return false, err
	}

	// Ищем сохраненные чаты между этими пользователями
	var count int64
	err = h.db.Model(&models.Chat{}).
		Joins("JOIN chat_users cu1 ON chats.id = cu1.chat_id").
		Joins("JOIN chat_users cu2 ON chats.id = cu2.chat_id").
		Where("chats.type = ? AND ((cu1.user_id = ? AND cu2.user_id = ?) OR (cu1.user_id = ? AND cu2.user_id = ?))", 
			models.ChatTypeSaved, userUUID1, userUUID2, userUUID2, userUUID1).
		Count(&count).Error

	return count > 0, err
}
