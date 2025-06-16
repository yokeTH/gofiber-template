package server

import (
	"context"
	"fmt"
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/swaggo/swag"
	"github.com/yokeTH/gofiber-scalar/scalar"
	"github.com/yokeTH/gofiber-template/docs"
	"github.com/yokeTH/gofiber-template/internal/adaptor/middleware"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

type Config struct {
	Env                  string `env:"ENV"`
	Name                 string `env:"NAME"`
	Port                 int    `env:"PORT"`
	BodyLimitMB          int    `env:"BODY_LIMIT_MB"`
	CorsAllowOrigins     string `env:"CORS_ALLOW_ORIGINS"`
	CorsAllowMethods     string `env:"CORS_ALLOW_METHODS"`
	CorsAllowHeaders     string `env:"CORS_ALLOW_HEADERS"`
	CorsAllowCredentials bool   `env:"CORS_ALLOW_CREDENTIALS"`
	SwaggerUser          string `env:"SWAGGER_USER"`
	SwaggerPass          string `env:"SWAGGER_PASS"`
}

const defaultEnv = "unknown"
const defaultName = "app"
const defaultPort = 8080
const defaultBodyLimitMB = 4
const defaultCorsAllowOrigins = "*"
const defaultCorsAllowMethods = "GET,POST,PUT,DELETE,PATCH,OPTIONS"
const defaultCorsAllowHeaders = "Origin,Content-Type,Accept,Authorization"
const defaultCorsAllowCredentials = false
const defalutSwaggerUser = "admin"
const defalutSwaggerPass = "1234"

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
		Env:                  defaultEnv,
		Name:                 defaultName,
		Port:                 defaultPort,
		BodyLimitMB:          defaultBodyLimitMB,
		CorsAllowOrigins:     defaultCorsAllowOrigins,
		CorsAllowMethods:     defaultCorsAllowMethods,
		CorsAllowHeaders:     defaultCorsAllowHeaders,
		CorsAllowCredentials: defaultCorsAllowCredentials,
		SwaggerUser:          defalutSwaggerUser,
		SwaggerPass:          defalutSwaggerPass,
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

	app.Use(requestid.New())
	app.Use(middleware.CtxRequestIDInjector())
	app.Use(middleware.RequestLogger())

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

	if server.config.Env == "dev" {
		swag.Register(docs.SwaggerInfo.InstanceName(), docs.SwaggerInfo)
		app.Get("/docs/*", basicauth.New(basicauth.Config{
			Users: map[string]string{
				server.config.SwaggerUser: server.config.SwaggerPass,
			},
		}), scalar.New())
	}

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
