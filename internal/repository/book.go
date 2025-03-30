package repository

import (
	"errors"

	"github.com/yokeTH/gofiber-template/internal/core/domain"
	"github.com/yokeTH/gofiber-template/internal/db"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/dto"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

func (r *BookRepository) CreateBook(book *domain.Book) error {
	if err := r.db.Create(book).Error; err != nil {
		return apperror.InternalServerError(err, "failed to create book")
	}
	return nil
}

func (r *BookRepository) GetBook(id int) (*domain.Book, error) {
	book := &domain.Book{}
	if err := r.db.First(book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError(err, "book not found")
		}
		return nil, apperror.InternalServerError(err, "failed to get book")
	}
	return book, nil
}

func (r *BookRepository) GetBooks(limit, page, total, last *int) ([]domain.Book, error) {
	var books []domain.Book

	if err := r.db.Scopes(db.Paginate(domain.Book{}, limit, page, total, last)).Find(&books).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError(err, "books not found")
		}
		return nil, apperror.InternalServerError(err, "failed to get books")
	}
	return books, nil
}

func (r *BookRepository) UpdateBook(id int, updateRequest *dto.UpdateBookRequest) (*domain.Book, error) {
	var book domain.Book
	if err := r.db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError(err, "book not found")
		}
		return nil, err
	}

	if err := r.db.Model(&book).Updates(updateRequest).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *BookRepository) DeleteBook(id int) error {
	if err := r.db.Delete(&domain.Book{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFoundError(err, "book not found")
		}
		return apperror.InternalServerError(err, "failed to delete book")
	}
	return nil
}
