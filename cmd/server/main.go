package main

import (
	"booking/internal/infra/repository/cache"
	router "booking/internal/router/booking"
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	bookinguc "booking/internal/application/usecases/booking"
	"booking/internal/config"
	bookinghandler "booking/internal/handlers/booking"
	expiredchecker "booking/internal/infra/background/expired"
	"booking/internal/infra/db"
	"booking/internal/infra/redis"
	bookingrepo "booking/internal/infra/repository/booking"
	pgrepo "booking/internal/infra/repository/postgres/booking"
	eventrepo "booking/internal/infra/repository/postgres/event"
	"booking/internal/infra/txmanager"
	"booking/pkg"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func main() {
	logger := pkg.NewLogger()

	cfg := config.NewConfig(logger)
	cfg.SetCFGFromEnvFile()
	cfg.LoadConfig()

	database, err := db.InitDataBaseConnectionWithConfig(cfg.DataBaseConfig)
	if err != nil {
		log.Fatal(err)
	}
	redisClient, err := redis.NewRedisClient(cfg.CacheConfig)
	if err != nil {
		log.Fatal(err)
	}

	txFactory := trmsqlx.NewDefaultFactory(database)
	txManager := manager.Must(txFactory)
	txGetter := trmsqlx.DefaultCtxGetter
	txRunner := txmanager.NewTRMRunner(txManager)

	pgRepo := pgrepo.NewBookingRepo(database, txGetter)
	eventRepo := eventrepo.NewEventRepo(database, txGetter)

	cahceRepo := cache.NewCacheRepo(redisClient.Client)

	bookingRepo := bookingrepo.NewCachedBookingRepo(pgRepo, cahceRepo, logger)

	bookingUsecase := bookinguc.NewBookingUsecase(bookingRepo, eventRepo, txRunner)
	bookingHandler := bookinghandler.NewBookingHandler(bookingUsecase, logger)
	checker := expiredchecker.NewExpiredChecker(bookingUsecase, cfg.Ticker, txRunner, logger)

	bookingRouter := router.NewRouter(bookingHandler)

	if err := runServer(bookingRouter, cfg, database, checker, logger, redisClient); err != nil {
		log.Fatal(err)
	}
}

func runServer(
	router *mux.Router,
	cfg *config.Config,
	database *sqlx.DB,
	checker *expiredchecker.ExpiredChecker,
	logger pkg.Logger,
	redisClient *redis.RedisClient,
) error {
	defer func() {
		if err := database.Close(); err != nil {
			logger.Error("database close failed", "error", err)
		}
		if err := redisClient.Client.Close(); err != nil {
			logger.Error("database close failed", "error", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go logBackgroundErrors(ctx, logger, "expired_checker", checker.GetErrChan())
	go runBackground(ctx, stop, logger, "expired_checker", checker.Start)

	server := &http.Server{
		Addr:    serverAddress(cfg.ServerConfig),
		Handler: router,
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("http server started", "addr", server.Addr)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			stop()
			return
		}
		serverErr <- nil
	}()

	var err error
	select {
	case <-ctx.Done():
	case err = <-serverErr:
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
		logger.Error("http server shutdown failed", "error", shutdownErr)
		if err == nil {
			err = shutdownErr
		}
	}

	if err != nil {
		return err
	}

	logger.Info("service stopped")
	return nil
}

func runBackground(
	ctx context.Context,
	stop context.CancelFunc,
	logger pkg.Logger,
	name string,
	start func(context.Context) error,
) {
	if err := start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("background worker stopped", "worker", name, "error", err)
		stop()
	}
}

func logBackgroundErrors(ctx context.Context, logger pkg.Logger, source string, errChan <-chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-errChan:
			if !ok {
				return
			}
			if err != nil {
				logger.Error("background worker error", "worker", source, "error", err)
			}
		}
	}
}

func serverAddress(cfg config.ServerConfig) string {
	return net.JoinHostPort(cfg.Address, cfg.Port)
}
