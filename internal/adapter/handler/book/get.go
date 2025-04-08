package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/internal/adapter/presenter"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

// GetBook godoc
// @summary GetBook
// @description get book by id
// @tags book
// @produce json
// @Param id path int true "Book ID"
// @response 200 {object} presenter.SuccessResponse[presenter.BookResponse] "OK"
// @response 400 {object} presenter.ErrorResponse "Bad Request"
// @response 500 {object} presenter.ErrorResponse "Internal Server Error"
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.BadRequestError(err, "id must be an integer")
	}

	book, err := h.bookUseCase.GetByID(id)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "get book service failed")
	}

	res := h.presenter.ToResponse(book)
	return c.JSON(presenter.Success(res))
}

// GetBooks godoc
// @summary GetBooks
// @description get books
// @tags book
// @produce json
// @Param limit query int false "Number of history to be retrieved"
// @Param page query int false "Page to retrieved"
// @response 200 {object} presenter.PaginationResponse[presenter.BookResponse] "OK"
// @response 400 {object} presenter.ErrorResponse "Bad Request"
// @response 500 {object} presenter.ErrorResponse "Internal Server Error"
// @Router /books [get]
func (h *BookHandler) GetBooks(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	if limit > 50 {
		return apperror.BadRequestError(errors.New("limit cannot exceed 50"), "limit cannot exceed 50")
	}
	page := c.QueryInt("page", 1)

	books, totalPage, totalRows, err := h.bookUseCase.List(limit, page)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "get books service failed")
	}

	convertedBooks := h.presenter.ToResponseList(books)
	return c.JSON(presenter.SuccessPagination(convertedBooks, page, totalPage, limit, totalRows))
}
