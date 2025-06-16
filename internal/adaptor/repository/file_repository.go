package repository

import (
	"context"
	"errors"

	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/db"
	"github.com/yokeTH/gofiber-template/pkg/logger"
	"gorm.io/gorm"
)

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *fileRepository {
	return &fileRepository{
		db: db,
	}
}

func (r *fileRepository) Create(c context.Context, book *domain.File) error {
	logger.Func(c, "fileRepository.Create")
	defer logger.Func(c, "fileRepository.Create", true)

	if err := r.db.Create(book).Error; err != nil {
		return apperror.InternalServerError(err, "failed to create book")
	}
	return nil
}

func (r *fileRepository) List(c context.Context, limit, page int) ([]domain.File, int, int, error) {
	logger.Func(c, "fileRepository.List")
	defer logger.Func(c, "fileRepository.List", true)

	var files []domain.File
	var total, last int

	if err := r.db.Scopes(db.Paginate(domain.File{}, &limit, &page, &total, &last)).Find(&files).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "files not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get files")
	}
	return files, last, total, nil
}

func (r *fileRepository) GetByID(c context.Context, id int) (*domain.File, error) {
	logger.Func(c, "fileRepository.GetByID")
	defer logger.Func(c, "fileRepository.GetByID", true)

	file := &domain.File{}
	if err := r.db.First(file, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError(err, "file not found")
		}
		return nil, apperror.InternalServerError(err, "failed to get file")
	}
	return file, nil
}
