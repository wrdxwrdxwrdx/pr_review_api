package handlers

import (
	"errors"
	"net/http"
	"pr_review_api/internal/domain/entities"
	customerrors "pr_review_api/internal/domain/errors"
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

	pr, err := h.PrService.Create(ctx, prJson.PullRequestId, prJson.PullRequestName, prJson.AuthorId)

	if err != nil {
		var domainErr *customerrors.DomainError
		if errors.As(err, &domainErr) {
			switch domainErr.Code {
			case customerrors.PRExists:
				c.JSON(http.StatusConflict, gin.H{
					"error": domainErr.Message,
				})
				return
			case customerrors.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Author not found",
				})
				return
			}
		}
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
			"error":   "Failed to create PullRequest",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, pr)
}

func (h *PrHandler) Reassign(c *gin.Context) {
	ctx := c.Request.Context()
	var prJson entities.ReassignPullRequestJson

	if err := c.ShouldBindJSON(&prJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON data",
			"details": err.Error(),
		})
		return
	}

	pr, err := h.PrService.Reassign(ctx, prJson.PullRequestId, prJson.OldReviewerId)

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
			c.JSON(http.StatusConflict, gin.H{
				"error": domainErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create PullRequest",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, pr)
}
