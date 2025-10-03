package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"swirl-backend/internal/models"
	"swirl-backend/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatHandler struct {
	db  *gorm.DB
	hub *websocket.Hub
}

func NewChatHandler(db *gorm.DB, hub *websocket.Hub) *ChatHandler {
	return &ChatHandler{db: db, hub: hub}
}

type CreateChatRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required"`
}

func (h *ChatHandler) GetChats(c *gin.Context) {
	userID := c.GetString("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var chats []models.Chat
	query := h.db.Preload("CreatedByUser").
		Preload("Participants.User").
		Joins("JOIN chat_users ON chats.id = chat_users.chat_id").
		Where("chat_users.user_id = ?", userID).
		Order("chats.updated_at DESC")

	if err := query.Offset(offset).Limit(limit).Find(&chats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chats": chats,
		"page":  page,
		"limit": limit,
	})
}

func (h *ChatHandler) CreateChat(c *gin.Context) {
	userID := c.GetString("user_id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создаем чат
	chat := models.Chat{
		Name:        req.Name,
		Description: req.Description,
		Type:        models.ChatType(req.Type),
		CreatedBy:   userUUID,
	}

	if err := h.db.Create(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	// Добавляем создателя как участника
	chatUser := models.ChatUser{
		ChatID:   chat.ID,
		UserID:   userUUID,
		IsActive: true,
	}

	if err := h.db.Create(&chatUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to chat"})
		return
	}

	// Загружаем связанные данные
	h.db.Preload("CreatedByUser").Preload("Participants.User").First(&chat, chat.ID)

	c.JSON(http.StatusCreated, chat)
}

func (h *ChatHandler) GetChat(c *gin.Context) {
	chatID := c.Param("id")
	userID := c.GetString("user_id")

	// Проверяем, является ли пользователь участником чата
	var chatUser models.ChatUser
	if err := h.db.Where("chat_id = ? AND user_id = ?", chatID, userID).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var chat models.Chat
	if err := h.db.Preload("CreatedByUser").
		Preload("Participants.User").
		Preload("Messages.User").
		First(&chat, "id = ?", chatID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	c.JSON(http.StatusOK, chat)
}

func (h *ChatHandler) DeleteChat(c *gin.Context) {
	chatID := c.Param("id")
	userID := c.GetString("user_id")

	// Проверяем, является ли пользователь создателем чата
	var chat models.Chat
	if err := h.db.Where("id = ? AND created_by = ?", chatID, userID).First(&chat).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.db.Delete(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}

func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
	// Получаем токен из query параметра
	token := c.Query("token")
	chatID := c.Query("chat_id")

	if token == "" || chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token and chat_id parameters required"})
		return
	}

	// Проверяем токен
	userID, err := h.validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Проверяем, является ли пользователь участником чата
	var chatUser models.ChatUser
	if err := h.db.Where("chat_id = ? AND user_id = ?", chatID, userID).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	websocket.HandleWebSocket(h.hub, c.Writer, c.Request, userID, chatID)
}

func (h *ChatHandler) validateToken(tokenString string) (string, error) {
	// Убираем "Bearer " префикс если есть
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	// Парсим JWT токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-super-secret-jwt-key-change-in-production"), nil // Используем тот же секрет что в auth.go
	})

	if err != nil || !token.Valid {
		return "", err
	}

	// Извлекаем user_id из claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrInvalidKey
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", jwt.ErrInvalidKey
	}

	return userID, nil
}
