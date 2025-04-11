package repository

import (
	"errors"

	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/db"
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

func (r *fileRepository) Create(book *domain.File) error {
	if err := r.db.Create(book).Error; err != nil {
		return apperror.InternalServerError(err, "failed to create book")
	}
	return nil
}

func (r *fileRepository) List(limit, page int) ([]domain.File, int, int, error) {
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

func (r *fileRepository) GetByID(id int) (*domain.File, error) {
	file := &domain.File{}
	if err := r.db.First(file, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError(err, "file not found")
		}
		return nil, apperror.InternalServerError(err, "failed to get file")
	}
	return file, nil
}
