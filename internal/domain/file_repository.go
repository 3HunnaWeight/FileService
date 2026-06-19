package domain

import "context"

type FileRepository interface {
	Create(ctx context.Context, file *File) error
	GetByPublicID(ctx context.Context, id string) (*File, error)
	Delete(ctx context.Context, publicID string) error
	List(ctx context.Context, limit, offset int) ([]File, error)
	Count(ctx context.Context) (int, error)
}
