package main

import (
	"context"
	"log"

	"github.com/yokeTH/gofiber-template/internal/core/service"
	"github.com/yokeTH/gofiber-template/internal/database"
	"github.com/yokeTH/gofiber-template/internal/handler"
	"github.com/yokeTH/gofiber-template/internal/repository"
	"github.com/yokeTH/gofiber-template/internal/server"
	"github.com/yokeTH/gofiber-template/pkg/config"
)

func main() {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	config := config.Load()

	db, err := database.NewPostgresDB(config.PSQL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	bookRepository := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	s := server.New(
		server.WithName(config.Server.Name),
		server.WithBodyLimitMB(config.Server.BodyLimitMB),
		server.WithPort(config.Server.Port),
	)

	s.Get("/books", bookHandler.GetBooks)
	s.Get("/books/:id", bookHandler.GetBook)
	s.Post("/books", bookHandler.CreateBook)
	s.Patch("/books/:id", bookHandler.UpdateBook)
	s.Delete("/books/:id", bookHandler.DeleteBook)

	s.Start(ctx, stop)
}
