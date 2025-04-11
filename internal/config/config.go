package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/yokeTH/gofiber-template/internal/server"
	"github.com/yokeTH/gofiber-template/pkg/db"
	"github.com/yokeTH/gofiber-template/pkg/storage"
)

type config struct {
	Server        server.Config  `envPrefix:"SERVER_" env:"required"`
	PSQL          db.DBConfig    `envPrefix:"POSTGRES_" env:"required"`
	PublicBucket  storage.Config `envPrefix:"PUBLIC_" env:"required"`
	PrivateBucket storage.Config `envPrefix:"PRIVATE_" env:"required"`
}

func Load() *config {
	config := &config{}

	if err := godotenv.Load(); err != nil {
		log.Warnf("Unable to load .env file: %s", err)
	}

	if err := env.Parse(config); err != nil {
		log.Fatalf("Unable to parse env vars: %s", err)
	}

	return config
}
