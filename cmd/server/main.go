package main

import (
	"context"
	"log"

	"github.com/gofiber/swagger"
	"github.com/swaggo/swag"
	"github.com/yokeTH/gofiber-template/docs"
	"github.com/yokeTH/gofiber-template/internal/core/service"
	"github.com/yokeTH/gofiber-template/internal/db"
	"github.com/yokeTH/gofiber-template/internal/handler"
	"github.com/yokeTH/gofiber-template/internal/repository"
	"github.com/yokeTH/gofiber-template/internal/server"
	"github.com/yokeTH/gofiber-template/pkg/config"
)

// @title GO-FIBER-TEMPLATE API
// @version 1.0
// @servers https http
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Bearer token authentication
func main() {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	config := config.Load()

	db, err := db.New(config.PSQL)
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

	swag.Register(docs.SwaggerInfo.InstanceName(), docs.SwaggerInfo)
	s.Get("/swagger/*", swagger.HandlerDefault)

	s.Get("/books", bookHandler.GetBooks)
	s.Get("/books/:id", bookHandler.GetBook)
	s.Post("/books", bookHandler.CreateBook)
	s.Patch("/books/:id", bookHandler.UpdateBook)
	s.Delete("/books/:id", bookHandler.DeleteBook)

	s.Start(ctx, stop)
}
