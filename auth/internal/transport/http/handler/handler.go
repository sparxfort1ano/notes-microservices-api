package handler

import (
	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/jwt"
	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/service"
)

type HTTPHandler struct {
	service    service.Service
	jwtManager *jwt.Manager
}

func NewHTTPHandler(service service.Service) *HTTPHandler {
	return &HTTPHandler{
		service:    service,
		jwtManager: jwt.NewManager(),
	}
}
