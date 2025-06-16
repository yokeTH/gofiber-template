package book

import (
	"context"

	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
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
	return uc.bookRepo.Create(c, book)
}

func (uc *bookUseCase) GetByID(c context.Context, id int) (*domain.Book, error) {
	return uc.bookRepo.GetByID(c, id)
}

func (uc *bookUseCase) List(c context.Context, limit, page int) ([]domain.Book, int, int, error) {
	return uc.bookRepo.List(c, limit, page)
}

func (uc *bookUseCase) Update(c context.Context, id int, bookUpdate *dto.UpdateBookRequest) (*domain.Book, error) {
	return uc.bookRepo.Update(c, id, bookUpdate)
}

func (uc *bookUseCase) Delete(c context.Context, id int) error {
	return uc.bookRepo.Delete(c, id)
}
