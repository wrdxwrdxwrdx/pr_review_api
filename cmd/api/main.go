package main

import (
	"pr_review_api/internal/handlers"
	"pr_review_api/internal/repository/inmemory"
	"pr_review_api/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация зависимостей
	userRepository := inmemory.NewUserRepository()
	userService := services.NewUserService(userRepository) // Предполагая, что у вас есть конструктор
	userHandler := handlers.NewUserHandler(userService)

	// Создание роутера Gin
	router := gin.Default()

	// Настройка middleware (опционально)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Настройка маршрутов
	api := router.Group("/api/v1")
	{
		api.POST("/users", userHandler.CreateUser)
		// Добавьте другие маршруты по мере необходимости
		// api.GET("/users", userHandler.GetUsers)
		// api.GET("/users/:id", userHandler.GetUserByID)
		// api.PUT("/users/:id", userHandler.UpdateUser)
		// api.DELETE("/users/:id", userHandler.DeleteUser)
	}

	// Запуск сервера
	router.Run(":8080") // localhost:8080
}
