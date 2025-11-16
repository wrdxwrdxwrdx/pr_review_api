package main

import (
	"context"
	"fmt"
	"log"
	"pr_review_api/internal/config"
	"pr_review_api/internal/handlers"
	"pr_review_api/internal/middleware"
	"pr_review_api/internal/repository/postgres"
	"pr_review_api/internal/services"
	"pr_review_api/pkg/auth"

	"github.com/gin-gonic/gin"
)

func initJWT(cfg *config.Config) *auth.JWTManager {
	jwtManager := auth.NewJWTManager(
		cfg.JWT.SecretKey,
		cfg.JWT.ExpiresIn,
		cfg.JWT.Issuer,
	)
	return jwtManager
}

func main() {
	cfg := config.Load()

	db := postgres.NewDB()

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.SSLMode,
	)

	err := db.Connect(connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	err = db.InitSchema(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}
	jwtManager := initJWT(cfg)

	userRepository := postgres.NewUserRepository(db.Pool)
	teamRepository := postgres.NewTeamRepository(db.Pool, *userRepository)
	prRepository := postgres.NewPRRepository(db.Pool)

	teamService := services.NewTeamService(teamRepository)
	userService := services.NewUserService(userRepository, prRepository)
	prService := services.NewPrService(prRepository, userRepository)

	authHandler := handlers.NewAuthHandler(userService, jwtManager)
	userHandler := handlers.NewUserHandler(userService)
	teamHandler := handlers.NewTeamHandler(teamService)
	prHandler := handlers.NewPrHandler(prService)

	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	router := gin.Default()

	public := router.Group("/api/v1")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		// public.GET("/health", healthHandler.HealthCheck)
	}

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	protected := router.Group("/api/v1")
	protected.Use(authMiddleware.Authenticate())
	{
		// User routes
		protected.POST("/users/setIsActive", userHandler.SetIsActive)
		protected.GET("/users/getReview", userHandler.GetReview)
		protected.GET("/auth/me", authHandler.Me)

		// Team routes
		protected.POST("/team/add", teamHandler.CreateTeam)
		protected.GET("/team/get", teamHandler.GetTeam)

		// PR routes
		protected.POST("/pullRequest/create", prHandler.CreatePr)
		protected.POST("/pullRequest/merge", prHandler.MergePr)
		protected.POST("/pullRequest/reassign", prHandler.Reassign)
	}

	addr := ":" + cfg.ServerPort
	router.Run(addr)
}
