package server

import (
	"context"
	"fmt"
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

type Config struct {
	Name                 string `env:"NAME"`
	Port                 int    `env:"PORT"`
	BodyLimitMB          int    `env:"BODY_LIMIT_MB"`
	CorsAllowOrigins     string `env:"CORS_ALLOW_ORIGINS"`
	CorsAllowMethods     string `env:"CORS_ALLOW_METHODS"`
	CorsAllowHeaders     string `env:"CORS_ALLOW_HEADERS"`
	CorsAllowCredentials bool   `env:"CORS_ALLOW_CREDENTIALS"`
}

const defaultName = "app"
const defaultPort = 8080
const defaultBodyLimitMB = 4
const defaultCorsAllowOrigins = "*"
const defaultCorsAllowMethods = "GET,POST,PUT,DELETE,PATCH,OPTIONS"
const defaultCorsAllowHeaders = "Origin,Content-Type,Accept,Authorization"
const defaultCorsAllowCredentials = false

type Server struct {
	app    *fiber.App
	config *Config
}

// New creates a new instance of the Server with the default configuration values.
// Additional configuration can be applied using functional options passed in the opts parameter.
//
// Default values are:
//   - Name: "app"
//   - Port: 8080
//   - BodyLimitMB: 4MB
//   - CORS settings: default allows all origins, methods, headers, and credentials
//
// opts: A variadic list of functional options to customize the server's configuration.
//
// Example usage:
//
//	server := server.New(
//	  server.WithPort(3000),
//	  server.WithCorsAllowOrigins("https://example.com"),
//	)
func New(opts ...ServerOption) *Server {

	defaultConfig := &Config{
		Name:                 defaultName,
		Port:                 defaultPort,
		BodyLimitMB:          defaultBodyLimitMB,
		CorsAllowOrigins:     defaultCorsAllowOrigins,
		CorsAllowMethods:     defaultCorsAllowMethods,
		CorsAllowHeaders:     defaultCorsAllowHeaders,
		CorsAllowCredentials: defaultCorsAllowCredentials,
	}

	server := &Server{
		config: defaultConfig,
	}

	for _, opt := range opts {
		opt(server)
	}

	app := fiber.New(fiber.Config{
		AppName:               server.config.Name,
		BodyLimit:             server.config.BodyLimitMB * 1024 * 1024,
		CaseSensitive:         true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
		ErrorHandler:          apperror.ErrorHandler,
	})

	server.app = app

	return server
}

func (s *Server) Use(middleware ...fiber.Handler) {
	s.app.Use(middleware)
}

func (s *Server) Start(ctx context.Context, stop context.CancelFunc) {
	go func() {
		if err := s.app.Listen(fmt.Sprintf(":%d", s.config.Port)); err != nil {
			log.Fatalf("failed to start server: %v", err)
			stop()
		}
	}()

	defer func() {
		if err := s.app.Shutdown(); err != nil {
			log.Printf("failed to shutdown server: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("shutting down server...")
}
