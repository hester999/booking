package db

import (
	"fmt"
	"log"
	"time"

	"booking/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func InitDataBaseConnectionWithConfig(cfg config.DataBaseConfig) (*sqlx.DB, error) {
	log.Println("Connecting to database ...")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Address,
		cfg.Port,
		cfg.Database,
	)

	database, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		database, err = connectWithRetry(dsn, cfg.MaxAttempts, cfg.Delay)
		if err != nil {
			return nil, fmt.Errorf("could not connect to postgres database: %w", err)
		}
	}

	database.SetMaxOpenConns(cfg.MaxConnections)
	database.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleLiveTime) * time.Minute)
	database.SetMaxIdleConns(5)
	database.SetConnMaxLifetime(30 * time.Minute)
	return database, nil
}

func connectWithRetry(dsn string, maxAttempts, delay int) (*sqlx.DB, error) {
	var err error
	for i := 1; i <= maxAttempts; i++ {
		database, connectErr := sqlx.Connect("pgx", dsn)
		if connectErr == nil {
			log.Println("Database connected successfully.")
			return database, nil
		}

		err = connectErr
		log.Printf("Database not ready (attempt %d/%d): %v", i, maxAttempts, err)
		time.Sleep(time.Second * time.Duration(delay))
	}

	return nil, fmt.Errorf("could not connect to postgres database after %d attempts: %w", maxAttempts, err)
}
