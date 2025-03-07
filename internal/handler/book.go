package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/internal/core/domain"
	"github.com/yokeTH/gofiber-template/internal/core/port"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/dto"
)

type BookHandler struct {
	BookService port.BookService
}

func NewBookHandler(bookService port.BookService) port.BookHandler {
	return &BookHandler{
		BookService: bookService,
	}
}

// CreateBook godoc
// @summary CreateBook
// @description create book by title and author
// @tags book
// @accept json
// @produce json
// @param Book body domain.Book true "Book Data"
// @response 201 {object} dto.SuccessResponse[domain.Book] "Created"
// @response 400 {object} apperror.AppError "Bad Request"
// @response 500 {object} apperror.AppError "Internal Server Error"
// @Router /books [post]
func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	body := new(domain.Book)
	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	if err := h.BookService.CreateBook(body); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "create book service failed")
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(*body))
}

// GetBook godoc
// @summary GetBook
// @description get book by id
// @tags book
// @produce json
// @Param id path int true "Book ID"
// @response 200 {object} dto.SuccessResponse[domain.Book] "OK"
// @response 400 {object} apperror.AppError "Bad Request"
// @response 500 {object} apperror.AppError "Internal Server Error"
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	book, err := h.BookService.GetBook(id)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "get book service failed")
	}

	return c.JSON(dto.Success(*book))
}

// GetBooks godoc
// @summary GetBooks
// @description get books
// @tags book
// @produce json
// @Param limit query int false "Number of history to be retrieved"
// @Param page query int false "Page to retrieved"
// @response 200 {object} dto.PaginationResponse[domain.Book] "OK"
// @response 400 {object} apperror.AppError "Bad Request"
// @response 500 {object} apperror.AppError "Internal Server Error"
// @Router /books [get]
func (h *BookHandler) GetBooks(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)
	books, totalPage, totalRows, err := h.BookService.GetBooks(limit, page)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "get books service failed")
	}

	convertedBooks := make([]domain.Book, len(books))
	for i, book := range books {
		convertedBooks[i] = *book
	}
	return c.JSON(dto.SuccessPagination(convertedBooks, page, totalPage, limit, totalRows))
}

// UpdateBook godoc
// @summary UpdateBook
// @description update book data
// @tags book
// @produce json
// @Param id path int true "Book ID"
// @param Book body domain.Book true "Book Data"
// @response 200 {object} dto.SuccessResponse[domain.Book] "OK"
// @response 400 {object} apperror.AppError "Bad Request"
// @response 500 {object} apperror.AppError "Internal Server Error"
// @Router /books/{id} [patch]
func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	body := new(domain.Book)
	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	book, err := h.BookService.UpdateBook(id, body)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "update book service failed")
	}

	return c.JSON(dto.Success(*book))
}

// DeleteBook godoc
// @summary UpdateBook
// @description update book data
// @tags book
// @produce json
// @Param id path int true "Book ID"
// @response 200 {object} dto.SuccessResponse[domain.Book] "OK"
// @response 400 {object} apperror.AppError "Bad Request"
// @response 500 {object} apperror.AppError "Internal Server Error"
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	if err := h.BookService.DeleteBook(id); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "delete book service failed")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
