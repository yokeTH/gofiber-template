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

func (r *BookRepository) GetBooks() ([]*domain.Book, error) {
	var books []*domain.Book
	if err := r.db.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepository) UpdateBook(id int, book *domain.Book) (*domain.Book, error) {
	if err := r.db.First(book, id).Error; err != nil {
		return nil, err
	}
	if err := r.db.Save(book).Error; err != nil {
		return nil, err
	}
	return book, nil
}

func (r *BookRepository) DeleteBook(id int) error {
	return r.db.Delete(&domain.Book{}, id).Error
}
