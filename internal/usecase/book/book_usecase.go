package book

import (
	"context"

	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/logger"
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
	logger.Func(c, "bookUseCase.Create")
	defer logger.Func(c, "bookUseCase.Create", true)
	return uc.bookRepo.Create(c, book)
}

func (uc *bookUseCase) GetByID(c context.Context, id int) (*domain.Book, error) {
	logger.Func(c, "bookUseCase.GetByID")
	defer logger.Func(c, "bookUseCase.GetByID", true)
	return uc.bookRepo.GetByID(c, id)
}

func (uc *bookUseCase) List(c context.Context, limit, page int) ([]domain.Book, int, int, error) {
	logger.Func(c, "bookUseCase.List")
	defer logger.Func(c, "bookUseCase.List", true)
	return uc.bookRepo.List(c, limit, page)
}

func (uc *bookUseCase) Update(c context.Context, id int, bookUpdate *dto.UpdateBookRequest) (*domain.Book, error) {
	logger.Func(c, "bookUseCase.Update")
	defer logger.Func(c, "bookUseCase.Update", true)
	return uc.bookRepo.Update(c, id, bookUpdate)
}

func (uc *bookUseCase) Delete(c context.Context, id int) error {
	logger.Func(c, "bookUseCase.Delete")
	defer logger.Func(c, "bookUseCase.Delete", true)
	return uc.bookRepo.Delete(c, id)
}
