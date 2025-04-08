package book

import (
	"github.com/yokeTH/gofiber-template/internal/adapter/presenter"
	"github.com/yokeTH/gofiber-template/internal/domain"
)

type bookUseCase struct {
	bookRepo BookRepository
}

func NewBookUseCase(bookRepo BookRepository) BookUseCase {
	return &bookUseCase{
		bookRepo: bookRepo,
	}
}

func (uc *bookUseCase) Create(book *domain.Book) error {
	return uc.bookRepo.Create(book)
}

func (uc *bookUseCase) GetByID(id int) (*domain.Book, error) {
	return uc.bookRepo.GetByID(id)
}

func (uc *bookUseCase) List(limit, page int) ([]domain.Book, int, int, error) {
	return uc.bookRepo.List(limit, page)
}

func (uc *bookUseCase) Update(id int, bookUpdate *presenter.UpdateBookRequest) (*domain.Book, error) {
	return uc.bookRepo.Update(id, bookUpdate)
}

func (uc *bookUseCase) Delete(id int) error {
	return uc.bookRepo.Delete(id)
}
