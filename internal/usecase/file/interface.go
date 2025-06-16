package file

import (
	"context"
	"mime/multipart"

	"github.com/yokeTH/gofiber-template/internal/domain"
)

type FileRepository interface {
	Create(c context.Context, file *domain.File) error
	List(c context.Context, limit, page int) ([]domain.File, int, int, error)
	GetByID(c context.Context, id int) (*domain.File, error)
}

type FileUseCase interface {
	CreatePrivateFile(ctx context.Context, file *multipart.FileHeader) (*domain.File, error)
	CreatePublicFile(ctx context.Context, file *multipart.FileHeader) (*domain.File, error)
	List(c context.Context, limit, page int) ([]domain.File, int, int, error)
	GetByID(c context.Context, id int) (*domain.File, error)
}
