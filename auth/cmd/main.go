package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/logger"
	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/service"
	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/transport/http/handler"
	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/transport/http/server"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log, err := logger.NewLogger(logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to initialize app logger", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Debug("initializing postgres connection pool")
	service, err := service.NewService()
	if err != nil {
		log.Fatal("failed to initialize postgres connection pool", zap.Error(err))
	}

	log.Debug("initializing HTTP server")
	r := handler.HTTPRouter(handler.NewHTTPHandler(service))
	s := server.NewHTTPServer(server.NewConfigMust(), r, log)

	if err := s.Run(ctx); err != nil {
		log.Error("HTTP server run error", zap.Error(err))
	}
}
