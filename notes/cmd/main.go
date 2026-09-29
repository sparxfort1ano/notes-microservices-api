package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/logger"
	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/repository/mongo"
	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/repository/redis"
	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/service"
	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/transport/http/handler"
	"github.com/sparxfort1ano/notes-microservices-api/notes/internal/transport/http/server"
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

	log.Debug("initializing mongodb client")
	db, err := mongo.NewDatabase(mongo.NewConfigMust())
	if err != nil {
		log.Fatal("failed to initialize mongodb client", zap.Error(err))
		return
	}

	redisCfg := redis.NewConfigMust()
	var svc service.Service
	if redisCfg.Enabled {
		log.Debug("initializing redis client")
		cache, err := redis.NewCachedDatabase(redis.NewConfigMust())
		if err != nil {
			log.Fatal("failed to initialize redis client", zap.Error(err))
			return
		}
		svc = service.NewCachedService(cache, service.NewService(db))
	} else {
		svc = service.NewService(db)
	}
	defer svc.Close()

	log.Debug("initializing HTTP server")
	r := handler.HTTPRouter(handler.NewHTTPHandler(svc), log)
	s := server.NewHTTPServer(server.NewConfigMust(), r, log)

	if err := s.Run(ctx); err != nil {
		log.Error("HTTP server run error", zap.Error(err))
	}
}
