package main

import (
	"booking/internal/application/usecases/event"
	"booking/internal/config"
	"booking/internal/infra/background/producer"
	"booking/internal/infra/db"
	"booking/internal/infra/kafka"
	eventrepo "booking/internal/infra/repository/postgres/event"
	"booking/internal/infra/txmanager"
	"booking/pkg"
	"context"
	"errors"
	"log"
	"os/signal"
	"syscall"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
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

	txFactory := trmsqlx.NewDefaultFactory(database)
	txManager := manager.Must(txFactory)
	txGetter := trmsqlx.DefaultCtxGetter
	txRunner := txmanager.NewTRMRunner(txManager)

	eventRepo := eventrepo.NewEventRepo(database, txGetter)

	eventUsecases := event.NewEventUsecases(txRunner, eventRepo)

	kafkaProd, err := kafka.NewKafkaProducer(cfg.BrokerConfig, logger)
	if err != nil {
		log.Fatal(err)
	}
	outboxProducer := producer.NewProducer(kafkaProd, eventUsecases, logger, cfg.InfraConfig)

	defer func() {
		if err := kafkaProd.Close(); err != nil {
			logger.Error("kafka producer close failed", "error", err)
		}
		if err := database.Close(); err != nil {
			logger.Error("database close failed", "error", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go runBackground(ctx, stop, logger, "producer", outboxProducer.Start)
	go logBackgroundErrors(ctx, logger, "producer", outboxProducer.GetErrChan())

	select {
	case <-ctx.Done():
	}

	logger.Info("worker stopped")
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
