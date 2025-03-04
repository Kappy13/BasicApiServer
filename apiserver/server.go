package apiserver

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/cap79/BasicApiServer/config"
)

type ApiServer struct {
	cfg    *config.Config
	logger *slog.Logger
	mux    *http.ServeMux
}

func New(cfg *config.Config, logger *slog.Logger) *ApiServer {
	return &ApiServer{
		cfg:    cfg,
		logger: logger,
		mux:    http.NewServeMux(),
	}
}

func (s *ApiServer) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (s *ApiServer) Start(ctx context.Context) error {
	s.mux.HandleFunc("/ping", s.ping)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(s.cfg.ApiServerHost, s.cfg.ApiServerPort),
		Handler: s.mux,
	}

	go func() {
		s.logger.Info("started apiserver on port:", "port", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("apiserver failed to listen and server", "error", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("apiserver failed to shutdown", "error", err)
		}
		s.logger.Info("apiserver shutting down", "shutdown", "success")
	}()
	wg.Wait()
	return nil
}
