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

	logger.Info(c, "creating file", "file", book)
	if err := r.db.Create(book).Error; err != nil {
		logger.Error(c, "failed to create file", "error", err, "file", book)
		return apperror.InternalServerError(err, "failed to create book")
	}
	logger.Info(c, "file created successfully", "file", book)
	return nil
}

func (r *fileRepository) List(c context.Context, limit, page int) ([]domain.File, int, int, error) {
	logger.Func(c, "fileRepository.List")
	defer logger.Func(c, "fileRepository.List", true)

	logger.Debug(c, "listing files", "limit", limit, "page", page)
	var files []domain.File
	var total, last int

	if err := r.db.Scopes(db.Paginate(domain.File{}, &limit, &page, &total, &last)).Find(&files).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn(c, "files not found", "limit", limit, "page", page)
			return nil, 0, 0, apperror.NotFoundError(err, "files not found")
		}
		logger.Error(c, "failed to get files", "error", err, "limit", limit, "page", page)
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get files")
	}
	logger.Info(c, "files listed successfully", "count", len(files), "last", last, "total", total)
	return files, last, total, nil
}

func (r *fileRepository) GetByID(c context.Context, id int) (*domain.File, error) {
	logger.Func(c, "fileRepository.GetByID")
	defer logger.Func(c, "fileRepository.GetByID", true)

	logger.Debug(c, "fetching file by id", "id", id)
	file := &domain.File{}
	if err := r.db.First(file, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn(c, "file not found", "id", id)
			return nil, apperror.NotFoundError(err, "file not found")
		}
		logger.Error(c, "failed to get file", "error", err, "id", id)
		return nil, apperror.InternalServerError(err, "failed to get file")
	}
	logger.Info(c, "file fetched successfully", "file", file)
	return file, nil
}
