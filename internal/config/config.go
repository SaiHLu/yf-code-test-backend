package config

import "github.com/caarlos0/env/v11"

type AppConfig struct {
	GO_ENV                 string `env:"GO_ENV" envDefault:"development"`
	PORT                   string `env:"PORT" envDefault:"8080"`
	MODE                   string `env:"MODE" envDefault:"debug"`
	POSTGRES_USERNAME      string `env:"POSTGRES_USERNAME"`
	POSTGRES_PASSWORD      string `env:"POSTGRES_PASSWORD"`
	POSTGRES_HOST          string `env:"POSTGRES_HOST"`
	POSTGRES_PORT          string `env:"POSTGRES_PORT" envDefault:"5432"`
	POSTGRES_DB            string `env:"POSTGRES_DB"`
	POSTGRES_SSLMODE       string `env:"POSTGRES_SSLMODE" envDefault:"disable"`
	MONGODB_URI            string `env:"MONGODB_URI" envDefault:"mongodb://localhost:27017"`
	REDIS_ADDRESS          string `env:"REDIS_ADDRESS" envDefault:"localhost:6379"`
	REDIS_PASSWORD         string `env:"REDIS_PASSWORD" envDefault:""`
	REDIS_DB               int    `env:"REDIS_DB" envDefault:"0"`
	REDIS_USER_LOG_CHANNEL string `env:"REDIS_USER_LOG_CHANNEL" envDefault:"user_log_channel"`
	ACCESS_TOKEN_KEY       string `env:"ACCESS_TOKEN_KEY" envDefault:"secret"`
	ACCESS_TOKEN_TTL       int    `env:"ACCESS_TOKEN_TTL" envDefault:"3600"`
	REFRESH_TOKEN_KEY      string `env:"REFRESH_TOKEN_KEY" envDefault:"refresh_secret"`
	REFRESH_TOKEN_TTL      int    `env:"REFRESH_TOKEN_TTL" envDefault:"86400"`
	CORS_ALLOWED_ORIGINS   string `env:"CORS_ALLOWED_ORIGINS"`
	CORS_ALLOWED_METHODS   string `env:"CORS_ALLOWED_METHODS" envDefault:"GET,POST,PUT,PATCH,DELETE,OPTIONS"`
	CORS_ALLOWED_HEADERS   string `env:"CORS_ALLOWED_HEADERS"`
	CORS_EXPOSED_HEADERS   string `env:"CORS_EXPOSED_HEADERS"`
	CORS_ALLOW_CREDENTIALS bool   `env:"CORS_ALLOW_CREDENTIALS" envDefault:"true"`
}

var config AppConfig

func NewAppConfig() *AppConfig {
	if err := env.Parse(&config); err != nil {
		panic(err)
	}

	return &config
}
