package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/logger"
	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/transport/http/middleware"
)

func HTTPRouter(h *HTTPHandler, log *logger.Logger) *gin.Engine {
	r := gin.New()
	r.Use(
		gin.Recovery(),
		middleware.RequestID(),
		middleware.TraceAndLog(log),
	)

	notes := r.Group("/notes")
	notes.Use(h.jwtInterceptor.RequireAuth())
	{
		notes.POST("", h.CreateNote)
		notes.GET("", h.GetAllNotes)
		notes.GET("/:id", h.GetNote)
		notes.PATCH("/:id", h.UpdateNote)
		notes.DELETE("/:id", h.DeleteNote)
	}

	return r
}
