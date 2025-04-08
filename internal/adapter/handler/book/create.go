package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/internal/adapter/presenter"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

// CreateBook godoc
// @summary CreateBook
// @description create book by title and author
// @tags book
// @accept json
// @produce json
// @param Book body presenter.CreateBookRequest true "Book Data"
// @response 201 {object} presenter.SuccessResponse[presenter.BookResponse] "Created"
// @response 400 {object} presenter.ErrorResponse "Bad Request"
// @response 500 {object} presenter.ErrorResponse "Internal Server Error"
// @Router /books [post]
func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	body := new(presenter.CreateBookRequest)
	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	book := &domain.Book{
		Author: body.Author,
		Title:  body.Title,
	}

	if err := h.bookUseCase.Create(book); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "create book service failed")
	}

	res := h.presenter.ToResponse(book)
	return c.Status(fiber.StatusCreated).JSON(presenter.Success(res))
}
