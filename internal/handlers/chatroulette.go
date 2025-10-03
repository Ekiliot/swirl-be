package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"swirl-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatrouletteHandler struct {
	db            *gorm.DB
	searchHandler *SearchQueueHandler
}

func NewChatrouletteHandler(db *gorm.DB) *ChatrouletteHandler {
	return &ChatrouletteHandler{
		db:            db,
		searchHandler: NewSearchQueueHandler(db),
	}
}

type FindRandomUserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	ChatID   string `json:"chat_id"`
}

func (h *ChatrouletteHandler) FindRandomUser(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	currentUserUUID, err := uuid.Parse(currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Получаем информацию о текущем пользователе
	var currentUser models.User
	if err := h.db.Where("id = ?", currentUserUUID).First(&currentUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Добавляем текущего пользователя в очередь поиска
	if err := h.searchHandler.AddToQueue(currentUserID, currentUser.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to search queue"})
		return
	}

	// Очищаем неактивных пользователей из очереди
	h.searchHandler.CleanupInactiveUsers()

	// Ищем случайного пользователя в очереди поиска, исключая тех, с кем уже есть сохраненный чат
	randomUser, err := h.findRandomUserWithoutExistingChat(currentUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users available in search queue"})
		return
	}

	// Создаем временный чат для чатрулетки
	chat := models.Chat{
		Name:        "Chatroulette Chat",
		Description: "Temporary chat for chatroulette",
		Type:        models.ChatTypePrivate,
		CreatedBy:   currentUserUUID,
	}

	if err := h.db.Create(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	// Добавляем обоих пользователей в чат
	chatUsers := []models.ChatUser{
		{
			ChatID:   chat.ID,
			UserID:   currentUserUUID,
			IsActive: true,
		},
		{
			ChatID:   chat.ID,
			UserID:   randomUser.UserID,
			IsActive: true,
		},
	}

	if err := h.db.Create(&chatUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add users to chat"})
		return
	}

	// Удаляем обоих пользователей из очереди поиска
	h.searchHandler.RemoveFromQueue(currentUserID)
	h.searchHandler.RemoveFromQueue(randomUser.UserID.String())

	c.JSON(http.StatusOK, FindRandomUserResponse{
		UserID:   randomUser.UserID.String(),
		Username: randomUser.Username,
		ChatID:   chat.ID.String(),
	})
}

func (h *ChatrouletteHandler) SaveChat(c *gin.Context) {
	chatID := c.Param("id")
	userID := c.GetString("user_id")

	// Проверяем, является ли пользователь участником чата
	var chatUser models.ChatUser
	if err := h.db.Where("chat_id = ? AND user_id = ?", chatID, userID).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Получаем чат с участниками
	var chat models.Chat
	if err := h.db.Preload("Participants").First(&chat, "id = ?", chatID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	// Находим ID другого пользователя в чате
	var otherUserID string
	for _, participant := range chat.Participants {
		if participant.UserID.String() != userID {
			otherUserID = participant.UserID.String()
			break
		}
	}

	// Проверяем, нет ли уже сохраненного чата между этими пользователями
	if otherUserID != "" {
		hasExistingChat, err := h.searchHandler.CheckExistingChat(userID, otherUserID)
		if err == nil && hasExistingChat {
			c.JSON(http.StatusConflict, gin.H{"error": "A saved chat already exists between these users"})
			return
		}
	}

	// Меняем тип чата на "saved"
	chat.Type = models.ChatTypeSaved
	chat.Name = "Saved Chat - " + chat.Name

	if err := h.db.Save(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Chat saved successfully",
		"chat":    chat,
	})
}

func (h *ChatrouletteHandler) SkipUser(c *gin.Context) {
	chatID := c.Param("id")
	userID := c.GetString("user_id")

	// Проверяем, является ли пользователь участником чата
	var chatUser models.ChatUser
	if err := h.db.Where("chat_id = ? AND user_id = ?", chatID, userID).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Деактивируем пользователя в чате
	chatUser.IsActive = false
	if err := h.db.Save(&chatUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to skip user"})
		return
	}

	// Если это временный чат чатрулетки, удаляем его через некоторое время
	go func() {
		time.Sleep(5 * time.Minute) // Ждем 5 минут
		
		var chat models.Chat
		if err := h.db.First(&chat, "id = ?", chatID).Error; err == nil {
			// Проверяем, активен ли кто-то еще в чате
			var activeUsers int64
			h.db.Model(&models.ChatUser{}).Where("chat_id = ? AND is_active = ?", chatID, true).Count(&activeUsers)
			
			if activeUsers == 0 {
				// Удаляем чат, если никто не активен
				h.db.Delete(&chat)
			}
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "User skipped"})
}

// UpdateSearchActivity обновляет активность пользователя в очереди поиска
func (h *ChatrouletteHandler) UpdateSearchActivity(c *gin.Context) {
	userID := c.GetString("user_id")
	
	if err := h.searchHandler.UpdateUserActivity(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Activity updated"})
}

// GetQueueStatus возвращает статус очереди поиска
func (h *ChatrouletteHandler) GetQueueStatus(c *gin.Context) {
	size, err := h.searchHandler.GetQueueSize()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get queue size"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"queue_size": size,
		"message": "Queue status retrieved",
	})
}

// ClearQueue очищает всю очередь поиска (административная функция)
func (h *ChatrouletteHandler) ClearQueue(c *gin.Context) {
	if err := h.searchHandler.ClearAllQueue(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear queue"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Queue cleared successfully",
	})
}

// findRandomUserWithoutExistingChat находит случайного пользователя, с которым еще нет сохраненного чата
func (h *ChatrouletteHandler) findRandomUserWithoutExistingChat(currentUserID string) (*models.SearchQueue, error) {
	// Получаем всех пользователей в очереди, кроме текущего
	var allUsers []models.SearchQueue
	if err := h.db.Where("user_id != ?", currentUserID).Find(&allUsers).Error; err != nil {
		return nil, err
	}

	// Фильтруем пользователей, с которыми уже есть сохраненный чат
	var availableUsers []models.SearchQueue
	for _, user := range allUsers {
		hasExistingChat, err := h.searchHandler.CheckExistingChat(currentUserID, user.UserID.String())
		if err != nil {
			continue // Пропускаем в случае ошибки
		}
		
		// Добавляем отладочную информацию
		if hasExistingChat {
			fmt.Printf("Пользователь %s исключен - уже есть сохраненный чат с %s\n", user.Username, currentUserID)
		} else {
			fmt.Printf("Пользователь %s доступен для чатрулетки\n", user.Username)
		}
		
		if !hasExistingChat {
			availableUsers = append(availableUsers, user)
		}
	}

	if len(availableUsers) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Возвращаем случайного пользователя из доступных
	randomIndex := rand.Intn(len(availableUsers))
	return &availableUsers[randomIndex], nil
}
