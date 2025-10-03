package handlers

import (
	"net/http"
	"strconv"

	"swirl-backend/internal/models"
	"swirl-backend/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageHandler struct {
	db  *gorm.DB
	hub *websocket.Hub
}

func NewMessageHandler(db *gorm.DB, hub *websocket.Hub) *MessageHandler {
	return &MessageHandler{db: db, hub: hub}
}

type SendMessageRequest struct {
	Type      string `json:"type" binding:"required"`
	Content   string `json:"content" binding:"required"`
	MediaURL  string `json:"media_url,omitempty"`
	ReplyToID string `json:"reply_to_id,omitempty"`
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	chatID := c.Param("id")
	userID := c.GetString("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	// Проверяем доступ к чату
	var chatUser models.ChatUser
	if err := h.db.Where("chat_id = ? AND user_id = ?", chatID, userID).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var messages []models.Message
	query := h.db.Preload("User").
		Preload("ReplyTo").
		Preload("ReplyTo.User").
		Where("chat_id = ?", chatID).
		Order("created_at DESC")

	if err := query.Offset(offset).Limit(limit).Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"page":     page,
		"limit":    limit,
	})
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	chatID := c.Param("id")
	userID := c.GetString("user_id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Проверяем доступ к чату
	var chatUser models.ChatUser
	if err := h.db.Where("chat_id = ? AND user_id = ?", chatID, userID).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создаем сообщение
	message := models.Message{
		ChatID:   uuid.MustParse(chatID),
		UserID:   userUUID,
		Type:     models.MessageType(req.Type),
		Content:  req.Content,
		MediaURL: req.MediaURL,
	}

	// Если есть ответ на сообщение
	if req.ReplyToID != "" {
		replyToUUID, err := uuid.Parse(req.ReplyToID)
		if err == nil {
			message.ReplyToID = &replyToUUID
		}
	}

	if err := h.db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Отмечаем как доставленное
	message.MarkAsDelivered()
	h.db.Save(&message)

	// Загружаем связанные данные
	h.db.Preload("User").Preload("ReplyTo").Preload("ReplyTo.User").First(&message, message.ID)

	// Отправляем сообщение через WebSocket
	websocketMessage := websocket.Message{
		Type:    "new_message",
		ChatID:  chatID,
		Payload: message,
	}
	h.hub.Broadcast <- websocketMessage

	c.JSON(http.StatusCreated, message)
}

func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	messageID := c.Param("id")
	userID := c.GetString("user_id")

	// Находим сообщение
	var message models.Message
	if err := h.db.Where("id = ?", messageID).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Проверяем, является ли пользователь автором сообщения
	if message.UserID.String() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.db.Delete(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	// Уведомляем через WebSocket
	websocketMessage := websocket.Message{
		Type:    "message_deleted",
		ChatID:  message.ChatID.String(),
		Payload: gin.H{"message_id": messageID},
	}
	h.hub.Broadcast <- websocketMessage

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}

// MarkMessageAsRead отмечает сообщение как прочитанное
func (h *MessageHandler) MarkMessageAsRead(c *gin.Context) {
	messageID := c.Param("id")
	userID := c.GetString("user_id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Находим сообщение
	var message models.Message
	if err := h.db.Where("id = ?", messageID).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Проверяем, что пользователь является участником чата
	var chatUser models.ChatUser
	if err := h.db.Where("chat_id = ? AND user_id = ?", message.ChatID, userID).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Отмечаем как прочитанное
	message.MarkAsRead(userUUID)
	
	if err := h.db.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message status"})
		return
	}

	// Отправляем обновление через WebSocket
	websocketMessage := websocket.Message{
		Type:    "message_status_update",
		ChatID:  message.ChatID.String(),
		Payload: gin.H{
			"message_id": message.ID,
			"status":     message.Status,
			"read_by":    message.ReadBy,
			"read_at":    message.ReadAt,
		},
	}
	h.hub.Broadcast <- websocketMessage

	c.JSON(http.StatusOK, gin.H{
		"message_id": message.ID,
		"status":     message.Status,
		"status_text": message.GetStatusText(),
		"read_by":    message.ReadBy,
		"read_at":    message.ReadAt,
	})
}

