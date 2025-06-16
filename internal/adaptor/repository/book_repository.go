package repository

import (
	"context"
	"errors"

	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/db"
	"github.com/yokeTH/gofiber-template/pkg/logger"
	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *bookRepository {
	return &bookRepository{
		db: db,
	}
}

func (r *bookRepository) Create(c context.Context, book *domain.Book) error {
	logger.Func(c, "bookRepository.Create")
	defer logger.Func(c, "bookRepository.Create", true)
	logger.Info(c, "creating book", "book", book)
	if err := r.db.Create(book).Error; err != nil {
		logger.Error(c, "failed to create book", "error", err, "book", book)
		return apperror.InternalServerError(err, "failed to create book")
	}
	logger.Info(c, "book created successfully", "book", book)
	return nil
}

func (r *bookRepository) GetByID(c context.Context, id int) (*domain.Book, error) {
	logger.Func(c, "bookRepository.GetByID")
	defer logger.Func(c, "bookRepository.GetByID", true)

	logger.Debug(c, "fetching book by id", "id", id)
	book := &domain.Book{}
	if err := r.db.First(book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn(c, "book not found", "id", id)
			return nil, apperror.NotFoundError(err, "book not found")
		}
		logger.Error(c, "failed to get book", "error", err, "id", id)
		return nil, apperror.InternalServerError(err, "failed to get book")
	}
	logger.Info(c, "book fetched successfully", "book", book)
	return book, nil
}

func (r *bookRepository) List(c context.Context, limit, page int) ([]domain.Book, int, int, error) {
	logger.Func(c, "bookRepository.List")
	defer logger.Func(c, "bookRepository.List", true)

	logger.Debug(c, "listing books", "limit", limit, "page", page)
	var books []domain.Book
	var total, last int

	if err := r.db.Scopes(db.Paginate(domain.Book{}, &limit, &page, &total, &last)).Find(&books).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn(c, "books not found", "limit", limit, "page", page)
			return nil, 0, 0, apperror.NotFoundError(err, "books not found")
		}
		logger.Error(c, "failed to get books", "error", err, "limit", limit, "page", page)
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get books")
	}
	logger.Info(c, "books listed successfully", "count", len(books), "last", last, "total", total)
	return books, last, total, nil
}

func (r *bookRepository) Update(c context.Context, id int, updateRequest *dto.UpdateBookRequest) (*domain.Book, error) {
	logger.Func(c, "bookRepository.Update")
	defer logger.Func(c, "bookRepository.Update", true)

	logger.Debug(c, "updating book", "id", id, "updateRequest", updateRequest)
	var book domain.Book
	if err := r.db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn(c, "book not found for update", "id", id)
			return nil, apperror.NotFoundError(err, "book not found")
		}
		logger.Error(c, "failed to get book for update", "error", err, "id", id)
		return nil, err
	}

	if err := r.db.Model(&book).Updates(updateRequest).Error; err != nil {
		logger.Error(c, "failed to update book", "error", err, "id", id, "updateRequest", updateRequest)
		return nil, err
	}

	logger.Info(c, "book updated successfully", "book", book)
	return &book, nil
}

func (r *bookRepository) Delete(c context.Context, id int) error {
	logger.Func(c, "bookRepository.Delete")
	defer logger.Func(c, "bookRepository.Delete", true)

	logger.Debug(c, "deleting book", "id", id)
	if err := r.db.Delete(&domain.Book{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn(c, "book not found for delete", "id", id)
			return apperror.NotFoundError(err, "book not found")
		}
		logger.Error(c, "failed to delete book", "error", err, "id", id)
		return apperror.InternalServerError(err, "failed to delete book")
	}
	logger.Info(c, "book deleted successfully", "id", id)
	return nil
}
