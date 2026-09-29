package handler

import (
	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/jwt"
	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/service"
)

type HTTPHandler struct {
	jwtInterceptor *jwt.Interceptor
	service        service.Service
}

func NewHTTPHandler(service service.Service) *HTTPHandler {
	return &HTTPHandler{
		jwtInterceptor: jwt.NewInterceptor(),
		service:        service,
	}
}
