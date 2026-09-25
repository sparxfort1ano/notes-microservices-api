package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/domain"
	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/errors"
	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/jwt"
)

func (h *HTTPHandler) RefreshToken(c *gin.Context) {
	var refreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	userID, err := h.jwtManager.ValidateToken(refreshRequest.RefreshToken, jwt.RefreshTokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.MsgRefreshToken,
		})
		return
	}

	if _, err := h.service.Read(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": errors.MsgUserNotFound,
		})
		return
	}

	accessToken, refreshToken, err := h.jwtManager.GenerateTokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errors.MsgTokenGeneration,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       errors.MsgTokensRefreshed,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *HTTPHandler) RegisterUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	createdUser, err := h.service.Create(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.MsgUserCreation,
		})
		return
	}

	createdUser.Password = ""
	c.JSON(http.StatusCreated, gin.H{
		"message": errors.MsgUserRegistered,
		"user":    createdUser,
	})
}

func (h *HTTPHandler) LoginUser(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	user, err := h.service.Authenticate(
		c.Request.Context(),
		loginRequest.Username,
		loginRequest.Password,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.MsgInvalidCredentials,
		})
		return
	}

	user.Password = ""

	accessToken, refreshToken, err := h.jwtManager.GenerateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errors.MsgTokenGeneration,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       errors.MsgLoginSuccess,
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
