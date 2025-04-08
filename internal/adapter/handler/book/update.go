package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/internal/adapter/presenter"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

// UpdateBook godoc
// @summary UpdateBook
// @description update book data
// @tags book
// @produce json
// @Param id path int true "Book ID"
// @param Book body presenter.UpdateBookRequest true "Book Data"
// @response 200 {object} presenter.SuccessResponse[presenter.BookResponse] "OK"
// @response 400 {object} presenter.ErrorResponse "Bad Request"
// @response 500 {object} presenter.ErrorResponse "Internal Server Error"
// @Router /books/{id} [patch]
func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.BadRequestError(err, "id must be an integer")
	}

	body := new(presenter.UpdateBookRequest)
	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	book, err := h.bookUseCase.Update(id, body)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "update book service failed")
	}

	res := h.presenter.ToResponse(book)
	return c.JSON(presenter.Success(res))
}
