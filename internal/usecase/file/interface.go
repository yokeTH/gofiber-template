package file

import (
	"context"
	"mime/multipart"

	"github.com/yokeTH/gofiber-template/internal/domain"
)

type FileRepository interface {
	Create(file *domain.File) error
	List(limit, page int) ([]domain.File, int, int, error)
	GetByID(id int) (*domain.File, error)
}

type FileUseCase interface {
	CreatePrivateFile(ctx context.Context, file *multipart.FileHeader) (*domain.File, error)
	CreatePublicFile(ctx context.Context, file *multipart.FileHeader) (*domain.File, error)
	List(limit, page int) ([]domain.File, int, int, error)
	GetByID(id int) (*domain.File, error)
}
