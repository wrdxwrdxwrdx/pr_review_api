package main

import (
	"context"
	"fmt"
	"log"
	"pr_review_api/internal/config"
	"pr_review_api/internal/handlers"
	"pr_review_api/internal/repository/postgres"
	"pr_review_api/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Инициализация БД
	db := postgres.NewDB()

	// Формируем connection string
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.SSLMode,
	)

	// Подключаемся к БД
	err := db.Connect(connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация схемы БД
	ctx := context.Background()
	err = db.InitSchema(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	userRepository := postgres.NewUserRepository(db.Pool)
	teamRepository := postgres.NewTeamRepository(db.Pool)
	prRepository := postgres.NewPRRepository(db.Pool)

	teamService := services.NewTeamService(teamRepository)
	userService := services.NewUserService(userRepository, prRepository)
	prService := services.NewPrService(prRepository, userRepository)

	userHandler := handlers.NewUserHandler(userService)
	teamHandler := handlers.NewTeamHandler(teamService)
	prHandler := handlers.NewPrHandler(prService)

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")
	{
		api.POST("/users/setIsActive", userHandler.SetIsActive)
		api.GET("/users/getReview", userHandler.GetReview)
		api.POST("/team/add", teamHandler.CreateTeam)
		api.POST("/pullRequest/create", prHandler.CreatePr)
		api.POST("/pullRequest/merge", prHandler.MergePr)
		api.POST("/pullRequest/reassign", prHandler.Reassign)
		api.GET("/team/get", teamHandler.GetTeam)
	}

	router.Run(":8080")
}
