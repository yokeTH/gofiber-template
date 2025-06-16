package logger

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type logger struct {
	zerolog zerolog.Logger
}

var (
	appLogger *logger
	appOnce   sync.Once
	logLevel  string
	logJSON   string
)

func init() {
	appOnce.Do(func() {
		_ = godotenv.Load()
		logLevel = os.Getenv("LOG_LEVEL")
		logJSON = os.Getenv("LOG_JSON")
		appLogger = &logger{
			zerolog: createLogger(),
		}
	})
}

func createLogger() zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	zerolog.LevelFieldMarshalFunc = func(level zerolog.Level) string {
		return strings.ToUpper(level.String())
	}
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z"
	zerolog.TimestampFieldName = "timestamp"

	var logger zerolog.Logger

	if parseBool(logJSON) {
		logger = zerolog.New(os.Stdout).With().
			Timestamp().
			Logger().
			Level(parseLogLevel(logLevel))
	} else {
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out: os.Stdout,
		}).
			Level(zerolog.TraceLevel).
			With().
			Timestamp().
			Logger()
	}

	return logger
}

func parseLogLevel(raw string) zerolog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "trace", "-1":
		return zerolog.TraceLevel
	case "debug", "0":
		return zerolog.DebugLevel
	case "info", "1":
		return zerolog.InfoLevel
	case "warn", "warning", "2":
		return zerolog.WarnLevel
	case "error", "3":
		return zerolog.ErrorLevel
	case "fatal", "4":
		return zerolog.FatalLevel
	case "panic", "5":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

func parseBool(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}

func logWithLevel(ctx context.Context, level zerolog.Level, msg string, fields ...any) {
	requestID, _ := ctx.Value("requestid").(string)
	event := appLogger.zerolog.WithLevel(level).
		Str("request_id", requestID)

	var errVal error
	if len(fields)%2 == 0 {
		for i := 0; i < len(fields); i += 2 {
			key, ok := fields[i].(string)
			if !ok {
				continue
			}

			if key == "error" || key == "err" {
				if e, ok := fields[i+1].(error); ok && errVal == nil {
					errVal = e
					continue
				}
			}
			event = event.Interface(key, fields[i+1])
		}
	}

	if level >= zerolog.ErrorLevel && errVal != nil {
		event = event.Err(errVal).Stack()
	}
	event.Msg(msg)
}

// Use for log enter and exit function
func Trace(ctx context.Context, msg string, fields ...any) {
	logWithLevel(ctx, zerolog.TraceLevel, msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...any) {
	logWithLevel(ctx, zerolog.DebugLevel, msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...any) {
	logWithLevel(ctx, zerolog.InfoLevel, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...any) {
	logWithLevel(ctx, zerolog.WarnLevel, msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...any) {
	logWithLevel(ctx, zerolog.ErrorLevel, msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...any) {
	logWithLevel(ctx, zerolog.FatalLevel, msg, fields...)
}

func Panic(ctx context.Context, msg string, fields ...any) {
	logWithLevel(ctx, zerolog.PanicLevel, msg, fields...)
}

// Use for log enter and exit function
func Func(ctx context.Context, name string, exiting ...bool) {
	state := "enter"
	if len(exiting) == 1 && exiting[0] {
		state = "exit"
	}
	message := fmt.Sprintf("%s %s", state, name)
	Trace(ctx, message)
}

func Req(ctx context.Context, req any) {
	Info(ctx, "incoming request", "req", req)
}

func Res(ctx context.Context, res any) {
	Info(ctx, "processed request", "res", res)
}
