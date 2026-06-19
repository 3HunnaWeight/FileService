package service

import (
	"context"
	"io"

	"github.com/3HunnaWeight/file-service/internal/domain"
	"github.com/google/uuid"
)

type Storage interface {
	Upload(ctx context.Context, key string, data []byte, mime string) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}

type FileService struct {
	repo    domain.FileRepository
	storage Storage
	bucket  string
}

func NewFileService(repo domain.FileRepository, storage Storage, bucket string) *FileService {
	return &FileService{repo: repo, storage: storage, bucket: bucket}
}

func (s *FileService) Upload(ctx context.Context, name string, data []byte, mime string) (string, error) {
	id := uuid.New().String()
	publicID := uuid.New().String()
	key := "uploads/" + id + "/" + name

	err := s.repo.Create(ctx, &domain.File{
		ID:              id,
		PublicID:        publicID,
		OriginalName:    name,
		MimeType:        mime,
		SizeBytes:       int64(len(data)),
		StorageProvider: "s3",
		StorageBucket:   s.bucket,
		StorageKey:      key,
	})
	if err != nil {
		return "", err
	}

	if err := s.storage.Upload(ctx, key, data, mime); err != nil {
		return "", err
	}

	return publicID, nil
}

func (s *FileService) GetByPublicID(ctx context.Context, publicID string) (*domain.File, error) {
	return s.repo.GetByPublicID(ctx, publicID)
}

func (s *FileService) Download(
	ctx context.Context,
	publicID string,
) (io.ReadCloser, *domain.File, error) {
	file, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		return nil, nil, err
	}

	reader, err := s.storage.Download(ctx, file.StorageKey)
	if err != nil {
		return nil, nil, err
	}

	return reader, file, nil
}

func (s *FileService) Delete(ctx context.Context, publicID string) error {
	file, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		return err
	}

	if err := s.storage.Delete(ctx, file.StorageKey); err != nil {
		return err
	}

	return s.repo.Delete(ctx, publicID)
}

func (s *FileService) List(ctx context.Context, limit, offset int) ([]domain.File, int, error) {
	files, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return files, total, nil
}
