package dto

import (
	"context"
	"time"

	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/storage"
)

type fileDto struct {
	private storage.Storage
	public  storage.Storage
}

type FileDto interface {
	ToResponse(domain.File) (*FileResponse, error)
	ToResponseList(files []domain.File) (*[]FileResponse, error)
}

func NewFileDto(pri, pub storage.Storage) *fileDto {
	return &fileDto{
		private: pri,
		public:  pub,
	}
}

func (f *fileDto) ToResponse(file domain.File) (*FileResponse, error) {
	if file.BucketType == domain.PrivateBucketType {
		url, err := f.private.GetSignedUrl(context.TODO(), file.Key, time.Hour*1)
		if err != nil {
			return &FileResponse{
				Name:      file.Name,
				Url:       "error",
				CreatedAt: &file.CreatedAt,
			}, nil
		}

		return &FileResponse{
			Name:      file.Name,
			Url:       url,
			CreatedAt: &file.CreatedAt,
		}, nil

	} else {
		url, err := f.public.GetPublicUrl(file.Key)
		if err != nil {
			return nil, apperror.InternalServerError(err, "generate url error")
		}
		return &FileResponse{
			Name:      file.Name,
			Url:       url,
			CreatedAt: &file.CreatedAt,
		}, nil
	}
}

func (f *fileDto) ToResponseList(files []domain.File) (*[]FileResponse, error) {
	response := make([]FileResponse, len(files))
	for i, file := range files {
		response[i] = FileResponse{
			ID:   int(file.ID),
			Name: file.Name,
		}
	}
	return &response, nil
}

type FileResponse struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Url       string     `json:"url,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
