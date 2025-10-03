package main

import (
	"log"
	"os"
	"time"

	"swirl-backend/internal/config"
	"swirl-backend/internal/database"
	"swirl-backend/internal/handlers"
	"swirl-backend/internal/middleware"
	"swirl-backend/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Инициализируем конфигурацию
	cfg := config.Load()

	// Подключаемся к базе данных
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Автомиграции
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Инициализируем WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Инициализируем SearchQueueHandler для очистки
	searchHandler := handlers.NewSearchQueueHandler(db)
	
	// Очищаем очередь при запуске сервера
	searchHandler.CleanupInactiveUsers()
	
	// Запускаем периодическую очистку очереди (каждые 30 секунд)
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			searchHandler.CleanupInactiveUsers()
		}
	}()

	// Создаем Gin роутер
	r := gin.Default()

	// Middleware безопасности
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RateLimiting())
	r.Use(middleware.InputValidation())

	// Настройка CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // В продакшене указать конкретные домены
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Инициализируем хендлеры
	authHandler := handlers.NewAuthHandler(db)
	chatHandler := handlers.NewChatHandler(db, hub)
	messageHandler := handlers.NewMessageHandler(db, hub)
	uploadHandler := handlers.NewUploadHandler("./uploads")
	chatrouletteHandler := handlers.NewChatrouletteHandler(db)

	// Публичные роуты
	api := r.Group("/api/v1")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
	}

	// Защищенные роуты
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Пользователи
		protected.GET("/profile", authHandler.GetProfile)
		protected.PUT("/profile", authHandler.UpdateProfile)
		protected.POST("/profile/online", authHandler.UpdateOnlineStatus)
		protected.GET("/users/:id/profile", authHandler.GetPublicProfile)

		// Чаты
		protected.GET("/chats", chatHandler.GetChats)
		protected.POST("/chats", chatHandler.CreateChat)
		protected.GET("/chats/:id", chatHandler.GetChat)
		protected.DELETE("/chats/:id", chatHandler.DeleteChat)

		// Сообщения
		protected.GET("/chats/:id/messages", messageHandler.GetMessages)
		protected.POST("/chats/:id/messages", messageHandler.SendMessage)
		protected.DELETE("/messages/:id", messageHandler.DeleteMessage)
		protected.PUT("/messages/:id/read", messageHandler.MarkMessageAsRead)
		protected.PUT("/messages/:id/edit", messageHandler.EditMessage)
		protected.GET("/messages/:id/status", messageHandler.GetMessageStatus)
		
		// Лайки сообщений
		protected.POST("/messages/:id/like", messageHandler.LikeMessage)
		protected.DELETE("/messages/:id/like", messageHandler.UnlikeMessage)
		protected.GET("/messages/:id/likes", messageHandler.GetMessageLikes)

		// Загрузка файлов
		protected.POST("/upload", uploadHandler.UploadFile)
		protected.DELETE("/uploads/:file_name", uploadHandler.DeleteFile)

		// Swirl (случайные встречи)
		protected.GET("/swirl/find", chatrouletteHandler.FindRandomUser)
		protected.POST("/swirl/:id/save", chatrouletteHandler.SaveChat)
		protected.POST("/swirl/:id/skip", chatrouletteHandler.SkipUser)
		protected.POST("/swirl/activity", chatrouletteHandler.UpdateSearchActivity)
		protected.GET("/swirl/status", chatrouletteHandler.GetQueueStatus)
		protected.DELETE("/swirl/clear", chatrouletteHandler.ClearQueue)
	}

	// WebSocket для real-time общения (без middleware авторизации)
	api.GET("/ws", chatHandler.HandleWebSocket)

	// Статические файлы (загрузки)
	r.GET("/uploads/:user_id/:file_name", uploadHandler.GetFile)

	// Запускаем сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
