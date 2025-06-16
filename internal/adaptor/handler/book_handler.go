package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/internal/usecase/book"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/logger"
)

type bookHandler struct {
	bookUseCase book.BookUseCase
	dto         dto.BookDto
}

func NewBookHandler(bookUseCase book.BookUseCase) *bookHandler {
	return &bookHandler{
		bookUseCase: bookUseCase,
		dto:         dto.NewBookDto(),
	}
}

// CreateBook godoc
//
//	@summary		CreateBook
//	@description	create book by title and author
//	@tags			book
//	@accept			json
//	@produce 		json
//	@param			Book	body	dto.CreateBookRequest	true	"Book Data"
//	@success 		201	{object}	dto.SuccessResponse[dto.BookResponse]	"Created"
//	@failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@failure 		500	{object}	dto.ErrorResponse	"Internal Server Error"
//
// @Router /books [post]
func (h *bookHandler) CreateBook(c *fiber.Ctx) error {
	logger.Func(c.UserContext(), "bookHandler.CreateBook")
	defer logger.Func(c.UserContext(), "bookHandler.CreateBook", true)

	body := new(dto.CreateBookRequest)
	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	book := &domain.Book{
		Author: body.Author,
		Title:  body.Title,
	}

	if err := h.bookUseCase.Create(c.UserContext(), book); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "create book service failed")
	}

	res := h.dto.ToResponse(book)
	return c.Status(fiber.StatusCreated).JSON(dto.Success(res))
}

// GetBook godoc
//
//	@summary		GetBook
//	@description	get book by id
//	@tags			book
//	@produce		json
//	@Param			id	path	int	true	"Book ID"
//	@success 		200	{object}	dto.SuccessResponse[dto.BookResponse]	"OK"
//	@failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@failure 		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router /books/{id} [get]
func (h *bookHandler) GetBook(c *fiber.Ctx) error {
	logger.Func(c.UserContext(), "bookHandler.GetBook")
	defer logger.Func(c.UserContext(), "bookHandler.GetBook", true)

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.BadRequestError(err, "id must be an integer")
	}

	book, err := h.bookUseCase.GetByID(c.UserContext(), id)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "get book service failed")
	}

	res := h.dto.ToResponse(book)
	return c.JSON(dto.Success(res))
}

// GetBooks godoc
//
//	@summary		GetBooks
//	@description	get books
//	@tags			book
//	@produce		json
//	@Param			limit	query	int	false	"Number of history to be retrieved"
//	@Param			page	query	int	false	"Page to retrieved"
//	@response		200	{object}	dto.PaginationResponse[dto.BookResponse]	"OK"
//	@response		400	{object}	dto.ErrorResponse	"Bad Request"
//	@response		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router /books [get]
func (h *bookHandler) GetBooks(c *fiber.Ctx) error {
	logger.Func(c.UserContext(), "bookHandler.GetBooks")
	defer logger.Func(c.UserContext(), "bookHandler.GetBooks", true)

	limit := c.QueryInt("limit", 10)
	if limit > 50 {
		return apperror.BadRequestError(errors.New("limit cannot exceed 50"), "limit cannot exceed 50")
	}
	page := c.QueryInt("page", 1)

	books, totalPage, totalRows, err := h.bookUseCase.List(c.UserContext(), limit, page)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "get books service failed")
	}

	convertedBooks := h.dto.ToResponseList(books)
	return c.JSON(dto.SuccessPagination(convertedBooks, page, totalPage, limit, totalRows))
}

// UpdateBook godoc
//
//	@summary 		UpdateBook
//	@description	update book data
//	@tags 			book
//	@produce		json
//	@Param 			id		path	int						true	"Book ID"
//	@param 			Book	body	dto.UpdateBookRequest	true	"Book Data"
//	@response 		200	{object}	dto.SuccessResponse[dto.BookResponse]	"OK"
//	@response		400	{object}	dto.ErrorResponse	"Bad Request"
//	@response		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router /books/{id} [patch]
func (h *bookHandler) UpdateBook(c *fiber.Ctx) error {
	logger.Func(c.UserContext(), "bookHandler.UpdateBook")
	defer logger.Func(c.UserContext(), "bookHandler.UpdateBook", true)

	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.BadRequestError(err, "id must be an integer")
	}

	body := new(dto.UpdateBookRequest)
	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	book, err := h.bookUseCase.Update(c.UserContext(), id, body)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "update book service failed")
	}

	res := h.dto.ToResponse(book)
	return c.JSON(dto.Success(res))
}

// DeleteBook godoc
//
//	@summary		DeleteBook
//	@description	delete book by id
//	@tags			book
//	@produce		json
//	@Param			id	path	int	true	"Book ID"
//	@response		204	"No Content"
//	@response		400	{object}	dto.ErrorResponse	"Bad Request"
//	@response		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/books/{id} [delete]
func (h *bookHandler) DeleteBook(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.BadRequestError(err, "id must be an integer")
	}

	if err := h.bookUseCase.Delete(c.UserContext(), id); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "delete book service failed")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
