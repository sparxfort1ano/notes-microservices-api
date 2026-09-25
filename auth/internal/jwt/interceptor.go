package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *Manager) Interceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := m.ExtractTokenFromHeader(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": MsgTokenRequired,
			})
			c.Abort()
			return
		}

		userID, err := m.ValidateToken(tokenString, AccessTokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": MsgInvalidToken,
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

const bearerPrefix = "Bearer "

func (m *Manager) ExtractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("%w", ErrMissingAuthHeader)
	}

	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("%w", ErrInvalidAuthFormat)
	}

	return authHeader[len(bearerPrefix):], nil
}

func GetCurrentUserID(c *gin.Context) (int, error) {
	userID, exists := c.Get("user_id")
	id, ok := userID.(int)
	if !ok || !exists {
		return 0, ErrMissingUserID
	}

	return id, nil
}
