# 🔌 WebSocket события в Chatroulette Backend

## 📡 **Подключение к WebSocket**

### **URL подключения:**
```
ws://localhost:8080/api/v1/ws?chat_id={chat_id}&token={jwt_token}
```

### **Параметры:**
- `chat_id` - ID чата для подключения
- `token` - JWT токен для аутентификации

### **JavaScript пример:**
```javascript
const ws = new WebSocket(`ws://localhost:8080/api/v1/ws?chat_id=${chatId}&token=${token}`);

ws.onopen = function(event) {
    console.log('WebSocket connected');
};

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    handleWebSocketMessage(data);
};

ws.onerror = function(error) {
    console.error('WebSocket error:', error);
};

ws.onclose = function(event) {
    console.log('WebSocket disconnected');
};
```

---

## 📨 **События сообщений**

### **1. Новое сообщение**
```json
{
  "type": "new_message",
  "chat_id": "uuid",
  "payload": {
    "id": "message_uuid",
    "chat_id": "chat_uuid",
    "user_id": "user_uuid",
    "type": "text",
    "content": "Hello!",
    "media_url": "",
    "status": "delivered",
    "is_edited": false,
    "edited_at": null,
    "read_at": null,
    "read_by": [],
    "likes_count": 0,
    "liked_by": [],
    "liked_at": null,
    "created_at": "2025-10-03T10:00:00Z",
    "user": {
      "id": "user_uuid",
      "username": "vasya"
    }
  }
}
```

**Обработка:**
```javascript
function handleNewMessage(payload) {
    const messageElement = createMessageElement(payload);
    document.getElementById('messages').appendChild(messageElement);
    scrollToBottom();
}
```

### **2. Сообщение отредактировано**
```json
{
  "type": "message_edited",
  "chat_id": "uuid",
  "payload": {
    "id": "message_uuid",
    "content": "Updated message text",
    "status": "edited",
    "is_edited": true,
    "edited_at": "2025-10-03T10:00:00Z"
  }
}
```

**Обработка:**
```javascript
function handleMessageEdited(payload) {
    const messageElement = document.querySelector(`[data-message-id="${payload.id}"]`);
    if (messageElement) {
        const contentElement = messageElement.querySelector('.message-content');
        contentElement.textContent = payload.content;
        
        // Добавляем индикатор редактирования
        if (payload.is_edited) {
            const editedBadge = document.createElement('span');
            editedBadge.className = 'edited-badge';
            editedBadge.textContent = '(изменено)';
            messageElement.appendChild(editedBadge);
        }
    }
}
```

### **3. Сообщение удалено**
```json
{
  "type": "message_deleted",
  "chat_id": "uuid",
  "payload": {
    "message_id": "message_uuid"
  }
}
```

**Обработка:**
```javascript
function handleMessageDeleted(payload) {
    const messageElement = document.querySelector(`[data-message-id="${payload.message_id}"]`);
    if (messageElement) {
        messageElement.remove();
    }
}
```

---

## 📊 **События статусов**

### **1. Обновление статуса сообщения**
```json
{
  "type": "message_status_update",
  "chat_id": "uuid",
  "payload": {
    "message_id": "message_uuid",
    "status": "read",
    "read_by": ["user1", "user2"],
    "read_at": "2025-10-03T10:00:00Z"
  }
}
```

**Обработка:**
```javascript
function handleStatusUpdate(payload) {
    const messageElement = document.querySelector(`[data-message-id="${payload.message_id}"]`);
    if (messageElement) {
        const statusElement = messageElement.querySelector('.message-status');
        statusElement.textContent = getStatusText(payload.status);
        
        // Обновляем индикаторы прочтения
        updateReadIndicators(messageElement, payload.read_by);
    }
}

