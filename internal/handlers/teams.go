package handlers

import (
	"net/http"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/services"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	TeamService services.TeamService
}

func NewTeamHandler(teamService services.TeamService) *TeamHandler {
	return &TeamHandler{
		TeamService: teamService,
	}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	ctx := c.Request.Context()
	var teamJson entities.Team

	// body, _ := c.GetRawData()
	// fmt.Printf("Raw request body: %s\n", string(body))

	if err := c.ShouldBindJSON(&teamJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON data",
			"details": err.Error(),
		})
		return
	}

	team, err := h.TeamService.Create(ctx, teamJson.TeamName, teamJson.Members)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, team)
}
