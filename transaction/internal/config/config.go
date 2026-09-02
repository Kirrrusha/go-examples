package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName     string `env:"SERVICE_NAME" required:"true" envDefault:"transaction"`
	AppEnv          string `env:"APP_ENV" required:"true" envDefault:"development"`
	Host            string `env:"GRPC_HOST" required:"true" envDefault:"localhost"`
	Port            int    `env:"GRPC_PORT" required:"true" envDefault:"50054"`
	LogLevel        string `env:"LOG_LEVEL" required:"true" envDefault:"info"`
	DbDsn           string `env:"DB_DSN" required:"true"`
	AccountGrpcHost string `env:"ACCOUNT_GRPC_HOST" required:"true"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
