package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/domain"
	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/errors"
	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/jwt"
)

func (h *HTTPHandler) RequireAuth() gin.HandlerFunc {
	return h.jwtManager.Interceptor()
}

func (h *HTTPHandler) GetUser(c *gin.Context) {
	userID, err := jwt.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.MsgAuthRequired,
		})
		return
	}

	user, err := h.service.Read(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": errors.MsgUserNotFound,
		})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *HTTPHandler) UpdateUser(c *gin.Context) {
	userID, err := jwt.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.MsgAuthRequired,
		})
		return
	}

	var updateUser struct {
		Username string `json:"username" binding:"omitempty,min=3,max=32"`
		Password string `json:"password" binding:"omitempty,min=6"`
	}

	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	user := domain.NewUser(userID, updateUser.Username, updateUser.Password)
	updatedUser, err := h.service.Update(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errors.MsgDatabaseOperation,
			"details": err.Error(),
		})
		return
	}

	updatedUser.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"message": errors.MsgUserUpdated,
		"user":    updatedUser,
	})
}

func (h *HTTPHandler) DeleteUser(c *gin.Context) {
	userID, err := jwt.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.MsgAuthRequired,
		})
		return
	}

	if _, err := h.service.Read(c.Request.Context(), userID); err != nil {
		c.JSON(404, gin.H{
			"error": errors.MsgUserNotFound,
		})
		return
	}

	if err := h.service.Delete(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errors.MsgDatabaseOperation,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": errors.MsgUserDeleted,
	})
}