function getStatusText(status) {
    const statusTexts = {
        'sent': 'Отправлено',
        'delivered': 'Доставлено',
        'read': 'Прочитано',
        'edited': 'Изменено',
        'deleted': 'Удалено'
    };
    return statusTexts[status] || 'Неизвестно';
}
```

---

## ❤️ **События лайков**

### **1. Сообщение лайкнуто**
```json
{
  "type": "message_liked",
  "chat_id": "uuid",
  "payload": {
    "message_id": "message_uuid",
    "likes_count": 5,
    "liked_by": ["user1", "user2", "user3", "user4", "user5"],
    "liked_at": "2025-10-03T10:00:00Z"
  }
}
```

**Обработка:**
```javascript
function handleMessageLiked(payload) {
    const messageElement = document.querySelector(`[data-message-id="${payload.message_id}"]`);
    if (messageElement) {
        const likesElement = messageElement.querySelector('.message-likes');
        likesElement.innerHTML = `
            <button class="like-button liked" onclick="unlikeMessage('${payload.message_id}')">
                ❤️ ${payload.likes_count}
            </button>
        `;
    }
}
```

### **2. Лайк убран**
```json
{
  "type": "message_unliked",
  "chat_id": "uuid",
  "payload": {
    "message_id": "message_uuid",
    "likes_count": 4,
    "liked_by": ["user1", "user2", "user3", "user4"],
    "liked_at": "2025-10-03T10:00:00Z"
  }
}
```

**Обработка:**
```javascript
function handleMessageUnliked(payload) {
    const messageElement = document.querySelector(`[data-message-id="${payload.message_id}"]`);
    if (messageElement) {
        const likesElement = messageElement.querySelector('.message-likes');
        likesElement.innerHTML = `
            <button class="like-button" onclick="likeMessage('${payload.message_id}')">
                ❤️ ${payload.likes_count}
            </button>
        `;
    }
}
```

---

## 🎲 **События чатрулетки**

### **1. Найден новый пользователь**
```json
{
  "type": "chatroulette_match",
  "chat_id": "uuid",
  "payload": {
    "user": {
      "id": "user_uuid",
      "username": "masha",
      "is_online": true,
      "last_seen": "2025-10-03T10:00:00Z"
    },
    "chat_id": "chat_uuid"
  }
}
```

**Обработка:**
```javascript
function handleChatrouletteMatch(payload) {
    // Обновляем UI для показа найденного пользователя
    document.getElementById('search-status').textContent = 'Найден собеседник!';
    document.getElementById('partner-info').innerHTML = `
        <h3>Собеседник: ${payload.user.username}</h3>
        <p>Статус: ${payload.user.is_online ? 'В сети' : 'Не в сети'}</p>
    `;
    
    // Подключаемся к чату
    connectToChat(payload.chat_id);
}
```

### **2. Пользователь покинул чатрулетку**
```json
{
  "type": "chatroulette_left",
  "chat_id": "uuid",
  "payload": {
    "user_id": "user_uuid",
    "username": "masha"
  }
}
```

**Обработка:**
```javascript
function handleChatrouletteLeft(payload) {
    document.getElementById('search-status').textContent = 'Собеседник покинул чат';
    // Возвращаемся к поиску
    startSearch();
}
```

---

## 🔄 **Универсальный обработчик**

### **JavaScript - Полный обработчик WebSocket**
```javascript
class WebSocketHandler {
    constructor(chatId, token) {
        this.chatId = chatId;
        this.token = token;
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
    }

