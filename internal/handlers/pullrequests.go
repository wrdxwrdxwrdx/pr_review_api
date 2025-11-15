package handlers

import (
	"net/http"
	"pr_review_api/internal/domain/entities"
	"pr_review_api/internal/services"

	"github.com/gin-gonic/gin"
)

type PrHandler struct {
	PrService services.PrService
}

func NewPrHandler(prService services.PrService) *PrHandler {
	return &PrHandler{
		PrService: prService,
	}
}

func (h *PrHandler) CreatePr(c *gin.Context) {
	ctx := c.Request.Context()
	var prJson entities.PullRequestJson

	if err := c.ShouldBindJSON(&prJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON data",
			"details": err.Error(),
		})
		return
	}

	pr, err := h.PrService.Create(ctx, prJson.PullRequestId, prJson.PullRequestId, prJson.AuthorId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create PullRequest",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, pr)
}

func (h *PrHandler) MergePr(c *gin.Context) {
	ctx := c.Request.Context()
	var prJson entities.MergePullRequestJson

	if err := c.ShouldBindJSON(&prJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON data",
			"details": err.Error(),
		})
		return
	}

	pr, err := h.PrService.Merge(ctx, prJson.PullRequestId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create PullRequest",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, pr)
}
