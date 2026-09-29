package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/logger"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-ID"
	requestIDStr    = "request_id"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(requestIDHeader)
		if reqID == "" {
			reqID = uuid.NewString()
		}

		c.Header(requestIDHeader, reqID)
		c.Set(requestIDStr, reqID)

		c.Next()
	}
}

func TraceAndLog(baseLog *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetString(requestIDStr)

		log := baseLog.With(
			zap.String("request_id", reqID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)

		c.Request = c.Request.WithContext(logger.IntoContext(c.Request.Context(), log))

		before := time.Now()

		c.Next()

		log.Debug("done HTTP request",
			zap.Int("status_code", c.Writer.Status()),
			zap.Duration("latency", time.Since(before)),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
