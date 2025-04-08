package book

import (
	"github.com/yokeTH/gofiber-template/internal/adapter/presenter"
	"github.com/yokeTH/gofiber-template/internal/domain"
)

type BookRepository interface {
	Create(book *domain.Book) error
	GetByID(id int) (*domain.Book, error)
	List(limit, page int) ([]domain.Book, int, int, error)
	Update(id int, book *presenter.UpdateBookRequest) (*domain.Book, error)
	Delete(id int) error
}

type BookUseCase interface {
	Create(book *domain.Book) error
	GetByID(id int) (*domain.Book, error)
	List(limit, page int) ([]domain.Book, int, int, error)
	Update(id int, book *presenter.UpdateBookRequest) (*domain.Book, error)
	Delete(id int) error
}
