package middleware

import (
	"bytes"
	"errors"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/logger"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.Req(c.UserContext(), map[string]any{
			"method": c.Method(),
			"path":   c.Path(),
			"query":  c.Queries(),
			"body":   parseLogBody(c.Body()),
		})

		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		logEntry := map[string]any{
			"method":     c.Method(),
			"path":       c.Path(),
			"status":     c.Response().StatusCode(),
			"process_ms": duration.Milliseconds(),
		}

		if err != nil {
			// logEntry["error"] = err
			var fiberErr *fiber.Error
			var appErr *apperror.AppError
			var code int
			if errors.As(err, &appErr) {
				code = appErr.Code
			} else if errors.As(err, &fiberErr) {
				code = fiberErr.Code
			} else {
				code = 500
			}
			logEntry["status"] = code
			logger.Error(c.UserContext(), "failed to process request", "data", logEntry, "err", err)
		} else {
			logEntry["body"] = parseLogBody(c.Response().Body())
			logger.Res(c.UserContext(), logEntry)
		}

		return err
	}
}

func parseLogBody(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	if json.Valid(data) {
		return string(data)
	}

	if utf8String := bytes.Runes(data); len(utf8String) > 0 && len(data) < 256 {
		return string(bytes.ReplaceAll(bytes.TrimSpace(data), []byte(" "), []byte("")))
	}

	return "body Neither JSON nor String (len < 256)"
}
