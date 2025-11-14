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

	teamService := services.NewTeamService(teamRepository)
	userService := services.NewUserService(userRepository)

	userHandler := handlers.NewUserHandler(userService)
	teamHandler := handlers.NewTeamHandler(teamService)

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")
	{
		api.POST("/users", userHandler.CreateUser)
		api.POST("/team/add", teamHandler.CreateTeam)
		api.GET("/team/get", teamHandler.GetTeam)
	}

	router.Run(":8080")
}
