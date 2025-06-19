package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"billsplitter-monolith/internal/transport/http/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"billsplitter-monolith/internal/cfg"
	authsvc "billsplitter-monolith/internal/domain/auth/impl"
	sessionstorage "billsplitter-monolith/internal/repository/storage/session"
	userstorage "billsplitter-monolith/internal/repository/storage/user"
	"billsplitter-monolith/internal/transport/http"
	authhttp "billsplitter-monolith/internal/transport/http/auth"
	"billsplitter-monolith/internal/utils"
)

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// init infra
	appCfg := mustLoadCfg(l)
	db := mustInitGormDB(l, appCfg.Storage.Postgres)

	// init storages
	userStorage := userstorage.NewStorage(db)
	sessionCache := sessionstorage.NewMemCache()
	sessionStorage := sessionstorage.NewStorage(db, sessionCache)

	// init service
	authSvc := authsvc.New(userStorage, sessionStorage, l)

	// init http server
	mw := middleware.NewMiddlewareManager(authSvc, l)
	authCtrl := authhttp.NewController(authSvc, l)
	httpServer := http.NewServer(mw, authCtrl, l)

	go func() {
		err := httpServer.Start(ctx, appCfg.Server.Http)
		if err != nil {
			l.WithGroup("main").ErrorContext(ctx, err.Error())
			quit <- os.Interrupt
		}
	}()

	// Graceful stop
	l.WithGroup("main").InfoContext(ctx, "waiting app to stop...")
	<-quit
	cancel()
	l.WithGroup("main").InfoContext(ctx, "cancel signal has been received, stopping app...")

	err := httpServer.Stop(ctx)
	if err != nil {
		l.WithGroup("main").ErrorContext(ctx, fmt.Sprintf("failed to stop http server: %s", err.Error()))
	}
}

func mustLoadCfg(l *slog.Logger) cfg.Config {
	c, err := cfg.LoadConfig()
	if err != nil {
		utils.LogFatalf(l, "failed to load config: %v", err)
	}

	l.WithGroup("main").Info("config loaded successfully")

	return c
}

func mustInitGormDB(l *slog.Logger, cfg cfg.Postgres) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	// Проверка соединения
	sqlDB, err := db.DB()
	if err != nil {
		utils.LogFatalf(l, "failed to get raw DB: %f", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		utils.LogFatalf(l, "failed to ping DB: %f", err)
	}

	l.WithGroup("main").Info("postgres connection loaded successfully")

	return db
}
