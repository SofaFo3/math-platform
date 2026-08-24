package main

import (
	"log"
	"math-platform/internal/config"
	"math-platform/internal/handlers"
	"math-platform/internal/repository"
	"math-platform/internal/services"
	"math-platform/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загружаем конфиг
	cfg := config.Load()

	// Подключаемся к БД
	db, err := database.Connect(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Инициализируем репозитории
	userRepo := repository.NewUserRepository(db)

	// Инициализируем сервисы
	jwtSvc := services.NewJWTService(cfg.JWTSecret)

	// Инициализируем хендлеры
	authHandler := handlers.NewAuthHandler(userRepo, jwtSvc)

	// Настраиваем роутер
	r := gin.Default()

	// Публичные эндпоинты
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Защищенные эндпоинты (требуют JWT)
	protected := r.Group("/api")
	protected.Use(authHandler.AuthMiddleware())
	{
		protected.GET("/profile", authHandler.GetProfile)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Запускаем сервер
	log.Printf(" Server running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
