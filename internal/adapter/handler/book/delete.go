package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

// DeleteBook godoc
// @summary DeleteBook
// @description delete book by id
// @tags book
// @produce json
// @Param id path int true "Book ID"
// @response 204 "No Content"
// @response 400 {object} presenter.ErrorResponse "Bad Request"
// @response 500 {object} presenter.ErrorResponse "Internal Server Error"
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.BadRequestError(err, "id must be an integer")
	}

	if err := h.bookUseCase.Delete(id); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "delete book service failed")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
