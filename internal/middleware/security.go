package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders добавляет заголовки безопасности
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Защита от XSS
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// HTTPS только в продакшене
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		
		c.Next()
	}
}

// RateLimiting простая защита от спама
func RateLimiting() gin.HandlerFunc {
	// В продакшене использовать Redis или более продвинутые решения
	return func(c *gin.Context) {
		// Простая проверка User-Agent
		userAgent := c.GetHeader("User-Agent")
		if userAgent == "" || len(userAgent) < 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// InputValidation валидация входных данных
func InputValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверка на SQL инъекции в query параметрах
		query := c.Request.URL.RawQuery
		if strings.Contains(strings.ToLower(query), "union") ||
		   strings.Contains(strings.ToLower(query), "select") ||
		   strings.Contains(strings.ToLower(query), "drop") ||
		   strings.Contains(strings.ToLower(query), "delete") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			c.Abort()
			return
		}
		
		c.Next()
	}
}
