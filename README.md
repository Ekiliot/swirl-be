# 🌪️ Swirl Backend

Бекенд для приложения Swirl - случайные встречи с возможностью сохранения чатов и общения как в обычном мессенджере.

## Возможности

- 🔐 Аутентификация и регистрация пользователей
- 💬 Создание и управление чатами
- 📱 Real-time общение через WebSocket
- 🎭 Поддержка различных типов сообщений:
  - Текстовые сообщения
  - Стикеры
  - GIF
  - Голосовые сообщения
  - Видео сообщения
  - Ответы на сообщения
- 💾 Сохранение чатов из Swirl в личные чаты

## Технологии

- **Go** - основной язык
- **Gin** - веб-фреймворк
- **GORM** - ORM для работы с базой данных
- **PostgreSQL** - база данных
- **WebSocket** - real-time общение
- **JWT** - аутентификация
- **Docker** - контейнеризация

## Установка и запуск

### С помощью Docker Compose (рекомендуется)

1. Клонируйте репозиторий
2. Скопируйте `env.example` в `.env` и настройте переменные
3. Запустите:
```bash
docker-compose up -d
```

### Локальная разработка

1. Установите PostgreSQL
2. Создайте базу данных:
```sql
CREATE DATABASE swirl;
```

3. Установите зависимости:
```bash
go mod tidy
```

4. Скопируйте `env.example` в `.env` и настройте переменные

5. Запустите приложение:
```bash
go run main.go
```

## API Endpoints

### Аутентификация
- `POST /api/v1/register` - регистрация
- `POST /api/v1/login` - авторизация

### Профиль (требует авторизации)
- `GET /api/v1/profile` - получить профиль
- `PUT /api/v1/profile` - обновить профиль

### Чаты (требует авторизации)
- `GET /api/v1/chats` - получить список чатов
- `POST /api/v1/chats` - создать чат
- `GET /api/v1/chats/:id` - получить чат
- `DELETE /api/v1/chats/:id` - удалить чат

### Сообщения (требует авторизации)
- `GET /api/v1/chats/:id/messages` - получить сообщения чата
- `POST /api/v1/chats/:id/messages` - отправить сообщение
- `DELETE /api/v1/messages/:id` - удалить сообщение

### WebSocket
- `GET /api/v1/ws?chat_id=:id` - подключение к real-time общению

## Типы сообщений

- `text` - текстовое сообщение
- `sticker` - стикер
- `gif` - GIF
- `voice` - голосовое сообщение
- `video` - видео сообщение
- `image` - изображение
- `reply` - ответ на сообщение

## Структура проекта

```
├── internal/
│   ├── config/          # Конфигурация
│   ├── database/        # Подключение к БД
│   ├── handlers/        # HTTP хендлеры
│   ├── middleware/      # Middleware
│   ├── models/          # Модели данных
│   └── websocket/       # WebSocket логика
├── main.go              # Точка входа
├── go.mod               # Зависимости
├── Dockerfile           # Docker образ
└── docker-compose.yml   # Docker Compose
```

## Переменные окружения

- `DATABASE_URL` - URL подключения к PostgreSQL
- `JWT_SECRET` - секретный ключ для JWT
- `PORT` - порт сервера (по умолчанию 8080)
- `ALLOWED_ORIGINS` - разрешенные CORS домены

## Разработка

Для разработки рекомендуется использовать hot reload:

```bash
go install github.com/cosmtrek/air@latest
air
```

## Лицензия

MIT
