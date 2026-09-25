package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenStr  = "accessToken"
	RefreshTokenStr = "refreshToken"
)

type Manager struct {
	cfg config
}

func NewManager() *Manager {
	return &Manager{
		cfg: NewConfigMust(),
	}
}

// GenerateTokens creates two JWTs for a user:
// - Access Token: for authorising API requests (short-living),
// - Refresh Token: to renew the access token (long-lived).
func (m *Manager) GenerateTokens(id int) (access, refresh string, err error) {
	accessToken, err := m.generateToken(id, AccessTokenStr, m.cfg.AccessTokenExpiration)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w: %v", accessToken, ErrTokenGeneration, err)
	}

	refreshToken, err := m.generateToken(id, RefreshTokenStr, m.cfg.RefreshTokenExpiration)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w: %v", RefreshTokenStr, ErrTokenGeneration, err)
	}

	return accessToken, refreshToken, nil
}

func (m *Manager) generateToken(id int, tokenType string, expDur time.Duration) (string, error) {
	now := time.Now()
	expTime := now.Add(expDur)

	claims := jwt.MapClaims{
		"id":   id,
		"type": tokenType,
		"iat":  now.Unix(),
		"exp":  expTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.cfg.SecretKey))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidSignature, err)
	}

	return tokenString, nil
}

// ValidateToken checks the signature, expiry date and token type.
func (m *Manager) ValidateToken(tokenString, tokenType string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrInvalidSignature, t.Header["alg"])
		}

		return []byte(m.cfg.SecretKey), nil
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
