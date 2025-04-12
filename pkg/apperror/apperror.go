package apperror

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	Code    int
	Message string
	Err     error
	Stack   string
}

func (e *AppError) Error() string {
	if e.Code/100 == 5 {
		return fmt.Sprintf("%s\nInternal Error: %v \nStack:\n%s", e.Message, e.Err, e.Stack)
	}
	return fmt.Sprintf("%s - Internal Error: %v", e.Message, e.Err)
}

func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

func New(code int, message string, err error) *AppError {
	stack := captureStack()
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Stack:   stack,
	}
}

// modified https://github.com/pkg/errors/blob/5dd12d0cfe7f152f80558d591504ce685299311e/stack.go
func captureStack() string {
	const depth = 32
	var pcs [depth]uintptr

	// skip 4 frames apperror.captureStack x2, apperror.New, apperror.InternalServerError or another
	n := runtime.Callers(4, pcs[:])

	if n == 0 {
		return "stack trace is not available"
	}

	stackTrace := make([]string, 0, n)
	for i := range n {
		fn := runtime.FuncForPC(pcs[i])
		file, line := fn.FileLine(pcs[i])
		stackTrace = append(stackTrace, fmt.Sprintf("%s\n\tat %s:%d", fn.Name(), file, line))
	}

	return strings.Join(stackTrace, "\n")
}

func InternalServerError(err error, msg string) *AppError {
	return New(fiber.StatusInternalServerError, msg, err)
}

func BadRequestError(err error, msg string) *AppError {
	return New(fiber.StatusBadRequest, msg, err)
}

func UnauthorizedError(err error, msg string) *AppError {
	return New(fiber.StatusUnauthorized, msg, err)
}

func ForbiddenError(err error, msg string) *AppError {
	return New(fiber.StatusForbidden, msg, err)
}

func NotFoundError(err error, msg string) *AppError {
	return New(fiber.StatusNotFound, msg, err)
}

func ConflictError(err error, msg string) *AppError {
	return New(fiber.StatusConflict, msg, err)
}

func UnprocessableEntityError(err error, msg string) *AppError {
	return New(fiber.StatusUnprocessableEntity, msg, err)
}

func ErrorHandler(c *fiber.Ctx, err error) error {

	// if is app error
	if IsAppError(err) {
		e := err.(*AppError)
		if err := c.Status(e.Code).JSON(fiber.Map{"error": e.Message}); err != nil {
			// if can't send error -- it should not be able
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		return nil
	}

	var e *fiber.Error
	if errors.As(err, &e) {
		if err := c.Status(e.Code).SendString(e.Error()); err != nil {
			// if can't send error -- it should not be able
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		return nil
	}

	// other case return error that is not fiber error or app error
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
}