// EditMessage редактирует сообщение
func (h *MessageHandler) EditMessage(c *gin.Context) {
	messageID := c.Param("id")
	userID := c.GetString("user_id")

	// Находим сообщение
	var message models.Message
	if err := h.db.Where("id = ?", messageID).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Проверяем, что пользователь является автором сообщения
	if message.UserID.String() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own messages"})
		return
	}

	var editData struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&editData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновляем сообщение
	message.Content = editData.Content
	message.MarkAsEdited()

	if err := h.db.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit message"})
		return
	}

	// Отправляем обновление через WebSocket
	websocketMessage := websocket.Message{
		Type:    "message_edited",
		ChatID:  message.ChatID.String(),
		Payload: message,
	}
	h.hub.Broadcast <- websocketMessage

	c.JSON(http.StatusOK, message)
}

// GetMessageStatus возвращает статус сообщения
func (h *MessageHandler) GetMessageStatus(c *gin.Context) {
	messageID := c.Param("id")
	userID := c.GetString("user_id")

	// Находим сообщение
	var message models.Message
	if err := h.db.Where("id = ?", messageID).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Проверяем доступ к чату
	var chatUser models.ChatUser
	if err := h.db.Where("chat_id = ? AND user_id = ?", message.ChatID, userID).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message_id": message.ID,
		"status":     message.Status,
		"status_text": message.GetStatusText(),
		"is_edited":  message.IsEdited,
		"edited_at":  message.EditedAt,
		"read_at":    message.ReadAt,
		"read_by":    message.ReadBy,
		"created_at": message.CreatedAt,
		"updated_at": message.UpdatedAt,
	})
}

// LikeMessage лайкает сообщение
func (h *MessageHandler) LikeMessage(c *gin.Context) {
	messageID := c.Param("id")
	userID := c.GetString("user_id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Находим сообщение
	var message models.Message
	if err := h.db.Where("id = ?", messageID).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Проверяем, что пользователь является участником чата
	var chatUser models.ChatUser
	if err := h.db.Where("chat_id = ? AND user_id = ?", message.ChatID, userID).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Добавляем лайк
	if !message.LikeMessage(userUUID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message already liked by this user"})
		return
	}

	// Сохраняем изменения
	if err := h.db.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like message"})
		return
	}

	// Отправляем обновление через WebSocket
	websocketMessage := websocket.Message{
		Type:    "message_liked",
		ChatID:  message.ChatID.String(),
		Payload: gin.H{
			"message_id": message.ID,
			"likes_count": message.LikesCount,
			"liked_by":    message.LikedBy,
			"liked_at":    message.LikedAt,
		},
	}
	h.hub.Broadcast <- websocketMessage

	c.JSON(http.StatusOK, gin.H{
		"message_id": message.ID,
		"likes_count": message.LikesCount,
		"liked_by":    message.LikedBy,
		"liked_at":    message.LikedAt,
		"is_liked":    true,
	})
}

// UnlikeMessage убирает лайк с сообщения
func (h *MessageHandler) UnlikeMessage(c *gin.Context) {
	messageID := c.Param("id")
	userID := c.GetString("user_id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Находим сообщение
	var message models.Message
	if err := h.db.Where("id = ?", messageID).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Проверяем, что пользователь является участником чата
	var chatUser models.ChatUser
	if err := h.db.Where("chat_id = ? AND user_id = ?", message.ChatID, userID).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Убираем лайк
	if !message.UnlikeMessage(userUUID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message not liked by this user"})
		return
	}

	// Сохраняем изменения
	if err := h.db.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike message"})
		return
	}

	// Отправляем обновление через WebSocket
	websocketMessage := websocket.Message{
		Type:    "message_unliked",
		ChatID:  message.ChatID.String(),
		Payload: gin.H{
			"message_id": message.ID,
			"likes_count": message.LikesCount,
			"liked_by":    message.LikedBy,
			"liked_at":    message.LikedAt,
		},
	}
	h.hub.Broadcast <- websocketMessage

	c.JSON(http.StatusOK, gin.H{
		"message_id": message.ID,
		"likes_count": message.LikesCount,
		"liked_by":    message.LikedBy,
		"liked_at":    message.LikedAt,
		"is_liked":    false,
	})
}

// GetMessageLikes возвращает информацию о лайках сообщения
func (h *MessageHandler) GetMessageLikes(c *gin.Context) {
	messageID := c.Param("id")
	userID := c.GetString("user_id")

	// Находим сообщение
	var message models.Message
	if err := h.db.Where("id = ?", messageID).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Проверяем доступ к чату
	var chatUser models.ChatUser
	if err := h.db.Where("chat_id = ? AND user_id = ?", message.ChatID, userID).First(&chatUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	userUUID, _ := uuid.Parse(userID)
	likesInfo := message.GetLikesInfo()
	likesInfo["is_liked_by_user"] = message.IsLikedBy(userUUID)

	c.JSON(http.StatusOK, likesInfo)
}
