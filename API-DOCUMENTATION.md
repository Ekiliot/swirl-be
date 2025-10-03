# 📚 Полная документация API Swirl Backend

## 🚀 **Обзор API**

**Base URL:** `http://localhost:8080/api/v1`  
**Версия:** 1.0  
**Аутентификация:** JWT Bearer Token  

### **📋 Содержание:**
1. [Аутентификация](#аутентификация)
2. [Пользователи](#пользователи)
3. [Чаты](#чаты)
4. [Сообщения](#сообщения)
5. [Медиа файлы](#медиа-файлы)
6. [Swirl (случайные встречи)](#swirl-случайные-встречи)
7. [Поисковая очередь](#поисковая-очередь)
8. [WebSocket](#websocket)
9. [Коды ошибок](#коды-ошибок)
10. [Примеры использования](#примеры-использования)

---

## 🔐 **Аутентификация**

### **Регистрация пользователя**
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "vasya",
  "password": "password123"
}
```

**Ответ (201 Created):**
```json
{
  "user": {
    "id": "uuid",
    "username": "vasya",
    "created_at": "2025-10-03T10:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Ошибки:**
- `400` - Неверные данные
- `409` - Пользователь уже существует

### **Вход в систему**
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "vasya",
  "password": "password123"
}
```

**Ответ (200 OK):**
```json
{
  "user": {
    "id": "uuid",
    "username": "vasya",
    "is_online": true,
    "last_seen": "2025-10-03T10:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Ошибки:**
- `401` - Неверные учетные данные
- `400` - Неверные данные

### **Получение профиля**
```http
GET /api/v1/auth/profile
Authorization: Bearer {token}
```

**Ответ (200 OK):**
```json
{
  "id": "uuid",
  "username": "vasya",
  "birthday": "1990-01-01T00:00:00Z",
  "profile_photo": "/uploads/user_id/photo.jpg",
  "is_online": true,
  "last_seen": "2025-10-03T10:00:00Z",
  "show_username": true,
  "show_birthday": true,
  "show_online_status": true,
  "created_at": "2025-10-03T10:00:00Z"
}
```

### **Обновление профиля**
```http
PUT /api/v1/auth/profile
Authorization: Bearer {token}
Content-Type: application/json

{
  "birthday": "1990-01-01T00:00:00Z",
  "profile_photo": "/uploads/user_id/photo.jpg",
  "show_username": true,
  "show_birthday": false,
  "show_online_status": true
}
```

### **Обновление онлайн статуса**
```http
POST /api/v1/profile/online
Authorization: Bearer {token}
Content-Type: application/json

{
  "is_online": true
}
```

### **Получение публичного профиля**
```http
GET /api/v1/users/{user_id}/profile
Authorization: Bearer {token}
```

---

## 👥 **Пользователи**

### **Получение списка пользователей**
```http
GET /api/v1/users
Authorization: Bearer {token}
```

**Ответ (200 OK):**
```json
[
  {
    "id": "uuid",
    "username": "vasya",
    "is_online": true,
    "last_seen": "2025-10-03T10:00:00Z",
    "status_text": "В сети"
  }
]
```

---

## 💬 **Чаты**

### **Получение списка чатов**
```http
GET /api/v1/chats
Authorization: Bearer {token}
```

**Ответ (200 OK):**
```json
[
  {
    "id": "uuid",
    "name": "Chat with vasya",
    "type": "saved",
    "description": "Personal chat",
    "created_by": "uuid",
    "created_at": "2025-10-03T10:00:00Z",
    "participants": [
      {
        "user_id": "uuid",
        "username": "vasya"
      }
    ]
  }
]
```

### **Создание чата**
```http
POST /api/v1/chats
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "New Chat",
  "type": "private",
  "description": "Chat description"
}
```

### **Получение чата**
```http
GET /api/v1/chats/{chat_id}
Authorization: Bearer {token}
```

### **Удаление чата**
```http
DELETE /api/v1/chats/{chat_id}
Authorization: Bearer {token}
```

---

## 📨 **Сообщения**

### **Получение сообщений чата**
```http
GET /api/v1/chats/{chat_id}/messages?limit=50&offset=0
Authorization: Bearer {token}
```

**Ответ (200 OK):**
```json
[
  {
    "id": "uuid",
    "chat_id": "uuid",
    "user_id": "uuid",
    "type": "text",
    "content": "Hello!",
    "media_url": "",
    "status": "read",
    "is_edited": false,
    "edited_at": null,
    "read_at": "2025-10-03T10:00:00Z",
    "read_by": ["uuid1", "uuid2"],
    "likes_count": 5,
    "liked_by": ["uuid1", "uuid2"],
    "liked_at": "2025-10-03T10:00:00Z",
    "created_at": "2025-10-03T10:00:00Z",
    "user": {
      "id": "uuid",
      "username": "vasya"
    }
  }
]
```

### **Отправка сообщения**
```http
POST /api/v1/chats/{chat_id}/messages
Authorization: Bearer {token}
Content-Type: application/json

{
  "type": "text",
  "content": "Hello!",
  "media_url": "",
  "reply_to_id": "uuid"
}
```

**Типы сообщений:**
- `text` - текстовое сообщение
- `image` - изображение
- `video` - видео
- `voice` - голосовое сообщение
- `sticker` - стикер
- `gif` - GIF анимация

### **Редактирование сообщения**
```http
PUT /api/v1/messages/{message_id}/edit
Authorization: Bearer {token}
Content-Type: application/json

{
  "content": "Updated message text"
}
```

### **Удаление сообщения**
```http
DELETE /api/v1/messages/{message_id}
Authorization: Bearer {token}
```

### **Отметить как прочитанное**
```http
PUT /api/v1/messages/{message_id}/read
Authorization: Bearer {token}
```

### **Получить статус сообщения**
```http
GET /api/v1/messages/{message_id}/status
Authorization: Bearer {token}
```

---

## ❤️ **Лайки сообщений**

### **Лайкнуть сообщение**
```http
POST /api/v1/messages/{message_id}/like
Authorization: Bearer {token}
```

**Ответ (200 OK):**
```json
{
  "message_id": "uuid",
  "likes_count": 5,
  "liked_by": ["uuid1", "uuid2", "uuid3"],
  "liked_at": "2025-10-03T10:00:00Z",
  "is_liked": true
}
```

### **Убрать лайк**
```http
DELETE /api/v1/messages/{message_id}/like
Authorization: Bearer {token}
```

### **Получить информацию о лайках**
```http
GET /api/v1/messages/{message_id}/likes
Authorization: Bearer {token}
```

---

## 📁 **Медиа файлы**

### **Загрузка файла**
```http
POST /api/v1/upload
Authorization: Bearer {token}
Content-Type: multipart/form-data

file: [файл]
```

**Ответ (200 OK):**
```json
{
  "file_id": "user_id_uuid.jpg",
  "file_url": "/uploads/user_id/filename.jpg",
  "media_type": "image",
  "size": 1024000,
  "created_at": "2025-10-03T10:00:00Z",
  "info": {
    "content_type": "image/jpeg",
    "is_image": true
  }
}
```

**Поддерживаемые форматы:**
- **Изображения:** JPEG, PNG, GIF, WebP (до 10MB)
- **Видео:** MP4, WebM, QuickTime (до 100MB)
- **Аудио:** MP3, WAV, OGG, M4A (до 25MB)

### **Получение файла**
```http
GET /uploads/{user_id}/{file_name}
```

### **Удаление файла**
```http
DELETE /api/v1/uploads/{file_name}
Authorization: Bearer {token}
```

---

## 🌪️ **Swirl (случайные встречи)**

### **Найти случайного пользователя**
```http
GET /api/v1/swirl/find
Authorization: Bearer {token}
```

**Ответ (200 OK):**
```json
{
  "user": {
    "id": "uuid",
    "username": "masha",
    "is_online": true,
    "last_seen": "2025-10-03T10:00:00Z"
  },
  "chat_id": "uuid"
}
```

**Ответ (404 Not Found):**
```json
{
  "error": "No users available"
}
```

### **Сохранить чат**
```http
POST /api/v1/swirl/{chat_id}/save
Authorization: Bearer {token}
```

**Ответ (200 OK):**
```json
{
  "message": "Chat saved successfully",
  "chat": {
    "id": "uuid",
    "name": "Saved Chat - Chatroulette Chat",
    "type": "saved"
  }
}
```

**Ответ (409 Conflict):**
```json
{
  "error": "Chat already saved"
}
```

### **Пропустить пользователя**
```http
POST /api/v1/swirl/{chat_id}/skip
Authorization: Bearer {token}
```

---

## 🔍 **Поисковая очередь**

### **Обновить активность поиска**
```http
POST /api/v1/swirl/activity
Authorization: Bearer {token}
```

### **Получить статус очереди**
```http
GET /api/v1/swirl/status
Authorization: Bearer {token}
```

**Ответ (200 OK):**
```json
{
  "queue_size": 5,
  "users_in_queue": [
    {
      "user_id": "uuid",
      "username": "vasya",
      "joined_at": "2025-10-03T10:00:00Z"
    }
  ]
}
```

### **Очистить очередь**
```http
DELETE /api/v1/swirl/clear
Authorization: Bearer {token}
```

---

## 🔌 **WebSocket**

### **Подключение**
```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws?chat_id=uuid&token=jwt_token');
```

### **События**

#### **Новое сообщение**
```json
{
  "type": "new_message",
  "chat_id": "uuid",
  "payload": {
    "id": "uuid",
    "content": "Hello!",
    "user": {
      "id": "uuid",
      "username": "vasya"
    }
  }
}
```

#### **Сообщение отредактировано**
```json
{
  "type": "message_edited",
  "chat_id": "uuid",
  "payload": {
    "id": "uuid",
    "content": "Updated text",
    "is_edited": true,
    "edited_at": "2025-10-03T10:00:00Z"
  }
}
```

#### **Обновление статуса сообщения**
```json
{
  "type": "message_status_update",
  "chat_id": "uuid",
  "payload": {
    "message_id": "uuid",
    "status": "read",
    "read_by": ["uuid1", "uuid2"],
    "read_at": "2025-10-03T10:00:00Z"
  }
}
```

#### **Сообщение лайкнуто**
```json
{
  "type": "message_liked",
  "chat_id": "uuid",
  "payload": {
    "message_id": "uuid",
    "likes_count": 5,
    "liked_by": ["uuid1", "uuid2"],
    "liked_at": "2025-10-03T10:00:00Z"
  }
}
```

#### **Лайк убран**
```json
{
  "type": "message_unliked",
  "chat_id": "uuid",
  "payload": {
    "message_id": "uuid",
    "likes_count": 4,
    "liked_by": ["uuid1", "uuid2"],
    "liked_at": "2025-10-03T10:00:00Z"
  }
}
```

---

## ❌ **Коды ошибок**

### **HTTP статус коды**
- `200` - OK
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `429` - Too Many Requests
- `500` - Internal Server Error

### **Структура ошибки**
```json
{
  "error": "Описание ошибки"
}
```

### **Частые ошибки**
- `"Invalid credentials"` - Неверные учетные данные
- `"User already exists"` - Пользователь уже существует
- `"Access denied"` - Доступ запрещен
- `"Message not found"` - Сообщение не найдено
- `"Chat not found"` - Чат не найден
- `"File too large"` - Файл слишком большой
- `"Too many requests"` - Слишком много запросов

---

## 💻 **Примеры использования**

### **JavaScript - Полный цикл работы**

```javascript
class ChatrouletteAPI {
    constructor(baseURL = 'http://localhost:8080/api/v1') {
        this.baseURL = baseURL;
        this.token = null;
    }

    // Аутентификация
    async register(username, password) {
        const response = await fetch(`${this.baseURL}/auth/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });
        const data = await response.json();
        this.token = data.token;
        return data;
    }

    async login(username, password) {
        const response = await fetch(`${this.baseURL}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });
        const data = await response.json();
        this.token = data.token;
        return data;
    }

    // Чаты
    async getChats() {
        return await this.request('GET', '/chats');
    }

    async getChatMessages(chatId, limit = 50, offset = 0) {
        return await this.request('GET', `/chats/${chatId}/messages?limit=${limit}&offset=${offset}`);
    }

    async sendMessage(chatId, content, type = 'text', mediaUrl = '') {
        return await this.request('POST', `/chats/${chatId}/messages`, {
            type, content, media_url: mediaUrl
        });
    }

    // Лайки
    async likeMessage(messageId) {
        return await this.request('POST', `/messages/${messageId}/like`);
    }

    async unlikeMessage(messageId) {
        return await this.request('DELETE', `/messages/${messageId}/like`);
    }

    // Чатрулетка
    async findRandomUser() {
        return await this.request('GET', '/chatroulette/find');
    }

    async saveChat(chatId) {
        return await this.request('POST', `/chatroulette/${chatId}/save`);
    }

    // Медиа
    async uploadFile(file) {
        const formData = new FormData();
        formData.append('file', file);
        
        const response = await fetch(`${this.baseURL}/upload`, {
            method: 'POST',
            headers: { 'Authorization': `Bearer ${this.token}` },
            body: formData
        });
        return await response.json();
    }

    // WebSocket
    connectWebSocket(chatId) {
        return new WebSocket(`ws://localhost:8080/api/v1/ws?chat_id=${chatId}&token=${this.token}`);
    }

    // Вспомогательный метод
    async request(method, endpoint, data = null) {
        const options = {
            method,
            headers: {
                'Authorization': `Bearer ${this.token}`,
                'Content-Type': 'application/json'
            }
        };

        if (data) {
            options.body = JSON.stringify(data);
        }

        const response = await fetch(`${this.baseURL}${endpoint}`, options);
        return await response.json();
    }
}

// Использование
const api = new ChatrouletteAPI();

// Регистрация и вход
await api.register('vasya', 'password123');
// или
await api.login('vasya', 'password123');

// Поиск случайного пользователя
const randomUser = await api.findRandomUser();
console.log('Found user:', randomUser);

// Отправка сообщения
await api.sendMessage(randomUser.chat_id, 'Hello!');

// Лайк сообщения
await api.likeMessage('message-uuid');

// WebSocket подключение
const ws = api.connectWebSocket('chat-uuid');
ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('WebSocket message:', data);
};
```

### **Python - Пример использования**

```python
import requests
import json

class ChatrouletteAPI:
    def __init__(self, base_url='http://localhost:8080/api/v1'):
        self.base_url = base_url
        self.token = None

    def register(self, username, password):
        response = requests.post(f'{self.base_url}/auth/register', 
                               json={'username': username, 'password': password})
        data = response.json()
        self.token = data['token']
        return data

    def login(self, username, password):
        response = requests.post(f'{self.base_url}/auth/login',
                               json={'username': username, 'password': password})
        data = response.json()
        self.token = data['token']
        return data

    def get_headers(self):
        return {'Authorization': f'Bearer {self.token}'}

    def get_chats(self):
        response = requests.get(f'{self.base_url}/chats', headers=self.get_headers())
        return response.json()

    def send_message(self, chat_id, content, message_type='text'):
        response = requests.post(f'{self.base_url}/chats/{chat_id}/messages',
                               headers=self.get_headers(),
                               json={'type': message_type, 'content': content})
        return response.json()

    def like_message(self, message_id):
        response = requests.post(f'{self.base_url}/messages/{message_id}/like',
                               headers=self.get_headers())
        return response.json()

# Использование
api = ChatrouletteAPI()
api.login('vasya', 'password123')
chats = api.get_chats()
```

---

## 🚀 **Быстрый старт**

### **1. Регистрация и вход**
```bash
# Регистрация
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "vasya", "password": "password123"}'

# Вход
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "vasya", "password": "password123"}'
```

### **2. Поиск пользователя в Swirl**
```bash
curl -X GET http://localhost:8080/api/v1/swirl/find \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### **3. Отправка сообщения**
```bash
curl -X POST http://localhost:8080/api/v1/chats/CHAT_ID/messages \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"type": "text", "content": "Hello!"}'
```

### **4. Лайк сообщения**
```bash
curl -X POST http://localhost:8080/api/v1/messages/MESSAGE_ID/like \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 📝 **Заключение**

**API Swirl Backend предоставляет полный набор функций для:**
- 🔐 **Аутентификации и управления пользователями**
- 💬 **Работы с чатами и сообщениями**
- ❤️ **Системы лайков и реакций**
- 📁 **Загрузки и управления медиа файлами**
- 🌪️ **Swirl (случайных встреч) и поиска пользователей**
- 🔌 **WebSocket для реального времени**

**Все функции готовы к использованию и хорошо документированы!** 🚀✨
