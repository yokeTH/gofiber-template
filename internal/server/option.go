package server

type ServerOption func(*Server)

// WithName sets a custom name for the server
func WithName(name string) ServerOption {
	return func(s *Server) {
		s.config.Name = name
	}
}

// WithPort sets a custom port for the server
func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.config.Port = port
	}
}

// WithCorsAllowOrigins sets the allowed CORS origins
func WithCorsAllowOrigins(origins string) ServerOption {
	return func(s *Server) {
		s.config.CorsAllowOrigins = origins
	}
}

// WithCorsAllowMethods sets the allowed CORS methods
func WithCorsAllowMethods(methods string) ServerOption {
	return func(s *Server) {
		s.config.CorsAllowMethods = methods
	}
}

// WithCorsAllowHeaders sets the allowed CORS headers
func WithCorsAllowHeaders(headers string) ServerOption {
	return func(s *Server) {
		s.config.CorsAllowHeaders = headers
	}
}

// WithCorsAllowCredentials sets whether CORS should allow credentials
func WithCorsAllowCredentials(allowed bool) ServerOption {
	return func(s *Server) {
		s.config.CorsAllowCredentials = allowed
	}
}

// WithBodyLimitMB sets the body limit in MB
func WithBodyLimitMB(limit int) ServerOption {
	return func(s *Server) {
		s.config.BodyLimitMB = limit
	}
}

// WithConfig is a functional option that sets a custom configuration for the server.
// It takes a pointer to a Config struct and applies it to the server's configuration.
//
// Example usage:
//
//	server.WithConfig(&server.Config{Port: 8081, Name: "custom-app"})
func WithConfig(config *Config) ServerOption {
	return func(s *Server) {
		s.config = config
	}
}
