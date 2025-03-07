package apperror

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/yokeTH/gofiber-template/pkg/dto"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	log.Warnf("your error is nil. passed message: %s", e.Message)
	return e.Message
}

func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
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

// ErrorHandler sends a JSON error response to the client based on the error type.
// If the error is an AppError, it uses the associated HTTP status code and message.
// For other errors, it either uses the status code from a fiber.Error or defaults to 500.
// If sending the JSON response fails, it responds with a plain text "Internal Server Error".
func ErrorHandler(c *fiber.Ctx, err error) error {

	// if is app error
	if IsAppError(err) {
		e := err.(*AppError)
		if err := c.Status(e.Code).JSON(dto.ErrorResponse{Error: e.Message}); err != nil {
			// if can't send error
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		return nil
	}

	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	if err := c.Status(code).JSON(dto.ErrorResponse{Error: err.Error()}); err != nil {
		// if can't send error
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return nil
}
