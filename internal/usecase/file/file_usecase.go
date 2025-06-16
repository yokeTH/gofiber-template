package file

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/logger"
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

func (u *fileUseCase) List(c context.Context, limit, page int) ([]domain.File, int, int, error) {
	logger.Func(c, "fileUseCase.List")
	defer logger.Func(c, "fileUseCase.List", true)
	logger.Debug(c, "listing files", "limit", limit, "page", page)
	files, last, total, err := u.fileRepo.List(c, limit, page)
	if err != nil {
		logger.Error(c, "failed to list files", "error", err, "limit", limit, "page", page)
		return nil, 0, 0, err
	}
	logger.Info(c, "files listed successfully", "count", len(files), "last", last, "total", total)
	return files, last, total, nil
}

func (u *fileUseCase) GetByID(c context.Context, id int) (*domain.File, error) {
	logger.Func(c, "fileUseCase.GetByID")
	defer logger.Func(c, "fileUseCase.GetByID", true)
	logger.Debug(c, "fetching file by id", "id", id)
	file, err := u.fileRepo.GetByID(c, id)
	if err != nil {
		logger.Error(c, "failed to get file", "error", err, "id", id)
		return nil, err
	}
	logger.Info(c, "file fetched successfully", "file", file)
	return file, nil
}

func (u *fileUseCase) CreatePrivateFile(ctx context.Context, file *multipart.FileHeader) (*domain.File, error) {
	logger.Func(ctx, "fileUseCase.CreatePrivateFile")
	defer logger.Func(ctx, "fileUseCase.CreatePrivateFile", true)
	logger.Info(ctx, "creating private file", "filename", file.Filename)
	result, err := u.create(ctx, file, domain.PrivateBucketType)
	if err != nil {
		logger.Error(ctx, "failed to create private file", "error", err, "filename", file.Filename)
		return nil, err
	}
	logger.Info(ctx, "private file created successfully", "file", result)
	return result, nil
}

func (u *fileUseCase) CreatePublicFile(ctx context.Context, file *multipart.FileHeader) (*domain.File, error) {
	logger.Func(ctx, "fileUseCase.CreatePublicFile")
	defer logger.Func(ctx, "fileUseCase.CreatePublicFile", true)
	logger.Info(ctx, "creating public file", "filename", file.Filename)
	result, err := u.create(ctx, file, domain.PublicBucketType)
	if err != nil {
		logger.Error(ctx, "failed to create public file", "error", err, "filename", file.Filename)
		return nil, err
	}
	logger.Info(ctx, "public file created successfully", "file", result)
	return result, nil
}

func (u *fileUseCase) getStorage(c context.Context, bucketType domain.BucketType) storage.Storage {
	logger.Func(c, "fileUseCase.getStorage")
	defer logger.Func(c, "fileUseCase.getStorage", true)
	logger.Debug(c, "getting storage for bucket type", "bucketType", bucketType)
	if bucketType == domain.PublicBucketType {
		return u.pubStorage
	}
	return u.priStorage
}

func (u *fileUseCase) create(ctx context.Context, file *multipart.FileHeader, bucketType domain.BucketType) (*domain.File, error) {
	logger.Func(ctx, "fileUseCase.create")
	defer logger.Func(ctx, "fileUseCase.create", true)
	logger.Debug(ctx, "opening file for upload", "filename", file.Filename, "bucketType", bucketType)
	fileData, err := file.Open()
	if err != nil {
		logger.Error(ctx, "error opening file", "error", err, "filename", file.Filename)
		return nil, apperror.InternalServerError(err, "error opening file")
	}
	defer func() {
		_ = fileData.Close()
	}()

	filename := strings.ReplaceAll(file.Filename, " ", "-")
	contentType := file.Header.Get("Content-Type")
	fileKey := fmt.Sprintf("upload/%s", filename)

	fileInfo := &domain.File{
		Name:       filename,
		Key:        fileKey,
		BucketType: bucketType,
	}

	logger.Debug(ctx, "uploading file to storage", "fileKey", fileKey, "contentType", contentType, "bucketType", bucketType)
	if err = u.getStorage(ctx, bucketType).UploadFile(ctx, fileKey, contentType, fileData); err != nil {
		logger.Error(ctx, "error uploading file", "error", err, "fileKey", fileKey)
		return nil, apperror.InternalServerError(err, "error uploading file")
	}

	logger.Debug(ctx, "creating file record in repository", "fileInfo", fileInfo)
	if err = u.fileRepo.Create(ctx, fileInfo); err != nil {
		logger.Error(ctx, "error creating file data", "error", err, "fileInfo", fileInfo)
		return nil, apperror.InternalServerError(err, "error create file data")
	}

	logger.Info(ctx, "file created and uploaded successfully", "file", fileInfo)
	return fileInfo, nil
}
