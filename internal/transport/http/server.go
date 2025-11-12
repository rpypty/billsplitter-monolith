package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"billsplitter-monolith/internal/cfg"
	"billsplitter-monolith/internal/transport/http/auth"
	"billsplitter-monolith/internal/transport/http/bill"
	"billsplitter-monolith/internal/transport/http/meet"
	mw "billsplitter-monolith/internal/transport/http/middleware"
	"billsplitter-monolith/internal/transport/http/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	httpswagger "github.com/swaggo/http-swagger"

	_ "billsplitter-monolith/docs" // важно: импорт с подчёркиванием!
)

type Server struct {
	authCtrl auth.Controller
	userCtrl user.Controller
	meetCtrl meet.Controller
	billCtrl bill.Controller
	mw       mw.Manager

	httpSrv *http.Server
	logger  *slog.Logger
}

func NewServer(
	mw mw.Manager,
	authCtrl auth.Controller,
	userCtrl user.Controller,
	meetCtrl meet.Controller,
	billCtrl bill.Controller,
	logger *slog.Logger,
) *Server {
	return &Server{
		authCtrl: authCtrl,
		userCtrl: userCtrl,
		meetCtrl: meetCtrl,
		billCtrl: billCtrl,
		logger:   logger,
		mw:       mw,
	}
}

func (s *Server) Start(_ context.Context, cfg cfg.Http) error {
	l := s.l()

	r := chi.NewRouter()

	// Настройка CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешить все origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Session-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Максимальное время кеширования preflight запроса
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 20))

	r.Get("/swagger/*", httpswagger.WrapHandler)

	// init routes
	auth.InitRoutes(r, s.authCtrl, s.mw)
	user.InitRoutes(r, s.userCtrl, s.mw)
	meet.InitRoutes(r, s.meetCtrl, s.mw)
	bill.InitRoutes(r, s.billCtrl, s.mw)

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
