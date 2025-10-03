package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	uploadPath string
}

func NewUploadHandler(uploadPath string) *UploadHandler {
	return &UploadHandler{uploadPath: uploadPath}
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	userID := c.GetString("user_id")
	
	// Получаем файл из формы
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Получаем тип файла
	contentType := header.Header.Get("Content-Type")

	// Проверяем размер файла в зависимости от типа
	maxSize := h.getMaxFileSize(contentType)
	if header.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File too large. Maximum size is %dMB", maxSize/(1024*1024))})
		return
	}
	allowedTypes := []string{
		"image/jpeg", "image/png", "image/gif", "image/webp",
		"video/mp4", "video/webm", "video/quicktime",
		"audio/mpeg", "audio/wav", "audio/ogg", "audio/mp4",
	}

	if !h.isAllowedType(contentType, allowedTypes) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed"})
		return
	}

	// Создаем уникальное имя файла
	fileExt := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s_%s%s", userID, uuid.New().String(), fileExt)
	
	// Создаем папку для пользователя, если её нет
	userDir := filepath.Join(h.uploadPath, userID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	// Путь к файлу
	filePath := filepath.Join(userDir, fileName)

	// Создаем файл
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer dst.Close()

	// Копируем содержимое файла
	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Определяем тип медиа
	mediaType := h.getMediaType(contentType)

	// Получаем дополнительную информацию о файле
	fileInfo := h.getFileInfo(filePath, contentType)

	// Возвращаем информацию о загруженном файле
	c.JSON(http.StatusOK, gin.H{
		"file_id":    fileName,
		"file_url":   fmt.Sprintf("/uploads/%s/%s", userID, fileName),
		"media_type": mediaType,
		"size":       header.Size,
		"created_at": time.Now(),
		"info":       fileInfo,
	})
}

func (h *UploadHandler) GetFile(c *gin.Context) {
	userID := c.Param("user_id")
	fileName := c.Param("file_name")

	// Проверяем безопасность пути
	if strings.Contains(fileName, "..") || strings.Contains(userID, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file path"})
		return
	}

	filePath := filepath.Join(h.uploadPath, userID, fileName)
	
	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Отдаем файл
	c.File(filePath)
}

// DeleteFile удаляет файл
func (h *UploadHandler) DeleteFile(c *gin.Context) {
	userID := c.GetString("user_id")
	fileName := c.Param("file_name")

	// Проверяем безопасность пути
	if strings.Contains(fileName, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file path"})
		return
	}

	filePath := filepath.Join(h.uploadPath, userID, fileName)
	
	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Удаляем файл
	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

func (h *UploadHandler) isAllowedType(contentType string, allowedTypes []string) bool {
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			return true
		}
	}
	return false
}

func (h *UploadHandler) getMediaType(contentType string) string {
	if strings.HasPrefix(contentType, "image/") {
		return "image"
	} else if strings.HasPrefix(contentType, "video/") {
		return "video"
	} else if strings.HasPrefix(contentType, "audio/") {
		return "voice"
	}
	return "unknown"
}

// getMaxFileSize возвращает максимальный размер файла в зависимости от типа
func (h *UploadHandler) getMaxFileSize(contentType string) int64 {
	if strings.HasPrefix(contentType, "image/") {
		return 10 * 1024 * 1024 // 10MB для изображений
	} else if strings.HasPrefix(contentType, "video/") {
		return 100 * 1024 * 1024 // 100MB для видео
	} else if strings.HasPrefix(contentType, "audio/") {
		return 25 * 1024 * 1024 // 25MB для аудио
	}
	return 50 * 1024 * 1024 // 50MB по умолчанию
}

// getFileInfo возвращает дополнительную информацию о файле
func (h *UploadHandler) getFileInfo(filePath, contentType string) map[string]interface{} {
	info := map[string]interface{}{
		"content_type": contentType,
		"file_path":    filePath,
	}

	// Для изображений можно добавить размеры
	if strings.HasPrefix(contentType, "image/") {
		info["is_image"] = true
		// TODO: Добавить получение размеров изображения
	}

	// Для видео можно добавить длительность
	if strings.HasPrefix(contentType, "video/") {
		info["is_video"] = true
		// TODO: Добавить получение длительности видео
	}

	// Для аудио можно добавить длительность
	if strings.HasPrefix(contentType, "audio/") {
		info["is_audio"] = true
		// TODO: Добавить получение длительности аудио
	}

	return info
}
