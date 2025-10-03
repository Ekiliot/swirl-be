# Инструкции по установке

## Установка Go

### Windows
1. Скачайте Go с официального сайта: https://golang.org/dl/
2. Установите Go, следуя инструкциям установщика
3. Проверьте установку: `go version`

### Альтернативно через Chocolatey
```powershell
choco install golang
```

### Альтернативно через Scoop
```powershell
scoop install go
```

## Установка PostgreSQL

### Windows
1. Скачайте PostgreSQL с официального сайта: https://www.postgresql.org/download/windows/
2. Установите PostgreSQL
3. Создайте базу данных:
```sql
CREATE DATABASE chatroulette;
```

### Альтернативно через Chocolatey
```powershell
choco install postgresql
```

## Запуск проекта

1. Установите зависимости:
```bash
go mod tidy
```

2. Скопируйте `env.example` в `.env` и настройте переменные

3. Запустите приложение:
```bash
go run main.go
```

## Или используйте Docker

```bash
docker-compose up -d
```
