package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	AdminTokenHeader = "X-Admin-Token"
)

type AdminMiddleware struct {
	adminToken string
}

func NewAdminMiddleware(adminToken string) *AdminMiddleware {
	return &AdminMiddleware{
		adminToken: adminToken,
	}
}

func (m *AdminMiddleware) AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminToken := c.GetHeader(AdminTokenHeader)
		if adminToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "admin token is required",
			})
			c.Abort()
			return
		}

		if adminToken != m.adminToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid admin token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AdminMiddleware) AdminOrUser(jwtMiddleware *AuthMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminToken := c.GetHeader(AdminTokenHeader)
		if adminToken != "" && adminToken == m.adminToken {
			c.Set("is_admin", true)
			c.Next()
			return
		}

		jwtMiddleware.Authenticate()(c)
		if c.IsAborted() {
			return
		}

		c.Set("is_admin", false)
		c.Next()
	}
}
