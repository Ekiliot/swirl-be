# 📸 Поддержка медиа файлов в Chatroulette Backend

## ✅ **Что уже работает:**

### **📁 Загрузка файлов:**
- **Фото:** JPEG, PNG, GIF, WebP (до 10MB)
- **Видео:** MP4, WebM, QuickTime (до 100MB)
- **Аудио:** MP3, WAV, OGG, M4A (до 25MB)
- **Безопасность:** Проверка типов файлов, защита от path traversal
- **Организация:** Файлы хранятся в папках по пользователям

### **💬 Типы сообщений:**
- **text** - текстовые сообщения
- **image** - изображения
- **video** - видео
- **voice** - голосовые сообщения
- **sticker** - стикеры
- **gif** - GIF анимации

## 🚀 **API для работы с медиа:**

### **1. Загрузка файла:**
```bash
POST /api/v1/upload
Content-Type: multipart/form-data

# Параметры:
file: [файл] - файл для загрузки

# Ответ:
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

### **2. Получение файла:**
```bash
GET /uploads/{user_id}/{file_name}

# Пример:
GET /uploads/123e4567-e89b-12d3-a456-426614174000/photo.jpg
```

### **3. Удаление файла:**
```bash
DELETE /api/v1/uploads/{file_name}
Authorization: Bearer {token}

# Ответ:
{
  "message": "File deleted successfully"
}
```

### **4. Отправка медиа сообщения:**
```bash
POST /api/v1/chats/{chat_id}/messages
Authorization: Bearer {token}
Content-Type: application/json

{
  "type": "image",
  "content": "Описание изображения",
  "media_url": "/uploads/user_id/filename.jpg"
}
```

## 📋 **Поддерживаемые форматы:**

### **🖼️ Изображения:**
- **JPEG** (.jpg, .jpeg) - до 10MB
- **PNG** (.png) - до 10MB
- **GIF** (.gif) - до 10MB
- **WebP** (.webp) - до 10MB

### **🎥 Видео:**
- **MP4** (.mp4) - до 100MB
- **WebM** (.webm) - до 100MB
- **QuickTime** (.mov) - до 100MB

### **🎵 Аудио:**
- **MP3** (.mp3) - до 25MB
- **WAV** (.wav) - до 25MB
- **OGG** (.ogg) - до 25MB
- **M4A** (.m4a) - до 25MB

## 🛡️ **Безопасность:**

### **Проверки:**
- ✅ **Тип файла** - только разрешенные MIME типы
- ✅ **Размер файла** - ограничения по типу медиа
- ✅ **Path traversal** - защита от `../` атак
- ✅ **Уникальные имена** - UUID для предотвращения конфликтов
- ✅ **Изоляция пользователей** - файлы в отдельных папках

### **Ограничения:**
- **Изображения:** максимум 10MB
- **Видео:** максимум 100MB
- **Аудио:** максимум 25MB
- **Общий лимит:** 50MB для других типов

## 🔧 **Как использовать:**

### **1. Загрузка файла:**
```javascript
const formData = new FormData();
formData.append('file', fileInput.files[0]);

const response = await fetch('/api/v1/upload', {
    method: 'POST',
    headers: {
        'Authorization': 'Bearer ' + token
    },
    body: formData
});

const result = await response.json();
console.log('File uploaded:', result.file_url);
```

### **2. Отправка медиа сообщения:**
```javascript
const message = {
    type: 'image',
    content: 'Посмотри на это фото!',
    media_url: result.file_url
};

await fetch(`/api/v1/chats/${chatId}/messages`, {
    method: 'POST',
    headers: {
        'Authorization': 'Bearer ' + token,
        'Content-Type': 'application/json'
    },
    body: JSON.stringify(message)
});
```

### **3. Отображение медиа в чате:**
```javascript
function displayMessage(message) {
    if (message.type === 'image') {
        return `<img src="${message.media_url}" alt="Image" style="max-width: 300px;">`;
    } else if (message.type === 'video') {
        return `<video controls style="max-width: 300px;"><source src="${message.media_url}"></video>`;
    } else if (message.type === 'voice') {
        return `<audio controls><source src="${message.media_url}"></audio>`;
    }
    return message.content;
}
```

## 🚀 **Возможные улучшения:**

### **1. Обработка медиа:**
- **Сжатие изображений** - автоматическое уменьшение размера
- **Превью видео** - генерация миниатюр
- **Метаданные** - извлечение EXIF, длительности
- **Конвертация** - автоматическое преобразование форматов

### **2. CDN интеграция:**
- **Облачное хранилище** - AWS S3, Google Cloud Storage
- **CDN** - быстрая доставка файлов
- **Кэширование** - оптимизация загрузки

### **3. Дополнительные форматы:**
- **Документы** - PDF, DOC, TXT
- **Архивы** - ZIP, RAR
- **3D модели** - OBJ, STL

### **4. Безопасность:**
- **Сканирование вирусов** - проверка загруженных файлов
- **Водяные знаки** - защита авторских прав
- **Шифрование** - защита конфиденциальных файлов

## 📊 **Статистика использования:**

### **Мониторинг:**
- **Размер хранилища** - общий объем файлов
- **Популярные форматы** - какие типы загружают чаще
- **Активность** - количество загрузок в день
- **Ошибки** - неудачные попытки загрузки

## ✅ **Заключение:**

**Система медиа файлов полностью функциональна!**
- 🎯 **Поддержка всех основных форматов**
- 🛡️ **Безопасность и защита**
- 📱 **Готово для мобильных приложений**
- 🚀 **Легко расширяется**

**Можно загружать фото, видео и голосовые сообщения прямо сейчас!** 📸🎥🎵
