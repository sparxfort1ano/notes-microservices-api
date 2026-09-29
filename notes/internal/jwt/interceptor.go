package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenStr  = "accessToken"
	RefreshTokenStr = "refreshToken"
)

type Interceptor struct {
	cfg config
}

func NewInterceptor() *Interceptor {
	return &Interceptor{
		cfg: NewConfigMust(),
	}
}

func (i *Interceptor) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := i.extractTokenFromHeader(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": MsgTokenRequired,
			})
			c.Abort()
			return
		}

		userID, err := i.validateToken(tokenString, AccessTokenStr)
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

func (i *Interceptor) extractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("%w", ErrMissingAuthHeader)
	}

	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("%w", ErrInvalidAuthFormat)
	}

	return authHeader[len(bearerPrefix):], nil
}

func (i *Interceptor) validateToken(tokenString, tokenType string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrInvalidSignature, t.Header["alg"])
		}

		return []byte(i.cfg.SecretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, fmt.Errorf("%w: %v", ErrTokenExpired, err)
		}

		return 0, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claimType, exists := claims["type"].(string); !exists || claimType != tokenType {
			return 0, fmt.Errorf("%w: expected %s, got %s", ErrInvalidTokenType, tokenType, claims["type"])
		}

		idValue, exists := claims["id"].(float64)
		if !exists {
			return 0, fmt.Errorf("%w", ErrMissingUserID)
		}

		return int(idValue), nil
	}

	return 0, fmt.Errorf("%w", ErrInvalidToken)
}

func GetCurrentUserID(c *gin.Context) (int, error) {
	userID, exists := c.Get("user_id")
	id, ok := userID.(int)
	if !ok || !exists {
		return 0, ErrMissingUserID
	}

	return id, nil
}
