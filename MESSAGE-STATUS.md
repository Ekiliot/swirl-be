# 📨 Статусы сообщений в Chatroulette Backend

## ✅ **Что добавлено:**

### **📊 Статусы сообщений:**
- **sent** - Отправлено (по умолчанию)
- **delivered** - Доставлено
- **read** - Прочитано
- **edited** - Изменено
- **deleted** - Удалено

### **🔍 Дополнительные поля:**
- **IsEdited** - Флаг редактирования
- **EditedAt** - Время редактирования
- **ReadAt** - Время первого прочтения
- **ReadBy** - Список пользователей, прочитавших сообщение

## 🚀 **API для работы со статусами:**

### **1. Отметить сообщение как прочитанное:**
```bash
PUT /api/v1/messages/{message_id}/read
Authorization: Bearer {token}

# Ответ:
{
  "message_id": "uuid",
  "status": "read",
  "status_text": "Прочитано",
  "read_by": ["user1", "user2"],
  "read_at": "2025-10-03T10:00:00Z"
}
```

### **2. Редактировать сообщение:**
```bash
PUT /api/v1/messages/{message_id}/edit
Authorization: Bearer {token}
Content-Type: application/json

{
  "content": "Новый текст сообщения"
}

# Ответ:
{
  "id": "uuid",
  "content": "Новый текст сообщения",
  "status": "edited",
  "is_edited": true,
  "edited_at": "2025-10-03T10:00:00Z"
}
```

### **3. Получить статус сообщения:**
```bash
GET /api/v1/messages/{message_id}/status
Authorization: Bearer {token}

# Ответ:
{
  "message_id": "uuid",
  "status": "read",
  "status_text": "Прочитано",
  "is_edited": false,
  "edited_at": null,
  "read_at": "2025-10-03T10:00:00Z",
  "read_by": ["user1", "user2"],
  "created_at": "2025-10-03T09:00:00Z",
  "updated_at": "2025-10-03T10:00:00Z"
}
```

## 📱 **WebSocket события:**

### **1. Обновление статуса сообщения:**
```json
{
  "type": "message_status_update",
  "chat_id": "chat_uuid",
  "payload": {
    "message_id": "message_uuid",
    "status": "read",
    "read_by": ["user1", "user2"],
    "read_at": "2025-10-03T10:00:00Z"
  }
}
```

### **2. Редактирование сообщения:**
```json
{
  "type": "message_edited",
  "chat_id": "chat_uuid",
  "payload": {
    "id": "message_uuid",
    "content": "Новый текст",
    "status": "edited",
    "is_edited": true,
    "edited_at": "2025-10-03T10:00:00Z"
  }
}
```

## 🔄 **Жизненный цикл сообщения:**

### **1. Отправка:**
```
Пользователь отправляет → Status: "sent"
```

### **2. Доставка:**
```
Сервер получает → Status: "delivered"
```

### **3. Прочтение:**
```
Пользователь читает → Status: "read"
```

### **4. Редактирование:**
```
Автор редактирует → Status: "edited"
```

### **5. Удаление:**
```
Автор удаляет → Status: "deleted"
```

## 💻 **Использование в коде:**

### **JavaScript - Отметить как прочитанное:**
```javascript
async function markAsRead(messageId) {
    const response = await fetch(`/api/v1/messages/${messageId}/read`, {
        method: 'PUT',
        headers: {
            'Authorization': 'Bearer ' + token
        }
    });
    
    const result = await response.json();
    console.log('Message marked as read:', result);
}
```

### **JavaScript - Редактировать сообщение:**
```javascript
async function editMessage(messageId, newContent) {
    const response = await fetch(`/api/v1/messages/${messageId}/edit`, {
        method: 'PUT',
        headers: {
            'Authorization': 'Bearer ' + token,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            content: newContent
        })
    });
    
    const result = await response.json();
    console.log('Message edited:', result);
}
```

### **JavaScript - Отображение статуса:**
```javascript
function displayMessageStatus(message) {
    const statusIcons = {
        'sent': '📤',
        'delivered': '📨',
        'read': '✅',
        'edited': '✏️',
        'deleted': '🗑️'
    };
    
    const statusText = {
        'sent': 'Отправлено',
        'delivered': 'Доставлено',
        'read': 'Прочитано',
        'edited': 'Изменено',
        'deleted': 'Удалено'
    };
    
    return `
        <div class="message-status">
            <span class="status-icon">${statusIcons[message.status]}</span>
            <span class="status-text">${statusText[message.status]}</span>
            ${message.is_edited ? '<span class="edited-badge">(изменено)</span>' : ''}
        </div>
    `;
}
```

### **JavaScript - WebSocket обработка:**
```javascript
ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    
    switch(data.type) {
        case 'message_status_update':
            updateMessageStatus(data.payload);
            break;
        case 'message_edited':
            updateMessageContent(data.payload);
            break;
    }
};

function updateMessageStatus(payload) {
    const messageElement = document.querySelector(`[data-message-id="${payload.message_id}"]`);
    if (messageElement) {
        const statusElement = messageElement.querySelector('.message-status');
        statusElement.innerHTML = displayMessageStatus(payload);
    }
}
```

## 🎯 **Особенности реализации:**

### **1. Автоматические статусы:**
- **Отправка:** Автоматически устанавливается "sent"
- **Доставка:** Автоматически меняется на "delivered"
- **Прочтение:** Требует явного вызова API

### **2. Безопасность:**
- ✅ **Редактирование:** Только автор может редактировать
- ✅ **Прочтение:** Только участники чата могут отмечать как прочитанное
- ✅ **Доступ:** Проверка прав доступа к чату

### **3. Отслеживание:**
- **ReadBy:** Массив ID пользователей, прочитавших сообщение
- **ReadAt:** Время первого прочтения
- **EditedAt:** Время последнего редактирования

## 🚀 **Возможные улучшения:**

### **1. Групповые чаты:**
- **Все прочитали:** Автоматическое обновление статуса
- **Кто прочитал:** Отображение списка прочитавших
- **Уведомления:** Push-уведомления о прочтении

### **2. Дополнительные статусы:**
- **typing** - Пользователь печатает
- **sending** - Отправляется
- **failed** - Ошибка отправки

### **3. Аналитика:**
- **Время доставки:** Метрики скорости
- **Процент прочтения:** Статистика по чатам
- **Активность:** Время ответа пользователей

## ✅ **Готово к использованию:**

**Система статусов сообщений полностью функциональна!**
- 🎯 **Все основные статусы поддерживаются**
- 🔄 **Автоматические обновления через WebSocket**
- 🛡️ **Безопасность и проверка прав**
- 📱 **Готово для мобильных приложений**

**Теперь можно отслеживать статус каждого сообщения!** 📨✨
