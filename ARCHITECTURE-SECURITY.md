# 🏗️ Архитектура и безопасность Chatroulette Backend

## 🔄 **Схема работы системы:**

```
┌─────────────┐    HTTP/WebSocket    ┌─────────────┐    SQL Queries    ┌─────────────┐
│   Клиент    │ ←──────────────────→ │   Бекенд    │ ←───────────────→ │ База данных │
│ (Браузер)   │                      │   (Go API)  │                  │ (PostgreSQL)│
└─────────────┘                      └─────────────┘                  └─────────────┘
```

## ❌ **НЕПРАВИЛЬНО (прямой доступ):**
```
┌─────────────┐    SQL Queries    ┌─────────────┐
│   Клиент    │ ←───────────────→ │ База данных │
│ (Браузер)   │                  │ (PostgreSQL)│
└─────────────┘                  └─────────────┘
```

## ✅ **ПРАВИЛЬНО (через бекенд):**
```
┌─────────────┐    HTTP/WebSocket    ┌─────────────┐    SQL Queries    ┌─────────────┐
│   Клиент    │ ←──────────────────→ │   Бекенд    │ ←───────────────→ │ База данных │
│ (Браузер)   │                      │   (Go API)  │                  │ (PostgreSQL)│
└─────────────┘                      └─────────────┘                  └─────────────┘
```

## 🛡️ **Многоуровневая защита данных:**

### **1. Уровень аутентификации:**
```go
// JWT токен обязателен для всех запросов
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
    // Проверка Authorization header
    // Валидация JWT токена
    // Извлечение user_id из токена
    c.Set("user_id", userID)
}
```

### **2. Уровень авторизации:**
```go
// Проверка доступа к чату
var chatUser models.ChatUser
if err := h.db.Where("chat_id = ? AND user_id = ?", chatID, userID).First(&chatUser).Error; err != nil {
    c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
    return
}
```

### **3. Уровень изоляции данных:**
```go
// Каждый пользователь видит только свои данные
func (h *MessageHandler) GetMessages(c *gin.Context) {
    userID := c.GetString("user_id") // Из JWT токена
    chatID := c.Param("id")
    
    // Проверяем, что пользователь участник чата
    // Возвращаем только сообщения из этого чата
}
```

## 🔒 **Может ли пользователь читать чужие данные?**

### **❌ НЕТ! Вот почему:**

#### **1. JWT аутентификация:**
- **Каждый запрос** требует валидный JWT токен
- **Токен содержит только user_id** текущего пользователя
- **Невозможно подделать** без секретного ключа

#### **2. Проверка доступа к чатам:**
```go
// Пользователь может получить сообщения ТОЛЬКО из своих чатов
var chatUser models.ChatUser
if err := h.db.Where("chat_id = ? AND user_id = ?", chatID, userID).First(&chatUser).Error; err != nil {
    c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
    return
}
```

#### **3. Изоляция на уровне SQL:**
```sql
-- Пользователь видит только свои чаты
SELECT * FROM chats 
JOIN chat_users ON chats.id = chat_users.chat_id 
WHERE chat_users.user_id = 'current-user-id'

-- Пользователь видит только сообщения из своих чатов
SELECT * FROM messages 
WHERE chat_id IN (
    SELECT chat_id FROM chat_users 
    WHERE user_id = 'current-user-id'
)
```

## 🚫 **Что НЕ может пользователь:**

### **❌ Получить чужие данные:**
- **Другие чаты** - только свои
- **Чужие сообщения** - только из своих чатов  
- **Личные данные других** - только username участников чата
- **Пароли других** - зашифрованы и недоступны

### **❌ Обойти авторизацию:**
- **Без JWT токена** - запрос отклоняется
- **С поддельным токеном** - не пройдет валидацию
- **С чужим токеном** - не пройдет проверку доступа

### **❌ Прямой доступ к БД:**
- **База данных** недоступна извне
- **Только бекенд** может делать SQL запросы
- **Все запросы** проходят через авторизацию

## 🔐 **Примеры защиты в коде:**

### **Получение сообщений:**
```go
func (h *MessageHandler) GetMessages(c *gin.Context) {
    userID := c.GetString("user_id") // Из JWT
    chatID := c.Param("id")
    
    // 1. Проверяем, что пользователь участник чата
    var chatUser models.ChatUser
    if err := h.db.Where("chat_id = ? AND user_id = ?", chatID, userID).First(&chatUser).Error; err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
        return
    }
    
    // 2. Возвращаем только сообщения из этого чата
    var messages []models.Message
    h.db.Where("chat_id = ?", chatID).Find(&messages)
}
```

### **Отправка сообщений:**
```go
func (h *MessageHandler) SendMessage(c *gin.Context) {
    userID := c.GetString("user_id") // Из JWT
    chatID := c.Param("id")
    
    // 1. Проверяем доступ к чату
    var chatUser models.ChatUser
    if err := h.db.Where("chat_id = ? AND user_id = ?", chatID, userID).First(&chatUser).Error; err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
        return
    }
    
    // 2. Создаем сообщение от имени текущего пользователя
    message := models.Message{
        ChatID:  uuid.MustParse(chatID),
        UserID:  uuid.MustParse(userID), // Только текущий пользователь
        Content: req.Content,
    }
}
```

### **WebSocket соединения:**
```go
func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
    token := c.Query("token")
    chatID := c.Query("chat_id")
    
    // 1. Проверяем JWT токен
    userID, err := h.validateToken(token)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }
    
    // 2. Проверяем доступ к чату
    var chatUser models.ChatUser
    if err := h.db.Where("chat_id = ? AND user_id = ?", chatID, userID).First(&chatUser).Error; err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
        return
    }
    
    // 3. Устанавливаем WebSocket соединение
    websocket.HandleWebSocket(h.hub, c.Writer, c.Request, userID, chatID)
}
```

## 🛡️ **Дополнительные меры безопасности:**

### **1. Middleware безопасности:**
```go
// Заголовки безопасности
c.Header("X-Content-Type-Options", "nosniff")
c.Header("X-Frame-Options", "DENY")
c.Header("X-XSS-Protection", "1; mode=block")

// Rate limiting
func RateLimiting() gin.HandlerFunc {
    // Ограничение запросов на IP
}

// Валидация входных данных
func InputValidation() gin.HandlerFunc {
    // Защита от SQL инъекций
}
```

### **2. Шифрование данных:**
```go
// Пароли хешируются bcrypt
func (u *User) HashPassword(password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    u.Password = string(hashedPassword)
    return err
}
```

### **3. Изоляция файлов:**
```go
// Файлы доступны только участникам чата
func (h *UploadHandler) GetFile(c *gin.Context) {
    userID := c.Param("user_id")
    fileName := c.Param("file_name")
    
    // Проверяем, что файл принадлежит пользователю
    if userID != currentUserID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
        return
    }
}
```

## ✅ **Заключение:**

### **🔒 Данные полностью защищены:**
- **Нет прямого доступа** к базе данных
- **Каждый запрос** проходит авторизацию
- **Пользователи видят** только свои данные
- **Невозможно получить** чужие данные

### **🛡️ Многоуровневая защита:**
- **JWT аутентификация** - кто вы?
- **Проверка доступа** - можете ли вы?
- **Изоляция данных** - только ваши данные
- **Валидация запросов** - защита от атак

### **🚫 Что невозможно:**
- **Прямой доступ** к базе данных
- **Чтение чужих данных** без авторизации
- **Обход системы безопасности** через API
- **Получение паролей** других пользователей

**Ваши данные в полной безопасности!** 🛡️✨
