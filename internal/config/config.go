package config

import (
	"strings"

	"booking/pkg"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Address string `env:"HTTP_ADDRESS" env-default:""`
	Port    string `env:"PORT" env-default:"8080"`
}

type DataBaseConfig struct {
	Address        string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port           string `env:"POSTGRES_PORT" env-default:"5432"`
	Database       string `env:"POSTGRES_DB" env-default:"booking"`
	User           string `env:"POSTGRES_USER" env-default:"admin"`
	Password       string `env:"POSTGRES_PASSWORD" env-default:"123"`
	MaxAttempts    int    `env:"MAX_ATTEMPTS" env-default:"1"`
	Delay          int    `env:"DELAY" env-default:"1"`
	MaxConnections int    `env:"MAX_CONNECTIONS" env-default:"1"`
}

type CacheConfig struct {
	Address  string `env:"CACHE_ADDRESS" env-default:"localhost"`
	Database string `env:"CACHE_DATABASE" env-default:"0"`
}

type BrokerConfig struct {
	Address        string `env:"KAFKA_BROKERS"`
	Host           string `env:"KAFKA_HOST" env-default:"localhost"`
	Port           string `env:"KAFKA_PORT" env-default:"9093"`
	Topic          string `env:"KAFKA_TOPIC" env-default:"booking-events"`
	MaxRetry       int    `env:"KAFKA_MAX_RETRY" env-default:"3"`
	Backoff        int64  `env:"KAFKA_BACKOFF_SECONDS" env-default:"1"`
	MaxConnections int    `env:"KAFKA_MAX_CONNECTIONS" env-default:"5"`
}

type InfraConfig struct {
	Ticker int64 `env:"BACKGROUND_TICKER_SECONDS" env-default:"10"`
}

type Config struct {
	ServerConfig
	DataBaseConfig
	CacheConfig
	BrokerConfig
	InfraConfig

	logger pkg.Logger `env:"-"`
}

func NewConfig(logger pkg.Logger) *Config {
	return &Config{logger: logger}
}

func (c *Config) SetCFGFromEnvFile() {
	paths := []string{".env", "../.env", "../../.env"}
	_ = godotenv.Load(paths...)
}

func (c *Config) LoadConfig() {
	if err := cleanenv.ReadEnv(c); err != nil {
		if c.logger != nil {
			c.logger.Error("config read error", "error", err)
		}
		panic(err)
	}

	c.normalize()
}

func (c *Config) normalize() {
	if c.ServerConfig.Port == "" {
		c.ServerConfig.Port = "8080"
	}
	if c.DataBaseConfig.MaxAttempts <= 0 {
		c.DataBaseConfig.MaxAttempts = 1
	}
	if c.DataBaseConfig.Delay <= 0 {
		c.DataBaseConfig.Delay = 1
	}
	if c.DataBaseConfig.MaxConnections <= 0 {
		c.DataBaseConfig.MaxConnections = 1
	}
	if c.BrokerConfig.MaxRetry < 0 {
		c.BrokerConfig.MaxRetry = 0
	}
	if c.BrokerConfig.Backoff <= 0 {
		c.BrokerConfig.Backoff = 1
	}
	if c.BrokerConfig.MaxConnections <= 0 {
		c.BrokerConfig.MaxConnections = 1
	}
	if c.BrokerConfig.Address == "" {
		c.BrokerConfig.Address = buildAddress(c.BrokerConfig.Host, c.BrokerConfig.Port)
	}
	if c.InfraConfig.Ticker <= 0 {
		c.InfraConfig.Ticker = 10
	}
}

func buildAddress(host, port string) string {
	if strings.Contains(host, ":") {
		return host
	}
	return host + ":" + port
}
