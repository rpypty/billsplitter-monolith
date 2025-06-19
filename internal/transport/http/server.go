package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"billsplitter-monolith/internal/cfg"
	"billsplitter-monolith/internal/transport/http/auth"
	mw "billsplitter-monolith/internal/transport/http/middleware"
)

type Server struct {
	authCtrl auth.Controller
	mw       mw.Manager

	httpSrv *http.Server
	logger  *slog.Logger
}

func NewServer(
	mw mw.Manager,
	authCtrl auth.Controller,
	logger *slog.Logger,
) *Server {
	return &Server{
		authCtrl: authCtrl,
		logger:   logger,
		mw:       mw,
	}
}

func (s *Server) Start(_ context.Context, cfg cfg.Http) error {
	l := s.l()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 20))

	// init routes
	auth.InitRoutes(r, s.authCtrl, s.mw)

	s.httpSrv = &http.Server{
		Addr:              cfg.Port,
		Handler:           r,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	l.Info("starting http server on port: " + cfg.Port)

	if err := s.httpSrv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.l().Info("shutting down http server")
	return s.httpSrv.Shutdown(ctx)
}

func (s *Server) l() *slog.Logger {
	return s.logger.WithGroup("auth-http-server")
}
