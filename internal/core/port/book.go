package port

import "github.com/yokeTH/gofiber-template/internal/core/domain"

type BookService interface {
	CreateBook(book *domain.Book) error
	GetBook(id int) (*domain.Book, error)
	GetBooks() ([]*domain.Book, error)
	UpdateBook(id int, book *domain.Book) (*domain.Book, error)
	DeleteBook(id int) error
}

type BookRepository interface {
	CreateBook(book *domain.Book) error
	GetBook(id int) (*domain.Book, error)
	GetBooks() ([]*domain.Book, error)
	UpdateBook(id int, book *domain.Book) (*domain.Book, error)
	DeleteBook(id int) error
}
