package book

import (
	"context"

	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
)

type BookRepository interface {
	Create(c context.Context, book *domain.Book) error
	GetByID(c context.Context, id int) (*domain.Book, error)
	List(c context.Context, limit, page int) ([]domain.Book, int, int, error)
	Update(c context.Context, id int, book *dto.UpdateBookRequest) (*domain.Book, error)
	Delete(c context.Context, id int) error
}

type BookUseCase interface {
	Create(c context.Context, book *domain.Book) error
	GetByID(c context.Context, id int) (*domain.Book, error)
	List(c context.Context, limit, page int) ([]domain.Book, int, int, error)
	Update(c context.Context, id int, book *dto.UpdateBookRequest) (*domain.Book, error)
	Delete(c context.Context, id int) error
}
