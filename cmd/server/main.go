package main

import (
	"context"

	"github.com/yokeTH/gofiber-template/internal/server"
	"github.com/yokeTH/gofiber-template/pkg/config"
)

func main() {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	config := config.Load()

	s := server.New(
		server.WithName(config.Server.Name),
		server.WithBodyLimitMB(config.Server.BodyLimitMB),
		server.WithPort(config.Server.Port),
	)

	s.Start(ctx, stop)
}
