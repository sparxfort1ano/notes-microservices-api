package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sparxfort1ano/notes-microservices-api/auth/internal/logger"
	"go.uber.org/zap"
)

type HTTPServer struct {
	router *gin.Engine
	cfg    config
	log    *logger.Logger
}

func NewHTTPServer(
	cfg config,
	r *gin.Engine,
	log *logger.Logger,
) *HTTPServer {
	return &HTTPServer{
		router: r,
		cfg:    cfg,
		log:    log,
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:         s.cfg.Addr,
		Handler:      s.router,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
		IdleTimeout:  s.cfg.ShutdownTimeout,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Info("start HTTP server", zap.String("addr", s.cfg.Addr))

		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		s.log.Info("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.cfg.ShutdownTimeout,
		)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			srv.Close()
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Info("HTTP server stopped")
	}

	return nil
}
