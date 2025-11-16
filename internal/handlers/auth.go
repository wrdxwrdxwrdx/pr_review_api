package handlers

import (
	"fmt"
	"net/http"
	customerrors "pr_review_api/internal/domain/errors"
	"pr_review_api/internal/middleware"
	"pr_review_api/internal/services"
	"pr_review_api/pkg/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService services.UserService
	jwtManager  *auth.JWTManager
}

func NewAuthHandler(userService services.UserService, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtManager:  jwtManager,
	}
}

type LoginRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	UserID  string `json:"user_id"`
	Expires int64  `json:"expires_in"`
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	fmt.Println(1)
	exists, err := h.userService.Exist(ctx, req.UserID)
	if exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "User already created",
		})
		return
	}
	fmt.Println(2)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println(3)

	user, err := h.userService.Create(ctx, req.UserID, req.Username, "EMPTYTEAM", true)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println(4)

	token, err := h.jwtManager.GenerateToken(user.UserId, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}
	ctx.JSON(http.StatusOK, LoginResponse{
		Token:   token,
		UserID:  user.UserId,
		Expires: int64(h.jwtManager.TokenDuration.Seconds()),
	})
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := h.userService.GetById(ctx.Request.Context(), req.UserID)
	if err != nil {
		if domainErr, ok := err.(*customerrors.DomainError); ok && domainErr.Code == customerrors.NotFound {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid credentials",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if user.Username != req.Username {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	if !user.IsActive {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "user account is inactive",
		})
		return
	}

	token, err := h.jwtManager.GenerateToken(user.UserId, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}
	ctx.JSON(http.StatusOK, LoginResponse{
		Token:   token,
		UserID:  user.UserId,
		Expires: int64(h.jwtManager.TokenDuration.Seconds()),
	})
}

type MeResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

func (h *AuthHandler) Me(ctx *gin.Context) {
	claims, exists := middleware.GetUserFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found in context",
		})
		return
	}
	fmt.Println(claims)
	fmt.Println(claims.UserID)
	user, err := h.userService.GetById(ctx.Request.Context(), claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get user data",
		})
		return
	}

	ctx.JSON(http.StatusOK, MeResponse{
		UserID:   user.UserId,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	})
}
