# ❤️ Система лайков сообщений в Chatroulette Backend

## ✅ **Что добавлено:**

### **💖 Лайки сообщений:**
- **Лайкнуть сообщение** - добавить лайк
- **Убрать лайк** - удалить свой лайк
- **Счетчик лайков** - количество лайков
- **Список лайкнувших** - кто лайкнул сообщение
- **Время лайка** - когда был поставлен лайк

### **🔍 Дополнительные поля:**
- **LikesCount** - количество лайков
- **LikedBy** - массив ID пользователей, лайкнувших сообщение
- **LikedAt** - время последнего лайка

## 🚀 **API для работы с лайками:**

### **1. Лайкнуть сообщение:**
```bash
POST /api/v1/messages/{message_id}/like
Authorization: Bearer {token}

# Ответ:
{
  "message_id": "uuid",
  "likes_count": 5,
  "liked_by": ["user1", "user2", "user3", "user4", "user5"],
  "liked_at": "2025-10-03T10:00:00Z",
  "is_liked": true
}
```

### **2. Убрать лайк:**
```bash
DELETE /api/v1/messages/{message_id}/like
Authorization: Bearer {token}

# Ответ:
{
  "message_id": "uuid",
  "likes_count": 4,
  "liked_by": ["user1", "user2", "user3", "user4"],
  "liked_at": "2025-10-03T10:00:00Z",
  "is_liked": false
}
```

### **3. Получить информацию о лайках:**
```bash
GET /api/v1/messages/{message_id}/likes
Authorization: Bearer {token}

# Ответ:
{
  "likes_count": 5,
  "liked_by": ["user1", "user2", "user3", "user4", "user5"],
  "liked_at": "2025-10-03T10:00:00Z",
  "is_liked_by_user": true
}
```

## 📱 **WebSocket события:**

### **1. Сообщение лайкнуто:**
```json
{
  "type": "message_liked",
  "chat_id": "chat_uuid",
  "payload": {
    "message_id": "message_uuid",
    "likes_count": 5,
    "liked_by": ["user1", "user2", "user3", "user4", "user5"],
    "liked_at": "2025-10-03T10:00:00Z"
  }
}
```

### **2. Лайк убран:**
```json
{
  "type": "message_unliked",
  "chat_id": "chat_uuid",
  "payload": {
    "message_id": "message_uuid",
    "likes_count": 4,
    "liked_by": ["user1", "user2", "user3", "user4"],
    "liked_at": "2025-10-03T10:00:00Z"
  }
}
```

## 💻 **Использование в коде:**

### **JavaScript - Лайкнуть сообщение:**
```javascript
async function likeMessage(messageId) {
    try {
        const response = await fetch(`/api/v1/messages/${messageId}/like`, {
            method: 'POST',
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });
        
        if (response.ok) {
            const result = await response.json();
            console.log('Message liked:', result);
            return result;
        } else {
            const error = await response.json();
            console.error('Failed to like message:', error);
            return null;
        }
    } catch (error) {
        console.error('Error liking message:', error);
        return null;
    }
}
```

### **JavaScript - Убрать лайк:**
```javascript
async function unlikeMessage(messageId) {
    try {
        const response = await fetch(`/api/v1/messages/${messageId}/like`, {
            method: 'DELETE',
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });
        
        if (response.ok) {
            const result = await response.json();
            console.log('Message unliked:', result);
            return result;
        } else {
            const error = await response.json();
            console.error('Failed to unlike message:', error);
            return null;
        }
    } catch (error) {
        console.error('Error unliking message:', error);
        return null;
    }
}
```

### **JavaScript - Переключить лайк:**
```javascript
async function toggleLike(messageId, isCurrentlyLiked) {
    if (isCurrentlyLiked) {
        return await unlikeMessage(messageId);
    } else {
        return await likeMessage(messageId);
    }
}
```

### **JavaScript - Отображение лайков:**
```javascript
function displayMessageLikes(message) {
    const likesCount = message.likes_count || 0;
    const isLiked = message.is_liked_by_user || false;
    
    return `
        <div class="message-likes">
            <button class="like-button ${isLiked ? 'liked' : ''}" 
                    onclick="toggleLike('${message.id}', ${isLiked})">
                ❤️ ${likesCount}
            </button>
            ${likesCount > 0 ? `<span class="likes-count">${likesCount} лайков</span>` : ''}
        </div>
    `;
}
```

### **JavaScript - WebSocket обработка:**
```javascript
ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    
    switch(data.type) {
        case 'message_liked':
            updateMessageLikes(data.payload);
            break;
        case 'message_unliked':
            updateMessageLikes(data.payload);
            break;
    }
};

function updateMessageLikes(payload) {
    const messageElement = document.querySelector(`[data-message-id="${payload.message_id}"]`);
    if (messageElement) {
        const likesElement = messageElement.querySelector('.message-likes');
        likesElement.innerHTML = displayMessageLikes(payload);
    }
}
```

## 🎯 **Особенности реализации:**

### **1. Безопасность:**
- ✅ **Только участники чата** могут лайкать сообщения
- ✅ **Один лайк на пользователя** - нельзя лайкнуть дважды
- ✅ **Проверка прав доступа** к чату

### **2. Логика лайков:**
- **Добавление:** Если пользователь еще не лайкал
- **Удаление:** Если пользователь уже лайкал
- **Счетчик:** Автоматически обновляется
- **Время:** Сохраняется время последнего лайка

### **3. WebSocket уведомления:**
- **message_liked** - когда кто-то лайкнул
- **message_unliked** - когда кто-то убрал лайк
- **Реальное время** - все участники чата видят изменения

## 🚀 **Возможные улучшения:**

### **1. Дополнительные реакции:**
- **Эмодзи реакции** - 😀 😂 ❤️ 👍 👎
- **Типы лайков** - разные виды реакций
- **Анимации** - эффекты при лайке

### **2. Аналитика:**
- **Популярные сообщения** - топ лайкнутых
- **Статистика пользователей** - кто больше лайкает
- **Временные метрики** - активность лайков

### **3. Уведомления:**
- **Push-уведомления** - когда лайкают твои сообщения
- **Email уведомления** - для важных сообщений
- **In-app уведомления** - в интерфейсе

## 📊 **Примеры использования:**

### **1. Простой лайк:**
```javascript
// Лайкнуть сообщение
likeMessage('message-uuid');
```

### **2. Переключение лайка:**
```javascript
// Переключить лайк (если лайкнуто - убрать, если нет - лайкнуть)
toggleLike('message-uuid', currentUserLiked);
```

### **3. Отображение в UI:**
```html
<div class="message" data-message-id="message-uuid">
    <div class="message-content">Текст сообщения</div>
    <div class="message-likes">
        <button class="like-button">❤️ 0</button>
    </div>
</div>
```

## ✅ **Готово к использованию:**

**Система лайков полностью функциональна!**
- 🎯 **Лайки и анлайки работают**
- 🔄 **WebSocket уведомления в реальном времени**
- 🛡️ **Безопасность и проверка прав**
- 📱 **Готово для мобильных приложений**

**Теперь можно лайкать сообщения!** ❤️✨
