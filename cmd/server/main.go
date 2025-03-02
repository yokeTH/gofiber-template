package main

import (
	"context"

	"github.com/yokeTH/gofiber-template/internal/server"
)

func main() {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	s := server.New(server.WithPort(3000), server.WithName("gofiber-template"))
	s.Start(ctx, stop)
}
