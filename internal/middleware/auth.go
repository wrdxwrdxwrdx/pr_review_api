package middleware

import (
	"net/http"
	"pr_review_api/pkg/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserContextKey      = "user"
)

type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		claims, err := m.jwtManager.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set(UserContextKey, claims)
		c.Next()
	}
}

func GetUserFromContext(c *gin.Context) (*auth.Claims, bool) {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return nil, false
	}

	claims, ok := user.(*auth.Claims)
	return claims, ok
}

func (m *AuthMiddleware) AuthenticateHandler() gin.HandlerFunc {
	return m.Authenticate()
}

func IsAdmin(c *gin.Context) bool {
	isAdmin, exists := c.Get("is_admin")
	return exists && isAdmin.(bool)
}
