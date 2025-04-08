package handler

import (
	"github.com/yokeTH/gofiber-template/internal/adapter/presenter"
	"github.com/yokeTH/gofiber-template/internal/usecase/book"
)

type BookHandler struct {
	bookUseCase book.BookUseCase
	presenter   *presenter.BookPresenter
}

func NewBookHandler(bookUseCase book.BookUseCase) *BookHandler {
	return &BookHandler{
		bookUseCase: bookUseCase,
		presenter:   presenter.NewBookPresenter(),
	}
}
