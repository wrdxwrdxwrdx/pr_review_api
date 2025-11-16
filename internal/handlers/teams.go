package handlers

import (
	"errors"
	"net/http"
	"pr_review_api/internal/domain/entities"
	customerrors "pr_review_api/internal/domain/errors"
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

	if err := c.ShouldBindJSON(&teamJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON data",
			"details": err.Error(),
		})
		return
	}

	team, err := h.TeamService.Create(ctx, teamJson.TeamName, teamJson.Members)

	if err != nil {
		var domainErr *customerrors.DomainError
		if errors.As(err, &domainErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": domainErr.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, team)
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	ctx := c.Request.Context()
	teamNameQuery := c.Query("TeamNameQuery")

	if teamNameQuery == "" {
		c.JSON(400, gin.H{
			"error": "No teamNameQuery",
		})
		return
	}

	team, err := h.TeamService.GetByID(ctx, teamNameQuery)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Team was not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, team)
}
