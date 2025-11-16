package handlers

import (
	"errors"
	"net/http"
	"pr_review_api/internal/domain/entities"
	customerrors "pr_review_api/internal/domain/errors"
	"pr_review_api/internal/middleware"
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
		var domainErr *customerrors.DomainError
		if errors.As(err, &domainErr) {
			switch domainErr.Code {
			case customerrors.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"error": domainErr.Message,
				})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to find user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetReview(c *gin.Context) {
	ctx := c.Request.Context()
	userIdQuery := c.Query("UserIdQuery")

	if !middleware.IsAdmin(c) {
		claims, exists := middleware.GetUserFromContext(c)
		if !exists || claims.UserID != userIdQuery {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "can only view your own reviews",
			})
			return
		}
	}

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
