package handlers

import (
	"net/http"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) SetIsActive(c *gin.Context) {
	ctx := c.Request.Context()
	var userJson entities.UserIsActiveJson

	if err := c.ShouldBindJSON(&userJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON data",
			"details": err.Error(),
		})
		return
	}

	user, err := h.userService.SetIsActive(ctx, userJson.UserId, userJson.IsActive)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetReview(c *gin.Context) {
	ctx := c.Request.Context()
	userIdQuery := c.Query("UserIdQuery")

	if userIdQuery == "" {
		c.JSON(400, gin.H{
			"error": "No teamNameQuery",
		})
		return
	}

	prs, err := h.userService.GetReview(ctx, userIdQuery)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, prs)
}
