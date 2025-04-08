package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/internal/adapter/presenter"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/internal/usecase/book"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

type BookHandler struct {
	bookUseCase book.BookUseCase
	presenter   *presenter.BookPresenter
}

func NewBookHandler(bookUseCase book.BookUseCase) *BookHandler {
	return &BookHandler{
		bookUseCase: bookUseCase,
		presenter:   presenter.NewBookPresenter(),
	}
}

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
