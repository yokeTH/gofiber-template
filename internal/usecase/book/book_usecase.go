package book

import (
	"context"

	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/logger"
)

type bookUseCase struct {
	bookRepo BookRepository
}

func NewBookUseCase(bookRepo BookRepository) *bookUseCase {
	return &bookUseCase{
		bookRepo: bookRepo,
	}
}

func (uc *bookUseCase) Create(c context.Context, book *domain.Book) error {
	logger.Func(c, "bookUseCase.Create")
	defer logger.Func(c, "bookUseCase.Create", true)
	logger.Info(c, "usecase: creating book", "book", book)
	err := uc.bookRepo.Create(c, book)
	if err != nil {
		logger.Error(c, "usecase: failed to create book", "error", err, "book", book)
		return err
	}
	logger.Info(c, "usecase: book created successfully", "book", book)
	return nil
}

func (uc *bookUseCase) GetByID(c context.Context, id int) (*domain.Book, error) {
	logger.Func(c, "bookUseCase.GetByID")
	defer logger.Func(c, "bookUseCase.GetByID", true)
	logger.Debug(c, "usecase: fetching book by id", "id", id)
	book, err := uc.bookRepo.GetByID(c, id)
	if err != nil {
		logger.Error(c, "usecase: failed to get book", "error", err, "id", id)
		return nil, err
	}
	logger.Info(c, "usecase: book fetched successfully", "book", book)
	return book, nil
}

func (uc *bookUseCase) List(c context.Context, limit, page int) ([]domain.Book, int, int, error) {
	logger.Func(c, "bookUseCase.List")
	defer logger.Func(c, "bookUseCase.List", true)
	logger.Debug(c, "usecase: listing books", "limit", limit, "page", page)
	books, last, total, err := uc.bookRepo.List(c, limit, page)
	if err != nil {
		logger.Error(c, "usecase: failed to list books", "error", err, "limit", limit, "page", page)
		return nil, 0, 0, err
	}
	logger.Info(c, "usecase: books listed successfully", "count", len(books), "last", last, "total", total)
	return books, last, total, nil
}

func (uc *bookUseCase) Update(c context.Context, id int, bookUpdate *dto.UpdateBookRequest) (*domain.Book, error) {
	logger.Func(c, "bookUseCase.Update")
	defer logger.Func(c, "bookUseCase.Update", true)
	logger.Debug(c, "usecase: updating book", "id", id, "updateRequest", bookUpdate)
	book, err := uc.bookRepo.Update(c, id, bookUpdate)
	if err != nil {
		logger.Error(c, "usecase: failed to update book", "error", err, "id", id, "updateRequest", bookUpdate)
		return nil, err
	}
	logger.Info(c, "usecase: book updated successfully", "book", book)
	return book, nil
}

func (uc *bookUseCase) Delete(c context.Context, id int) error {
	logger.Func(c, "bookUseCase.Delete")
	defer logger.Func(c, "bookUseCase.Delete", true)
	logger.Debug(c, "usecase: deleting book", "id", id)
	err := uc.bookRepo.Delete(c, id)
	if err != nil {
		logger.Error(c, "usecase: failed to delete book", "error", err, "id", id)
		return err
	}
	logger.Info(c, "usecase: book deleted successfully", "id", id)
	return nil
}
