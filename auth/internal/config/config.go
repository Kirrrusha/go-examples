package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName           string `env:"SERVICE_NAME" json:"service_name" required:"true" envDefault:"auth-service"`
	AppEnv                string `env:"APP_ENV" json:"app_environment" required:"true" envDefault:"development"`
	Host                  string `env:"GRPC_HOST" json:"host" required:"true" envDefault:"localhost"`
	Port                  int    `env:"GRPC_PORT" json:"port" required:"true" envDefault:"50052"`
	LogLevel              string `env:"LOG_LEVEL" json:"log_level" required:"true" envDefault:"info"`
	DbDsn                 string `env:"DB_DSN" json:"db_dsn" required:"true"`
	JwtSecret             string `env:"JWT_SECRET" json:"jwt_secret" required:"true"`
	AccessTokenTTLMinutes int    `env:"ACCESS_TOKEN_TTL_MINUTES" json:"access_ttl_min" required:"true" envDefault:"60"`
	RefreshTokenTTLDays   int    `env:"REFRESH_TOKEN_TTL_DAYS" json:"refresh_ttl_days" required:"true" envDefault:"30"`
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
