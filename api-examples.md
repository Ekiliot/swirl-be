# Примеры API запросов

## Аутентификация

### Регистрация
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### Авторизация
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

## Профиль

### Получить профиль
```bash
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Обновить профиль
```bash
curl -X PUT http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newusername"
  }'
```

## Чаты

### Получить список чатов
```bash
curl -X GET http://localhost:8080/api/v1/chats \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Создать чат
```bash
curl -X POST http://localhost:8080/api/v1/chats \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Chat",
    "description": "Chat description",
    "type": "private"
  }'
```

### Получить чат
```bash
curl -X GET http://localhost:8080/api/v1/chats/CHAT_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Удалить чат
```bash
curl -X DELETE http://localhost:8080/api/v1/chats/CHAT_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Сообщения

### Получить сообщения чата
```bash
curl -X GET http://localhost:8080/api/v1/chats/CHAT_ID/messages \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Отправить текстовое сообщение
```bash
curl -X POST http://localhost:8080/api/v1/chats/CHAT_ID/messages \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "text",
    "content": "Hello, world!"
  }'
```

### Отправить сообщение с ответом
```bash
curl -X POST http://localhost:8080/api/v1/chats/CHAT_ID/messages \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "text",
    "content": "This is a reply",
    "reply_to_id": "MESSAGE_ID"
  }'
```

### Удалить сообщение
```bash
curl -X DELETE http://localhost:8080/api/v1/messages/MESSAGE_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Загрузка файлов

### Загрузить файл
```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/your/file.jpg"
```

### Отправить сообщение с медиа
```bash
curl -X POST http://localhost:8080/api/v1/chats/CHAT_ID/messages \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "image",
    "content": "Check this out!",
    "media_url": "/uploads/USER_ID/FILE_NAME"
  }'
```

## Чатрулетка

### Найти случайного пользователя
```bash
curl -X GET http://localhost:8080/api/v1/chatroulette/find \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Сохранить чат из чатрулетки
```bash
curl -X POST http://localhost:8080/api/v1/chatroulette/CHAT_ID/save \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Пропустить пользователя
```bash
curl -X POST http://localhost:8080/api/v1/chatroulette/CHAT_ID/skip \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## WebSocket

### Подключение к WebSocket
```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws?chat_id=CHAT_ID');

ws.onopen = function() {
    console.log('Connected to chat');
};

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    console.log('New message:', data);
};

ws.onclose = function() {
    console.log('Disconnected from chat');
};
```

## Типы сообщений

- `text` - текстовое сообщение
- `sticker` - стикер
- `gif` - GIF анимация
- `voice` - голосовое сообщение
- `video` - видео сообщение
- `image` - изображение
- `reply` - ответ на сообщение

## Типы чатов

- `private` - приватный чат
- `group` - групповой чат
- `saved` - сохраненный чат из чатрулетки
