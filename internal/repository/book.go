package repository

import (
	"github.com/yokeTH/gofiber-template/internal/core/domain"
	"github.com/yokeTH/gofiber-template/internal/core/port"
	"github.com/yokeTH/gofiber-template/internal/database"
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
	return r.db.Create(book).Error
}

func (r *BookRepository) GetBook(id int) (*domain.Book, error) {
	book := &domain.Book{}
	if err := r.db.First(book, id).Error; err != nil {
		return nil, err
	}
	return book, nil
}

func (r *BookRepository) GetBooks(limit int, page int) ([]*domain.Book, int, int, error) {
	var books []*domain.Book

	totalPage, totalRows, err := r.db.Paginate(&books, limit, page, "id asc")
	if err != nil {
		return nil, 0, 0, err
	}
	return books, totalPage, totalRows, nil
}

func (r *BookRepository) UpdateBook(id int, book *domain.Book) (*domain.Book, error) {
	if err := r.db.Where("id = ?", id).Updates(book).First(book, id).Error; err != nil {
		return nil, err
	}
	return book, nil
}

func (r *BookRepository) DeleteBook(id int) error {
	return r.db.Delete(&domain.Book{}, id).Error
}
