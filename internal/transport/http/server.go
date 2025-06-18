package http

import (
	"context"
	"log/slog"
	"net/http"

	"billsplitter-monolith/internal/cfg"
	"billsplitter-monolith/internal/transport/http/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	logger   *slog.Logger
	authCtrl *auth.Controller

	httpSrv *http.Server
}

func NewServer(authCtrl *auth.Controller, logger *slog.Logger) *Server {
	return &Server{
		logger:   logger,
		authCtrl: authCtrl,
	}
}

func (s *Server) Start(ctx context.Context, cfg cfg.Http) error {
	l := s.l()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60))

	// init routes
	auth.InitRoutes(r, s.authCtrl)

	s.httpSrv = &http.Server{
		Addr:    cfg.Port,
		Handler: r,
	}

	l.InfoContext(ctx, "starting http server on port: "+cfg.Port)
	if err := s.httpSrv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}

func (s *Server) l() *slog.Logger {
	return s.logger.WithGroup("auth-http-server")
}
