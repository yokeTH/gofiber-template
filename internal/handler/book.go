package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/internal/core/domain"
	"github.com/yokeTH/gofiber-template/internal/core/port"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/response"
)

type BookHandler struct {
	BookService port.BookService
}

func NewBookHandler(bookService port.BookService) port.BookHandler {
	return &BookHandler{
		BookService: bookService,
	}
}

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

	return c.Status(fiber.StatusCreated).JSON(response.Success(*body))
}

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

	return c.JSON(response.Success(*book))
}

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
	return c.JSON(response.SuccessPagination(convertedBooks, page, totalPage, limit, totalRows))
}

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

	return c.JSON(response.Success(*book))
}

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