    connect() {
        const wsUrl = `ws://localhost:8080/api/v1/ws?chat_id=${this.chatId}&token=${this.token}`;
        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = (event) => {
            console.log('WebSocket connected');
            this.reconnectAttempts = 0;
        };

        this.ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                this.handleMessage(data);
            } catch (error) {
                console.error('Error parsing WebSocket message:', error);
            }
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        this.ws.onclose = (event) => {
            console.log('WebSocket disconnected');
            this.handleReconnect();
        };
    }

    handleMessage(data) {
        switch (data.type) {
            case 'new_message':
                this.handleNewMessage(data.payload);
                break;
            case 'message_edited':
                this.handleMessageEdited(data.payload);
                break;
            case 'message_deleted':
                this.handleMessageDeleted(data.payload);
                break;
            case 'message_status_update':
                this.handleStatusUpdate(data.payload);
                break;
            case 'message_liked':
                this.handleMessageLiked(data.payload);
                break;
            case 'message_unliked':
                this.handleMessageUnliked(data.payload);
                break;
            case 'chatroulette_match':
                this.handleChatrouletteMatch(data.payload);
                break;
            case 'chatroulette_left':
                this.handleChatrouletteLeft(data.payload);
                break;
            default:
                console.log('Unknown WebSocket event:', data.type);
        }
    }

    handleNewMessage(payload) {
        // Добавляем новое сообщение в чат
        const messageElement = this.createMessageElement(payload);
        document.getElementById('messages').appendChild(messageElement);
        this.scrollToBottom();
    }

    handleMessageEdited(payload) {
        // Обновляем отредактированное сообщение
        const messageElement = document.querySelector(`[data-message-id="${payload.id}"]`);
        if (messageElement) {
            const contentElement = messageElement.querySelector('.message-content');
            contentElement.textContent = payload.content;
            
            if (payload.is_edited) {
                const editedBadge = document.createElement('span');
                editedBadge.className = 'edited-badge';
                editedBadge.textContent = '(изменено)';
                messageElement.appendChild(editedBadge);
            }
        }
    }

    handleMessageDeleted(payload) {
        // Удаляем сообщение из чата
        const messageElement = document.querySelector(`[data-message-id="${payload.message_id}"]`);
        if (messageElement) {
            messageElement.remove();
        }
    }

    handleStatusUpdate(payload) {
        // Обновляем статус сообщения
        const messageElement = document.querySelector(`[data-message-id="${payload.message_id}"]`);
        if (messageElement) {
            const statusElement = messageElement.querySelector('.message-status');
            statusElement.textContent = this.getStatusText(payload.status);
        }
    }

    handleMessageLiked(payload) {
        // Обновляем лайки сообщения
        const messageElement = document.querySelector(`[data-message-id="${payload.message_id}"]`);
        if (messageElement) {
            const likesElement = messageElement.querySelector('.message-likes');
            likesElement.innerHTML = `
                <button class="like-button liked" onclick="unlikeMessage('${payload.message_id}')">
                    ❤️ ${payload.likes_count}
                </button>
            `;
        }
    }

    handleMessageUnliked(payload) {
        // Обновляем лайки сообщения
        const messageElement = document.querySelector(`[data-message-id="${payload.message_id}"]`);
        if (messageElement) {
            const likesElement = messageElement.querySelector('.message-likes');
            likesElement.innerHTML = `
                <button class="like-button" onclick="likeMessage('${payload.message_id}')">
                    ❤️ ${payload.likes_count}
                </button>
            `;
        }
    }

    handleChatrouletteMatch(payload) {
        // Обрабатываем найденного собеседника
        document.getElementById('search-status').textContent = 'Найден собеседник!';
        document.getElementById('partner-info').innerHTML = `
            <h3>Собеседник: ${payload.user.username}</h3>
            <p>Статус: ${payload.user.is_online ? 'В сети' : 'Не в сети'}</p>
        `;
    }

    handleChatrouletteLeft(payload) {
        // Обрабатываем уход собеседника
        document.getElementById('search-status').textContent = 'Собеседник покинул чат';
    }

    handleReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            console.log(`Attempting to reconnect... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
            setTimeout(() => this.connect(), 3000);
        } else {
            console.error('Max reconnection attempts reached');
        }
    }

    createMessageElement(message) {
        const messageDiv = document.createElement('div');
        messageDiv.className = 'message';
        messageDiv.setAttribute('data-message-id', message.id);
        
        messageDiv.innerHTML = `
            <div class="message-header">
                <span class="username">${message.user.username}</span>
                <span class="timestamp">${new Date(message.created_at).toLocaleTimeString()}</span>
            </div>
            <div class="message-content">${message.content}</div>
            <div class="message-status">${this.getStatusText(message.status)}</div>
            <div class="message-likes">
                <button class="like-button" onclick="likeMessage('${message.id}')">
                    ❤️ ${message.likes_count || 0}
                </button>
            </div>
        `;
        
        return messageDiv;
    }

    getStatusText(status) {
        const statusTexts = {
            'sent': 'Отправлено',
            'delivered': 'Доставлено',
            'read': 'Прочитано',
            'edited': 'Изменено',
            'deleted': 'Удалено'
        };
        return statusTexts[status] || 'Неизвестно';
    }

    scrollToBottom() {
        const messagesContainer = document.getElementById('messages');
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }

    disconnect() {
        if (this.ws) {
            this.ws.close();
        }
    }
}

// Использование
const wsHandler = new WebSocketHandler(chatId, token);
wsHandler.connect();
```

---

## 🎯 **Лучшие практики**

### **1. Обработка ошибок**
```javascript
ws.onerror = function(error) {
    console.error('WebSocket error:', error);
    // Показать пользователю уведомление об ошибке
    showNotification('Ошибка соединения', 'error');
};
```

### **2. Переподключение**
```javascript
ws.onclose = function(event) {
    if (event.code !== 1000) { // Не нормальное закрытие
        setTimeout(() => {
            connectWebSocket();
        }, 3000);
    }
};
```

### **3. Обработка больших сообщений**
```javascript
ws.onmessage = function(event) {
    try {
        const data = JSON.parse(event.data);
        // Проверяем размер payload
        if (JSON.stringify(data.payload).length > 1000000) {
            console.warn('Large message received');
        }
        handleMessage(data);
    } catch (error) {
        console.error('Error parsing message:', error);
    }
};
```

---

## ✅ **Заключение**

**WebSocket события обеспечивают:**
- 🔄 **Реальное время** - мгновенные обновления
- 📨 **Синхронизацию сообщений** - все участники видят изменения
- ❤️ **Интерактивность** - лайки, статусы, редактирование
- 🎲 **Чатрулетку** - поиск и подключение к собеседникам

**Система WebSocket полностью готова к использованию!** 🔌✨
