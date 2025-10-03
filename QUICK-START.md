# 🚀 Быстрый старт с Swirl Backend

## ⚡ **За 5 минут до работающего API**

### **1. Запуск сервера**
```bash
# Установка зависимостей
go mod tidy

# Запуск сервера
go run main.go
```

**Сервер запустится на:** `http://localhost:8080`

### **2. Первые API вызовы**

#### **Регистрация пользователя**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "vasya", "password": "password123"}'
```

#### **Вход в систему**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "vasya", "password": "password123"}'
```

#### **Поиск собеседника в чатрулетке**
```bash
curl -X GET http://localhost:8080/api/v1/chatroulette/find \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### **Отправка сообщения**
```bash
curl -X POST http://localhost:8080/api/v1/chats/CHAT_ID/messages \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"type": "text", "content": "Hello!"}'
```

---

## 🎯 **Основные функции**

### **✅ Что уже работает:**

#### **🔐 Аутентификация**
- Регистрация и вход
- JWT токены
- Профили пользователей

#### **💬 Чаты и сообщения**
- Создание и управление чатами
- Отправка текстовых сообщений
- Статусы сообщений (отправлено, доставлено, прочитано)
- Редактирование сообщений

#### **❤️ Лайки**
- Лайкнуть/убрать лайк
- Счетчик лайков
- WebSocket уведомления

#### **📁 Медиа файлы**
- Загрузка фото, видео, аудио
- Разные лимиты размеров
- Безопасная обработка

#### **🎲 Чатрулетка**
- Поиск случайных пользователей
- Сохранение чатов
- Очередь поиска

#### **🔌 WebSocket**
- Реальное время
- События сообщений
- Синхронизация

---

## 📱 **Тестирование с HTML**

### **1. Откройте `chatroulette-app.html`**
```bash
# В браузере откройте файл
open chatroulette-app.html
```

### **2. Зарегистрируйте двух пользователей**
- `vasya` / `password123`
- `masha` / `password123`

### **3. Протестируйте Swirl**
- Войдите как `vasya`
- Нажмите "Найти собеседника"
- Войдите как `masha` в другом окне
- Нажмите "Найти собеседника"
- Они найдут друг друга!

### **4. Протестируйте функции**
- Отправьте сообщения
- Лайкните сообщения
- Редактируйте сообщения
- Загрузите медиа файлы

---

## 🔧 **Настройка окружения**

### **Переменные окружения**
Создайте файл `.env`:
```env
DATABASE_URL=postgres://username:password@localhost:5432/chatroulette?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-in-production
PORT=8080
```

### **База данных PostgreSQL**
```bash
# Установка PostgreSQL (Ubuntu/Debian)
sudo apt-get install postgresql postgresql-contrib

# Создание базы данных
sudo -u postgres createdb chatroulette
sudo -u postgres createuser --interactive
```

---

## 📚 **Документация**

### **Полная документация:**
- **[API-DOCUMENTATION.md](API-DOCUMENTATION.md)** - Полное руководство по API
- **[WEBSOCKET-EVENTS.md](WEBSOCKET-EVENTS.md)** - WebSocket события
- **[MEDIA-SUPPORT.md](MEDIA-SUPPORT.md)** - Работа с медиа файлами
- **[MESSAGE-STATUS.md](MESSAGE-STATUS.md)** - Статусы сообщений
- **[MESSAGE-LIKES.md](MESSAGE-LIKES.md)** - Система лайков

### **Безопасность:**
- **[SECURITY.md](SECURITY.md)** - Меры безопасности
- **[ARCHITECTURE-SECURITY.md](ARCHITECTURE-SECURITY.md)** - Архитектура безопасности

### **Специальные функции:**
- **[ANTI-DUPLICATE.md](ANTI-DUPLICATE.md)** - Предотвращение дублирования

---

## 🚀 **Готовые примеры**

### **JavaScript - Полный пример**
```javascript
// Подключение к API
const api = new ChatrouletteAPI();

// Регистрация
await api.register('vasya', 'password123');

// Поиск собеседника
const match = await api.findRandomUser();
console.log('Found:', match.user.username);

// Отправка сообщения
await api.sendMessage(match.chat_id, 'Hello!');

// WebSocket подключение
const ws = api.connectWebSocket(match.chat_id);
ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('New message:', data);
};
```

### **Python - Быстрый старт**
```python
import requests

# Регистрация
response = requests.post('http://localhost:8080/api/v1/auth/register', 
                        json={'username': 'vasya', 'password': 'password123'})
token = response.json()['token']

# Поиск собеседника
headers = {'Authorization': f'Bearer {token}'}
response = requests.get('http://localhost:8080/api/v1/chatroulette/find', headers=headers)
match = response.json()

# Отправка сообщения
requests.post(f'http://localhost:8080/api/v1/chats/{match["chat_id"]}/messages',
             headers=headers, json={'type': 'text', 'content': 'Hello!'})
```

---

## 🎯 **Что дальше?**

### **1. Интеграция с мобильным приложением**
- Используйте API endpoints
- Подключите WebSocket для реального времени
- Реализуйте push-уведомления

### **2. Расширение функциональности**
- Групповые чаты
- Каналы
- Боты
- Интеграции

### **3. Масштабирование**
- Кластеризация
- Load balancing
- CDN для медиа файлов

---

## ✅ **Готово к использованию!**

**Swirl Backend предоставляет:**
- 🎯 **Полный набор API** для мессенджера
- 🔌 **WebSocket** для реального времени
- 📁 **Медиа поддержка** для всех типов файлов
- ❤️ **Лайки и реакции** для интерактивности
- 🌪️ **Swirl** для поиска собеседников
- 🛡️ **Безопасность** и защита данных

**Начните разработку прямо сейчас!** 🚀✨
