package server

import (
	"context"
	"fmt"
	"log"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

type Config struct {
	Name                 string `env:"NAME,required"`
	Port                 int    `env:"PORT,required"`
	BodyLimitMB          int    `env:"BODY_LIMIT_MB,required"`
	CorsAllowOrigins     string `env:"CORS_ALLOW_ORIGINS,required"`
	CorsAllowMethods     string `env:"CORS_ALLOW_METHODS,required"`
	CorsAllowHeaders     string `env:"CORS_ALLOW_HEADERS,required"`
	CorsAllowCredentials bool   `env:"CORS_ALLOW_CREDENTIALS,required"`
}

const defaultName = "app"
const defaultPort = 8080
const defaultBodyLimitMB = 4
const defaultCorsAllowOrigins = "*"
const defaultCorsAllowMethods = "GET,POST,PUT,DELETE,PATCH,OPTIONS"
const defaultCorsAllowHeaders = "Origin,Content-Type,Accept,Authorization"
const defaultCorsAllowCredentials = false

type Server struct {
	*fiber.App
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
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
		DisableStartupMessage: true,
		ErrorHandler:          apperror.ErrorHandler,
	})

	app.Use(requestid.New())

	app.Use(logger.New(logger.Config{
		DisableColors: true,
		TimeFormat:    "2006-01-02 15:04:05",
		Format:        "${time} | ${locals:requestid} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
	}))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessEndpoint: "/health",
		LivenessProbe: func(c *fiber.Ctx) bool {
			if err := c.JSON(fiber.Map{"status": "ok"}); err != nil {
				return false
			}
			return true
		},
	}))

	server.App = app

	return server
}

func (s *Server) Start(ctx context.Context, stop context.CancelFunc) {
	go func() {
		if err := s.Listen(fmt.Sprintf(":%d", s.config.Port)); err != nil {
			log.Fatalf("failed to start server: %v", err)
			stop()
		}
	}()

	defer func() {
		if err := s.Shutdown(); err != nil {
			log.Printf("failed to shutdown server: %v.", err)
		}
	}()

	<-ctx.Done()

	log.Println("shutting down server...")
}
