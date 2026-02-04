package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"billsplitter-monolith/internal/cfg"
	billimpl "billsplitter-monolith/internal/domain/bill/impl"
	eventimpl "billsplitter-monolith/internal/domain/event/impl"
	paymentmethodimpl "billsplitter-monolith/internal/domain/payment_method/impl"
	sessionimpl "billsplitter-monolith/internal/domain/session/impl"
	userimpl "billsplitter-monolith/internal/domain/user/impl"
	billrepo "billsplitter-monolith/internal/infrastructure/postgres/bill"
	eventrepo "billsplitter-monolith/internal/infrastructure/postgres/event"
	paymentmethodrepo "billsplitter-monolith/internal/infrastructure/postgres/payment_method"
	sessionrepo "billsplitter-monolith/internal/infrastructure/postgres/session"
	userrepo "billsplitter-monolith/internal/infrastructure/postgres/user"
	"billsplitter-monolith/internal/transport/http"
	authhttp "billsplitter-monolith/internal/transport/http/auth"
	billhttp "billsplitter-monolith/internal/transport/http/bill"
	meethttp "billsplitter-monolith/internal/transport/http/meet"
	"billsplitter-monolith/internal/transport/http/middleware"
	userhttp "billsplitter-monolith/internal/transport/http/user"
	billuc "billsplitter-monolith/internal/usecase/bill"
	eventuc "billsplitter-monolith/internal/usecase/event"
	useruc "billsplitter-monolith/internal/usecase/user"
	"billsplitter-monolith/internal/utils"
)

// @title           BillSplitter API
// @version         1.0
// @description     API для Telegram Mini App по разделению счетов
// @host            localhost:5001
// @BasePath        /
// @schemes         http
// @securityDefinitions.apikey SessionAuth
// @in header
// @name X-Session-ID
// @description     Сессионный токен. Значение из /auth/login/telegram передавайте через Authorize -> X-Session-ID header.
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
	userRepo := userrepo.NewRepository(db)
	paymentMethodRepo := paymentmethodrepo.New(db)
	sessionRepo := sessionrepo.NewRepository(db)
	eventRepo := eventrepo.NewRepository(db)
	billRepo := billrepo.NewRepository(db)

	// init service
	sessionSvc := sessionimpl.New(sessionRepo, l)
	userSvc := userimpl.New(userRepo)
	paymentMethodSvc := paymentmethodimpl.New(paymentMethodRepo)
	eventSvc := eventimpl.New(eventRepo)
	billSvc := billimpl.New(billRepo)

	// init use case
	userUC := useruc.New(userSvc, sessionSvc)
	billsUC := billuc.New(billSvc, eventSvc)
	meetUC := eventuc.New(eventSvc, userSvc, billSvc)

	// init http server
	mw := middleware.NewMiddlewareManager(sessionSvc, l)
	authCtrl := authhttp.NewController(userUC, userSvc, sessionSvc, l)
	userCtrl := userhttp.NewController(userSvc, paymentMethodSvc, l)
	meetCtrl := meethttp.NewController(meetUC, l)
	billCtrl := billhttp.NewController(billsUC, l)
	httpServer := http.NewServer(mw, authCtrl, userCtrl, meetCtrl, billCtrl, l)

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
	c, err := cfg.LoadConfig("config", "yaml")
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

	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		utils.LogFatalf(l, "failed to ping DB: %f", err)
	}

	l.WithGroup("main").Info("postgres connection loaded successfully")

	return db
}
