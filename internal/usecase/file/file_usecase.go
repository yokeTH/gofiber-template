package file

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/storage"
)

type fileUseCase struct {
	fileRepo   FileRepository
	pubStorage storage.Storage
	priStorage storage.Storage
}

func NewFileUseCase(fileRepo FileRepository, pub, pri storage.Storage) *fileUseCase {
	return &fileUseCase{
		fileRepo:   fileRepo,
		priStorage: pri,
		pubStorage: pub,
	}
}

func (u *fileUseCase) List(limit, page int) ([]domain.File, int, int, error) {
	return u.fileRepo.List(limit, page)
}

func (u *fileUseCase) GetByID(id int) (*domain.File, error) {
	return u.fileRepo.GetByID(id)
}

func (u *fileUseCase) CreatePrivateFile(ctx context.Context, file *multipart.FileHeader) (*domain.File, error) {
	return u.create(ctx, file, domain.PrivateBucketType)
}

func (u *fileUseCase) CreatePublicFile(ctx context.Context, file *multipart.FileHeader) (*domain.File, error) {
	return u.create(ctx, file, domain.PublicBucketType)
}

func (u *fileUseCase) getStorage(bucketType domain.BucketType) storage.Storage {
	if bucketType == domain.PublicBucketType {
		return u.pubStorage
	}
	return u.priStorage
}

func (u *fileUseCase) create(ctx context.Context, file *multipart.FileHeader, bucketType domain.BucketType) (*domain.File, error) {
	fileData, err := file.Open()
	if err != nil {
		return nil, apperror.InternalServerError(err, "error opening file")
	}
	defer fileData.Close()

	filename := strings.ReplaceAll(file.Filename, " ", "-")
	contentType := file.Header.Get("Content-Type")
	fileKey := fmt.Sprintf("upload/%s", filename)

	fileInfo := &domain.File{
		Name:       filename,
		Key:        fileKey,
		BucketType: bucketType,
	}

	if err = u.getStorage(bucketType).UploadFile(ctx, fileKey, contentType, fileData); err != nil {
		return nil, apperror.InternalServerError(err, "error uploading file")
	}

	if err = u.fileRepo.Create(fileInfo); err != nil {
		return nil, apperror.InternalServerError(err, "error create file data")
	}

	return fileInfo, nil
}
