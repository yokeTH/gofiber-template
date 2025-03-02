package repository

import (
	"errors"

	"github.com/yokeTH/gofiber-template/internal/core/domain"
	"github.com/yokeTH/gofiber-template/internal/core/port"
	"github.com/yokeTH/gofiber-template/internal/database"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *database.Database
}

func NewBookRepository(db *database.Database) port.BookRepository {
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

func (r *BookRepository) GetBooks(limit int, page int) ([]*domain.Book, int, int, error) {
	var books []*domain.Book

	totalPage, totalRows, err := r.db.Paginate(&books, limit, page, "id asc")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "books not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get books")
	}
	return books, totalPage, totalRows, nil
}

func (r *BookRepository) UpdateBook(id int, book *domain.Book) (*domain.Book, error) {
	if err := r.db.Where("id = ?", id).Updates(book).First(book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError(err, "book not found")
		}
		return nil, err
	}
	return book, nil
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
