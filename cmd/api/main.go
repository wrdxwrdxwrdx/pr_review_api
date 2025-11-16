package main

import (
	"pr_review_api/internal/handlers"
	"pr_review_api/internal/repository/inmemory"
	"pr_review_api/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	userRepository := inmemory.NewUserRepository()
	teamRepository := inmemory.NewTeamRepository()
	prRepository := inmemory.NewPrRepository()

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
