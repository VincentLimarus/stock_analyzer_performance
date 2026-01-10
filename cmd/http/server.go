package http

import (
	configs "VincentLimarus/stock-analyzer-performance/config"
	delivery "VincentLimarus/stock-analyzer-performance/delivery/http"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type IServer interface {
	Serve(ctx context.Context)
}

type Server struct {
	router Router
	srv    *http.Server
}

func NewServer(delivery delivery.IRegistry, middlewareLimiter gin.HandlerFunc) *Server {
	return &Server{
		router: NewRouter(
			delivery,
			middlewareLimiter,
		),
		srv: &http.Server{
			Addr:              fmt.Sprintf(":%s", configs.Env.AppPort), // ubah %d jadi %s
			ReadHeaderTimeout: time.Duration(configs.Env.AppReadHeaderTimeoutInSeconds) * time.Second,
		},
	}
}

func (s *Server) Serve(ctx context.Context) {
	s.srv.Handler = s.router.Register()
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("Failed to start HTTP server: %v", err))
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
