package service

import (
	"github.com/yokeTH/gofiber-template/internal/core/domain"
	"github.com/yokeTH/gofiber-template/internal/core/port"
)

type BookService struct {
	BookRepository port.BookRepository
}

func NewBookService(r port.BookRepository) port.BookService {
	return &BookService{
		BookRepository: r,
	}
}

func (s *BookService) CreateBook(book *domain.Book) error {
	return s.BookRepository.CreateBook(book)
}

func (s *BookService) GetBook(id int) (*domain.Book, error) {
	return s.BookRepository.GetBook(id)
}

func (s *BookService) GetBooks(limit int, page int) ([]*domain.Book, int, int, error) {
	return s.BookRepository.GetBooks(limit, page)
}

func (s *BookService) UpdateBook(id int, book *domain.Book) (*domain.Book, error) {
	return s.BookRepository.UpdateBook(id, book)
}

func (s *BookService) DeleteBook(id int) error {
	return s.BookRepository.DeleteBook(id)
}
