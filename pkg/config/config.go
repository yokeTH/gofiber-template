package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/yokeTH/gofiber-template/internal/database"
	"github.com/yokeTH/gofiber-template/internal/server"
	"github.com/yokeTH/gofiber-template/pkg/storage"
)

type Config struct {
	Server server.Config           `envPrefix:"SERVER_"`
	PSQL   database.PostgresConfig `envPrefix:"POSTGRES_"`
	R2     storage.R2Config        `envPrefix:"R2_"`
}

func Load() *Config {
	config := &Config{}

	if err := godotenv.Load(); err != nil {
		log.Warnf("Unable to load .env file: %s", err)
	}

	if err := env.Parse(config); err != nil {
		log.Fatalf("Unable to parse env vars: %s", err)
	}

	return config
}
