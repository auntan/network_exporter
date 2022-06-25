package http_server

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Config struct {
	Port int
}

type HTTPServer struct {
	server *http.Server
}

func New(conf *Config) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{Addr: ":" + strconv.Itoa(conf.Port)},
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		<-ctx.Done()
		err := s.server.Shutdown(context.TODO())
		if err != nil {
			zap.L().Warn(fmt.Sprintf("HTTP server shutdown error: %v", err))
		}
	}()

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}
