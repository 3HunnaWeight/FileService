package postgres

import (
	"context"

	"github.com/3HunnaWeight/file-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository struct {
	db *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(ctx context.Context, f *domain.File) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO files (
			id, public_id, original_name, mime_type,
			size_bytes, storage_provider, storage_bucket, storage_key, created_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`,
		f.ID,
		f.PublicID,
		f.OriginalName,
		f.MimeType,
		f.SizeBytes,
		f.StorageProvider,
		f.StorageBucket,
		f.StorageKey,
		f.CreatedAt,
	)
	return err
}

func (r *FileRepository) GetByPublicID(
	ctx context.Context,
	publicID string,
) (*domain.File, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, public_id, original_name, mime_type,
		       size_bytes, storage_provider, storage_bucket,
		       storage_key, created_at
		FROM files
		WHERE public_id = $1 AND deleted_at IS NULL
	`, publicID)

	var f domain.File
	err := row.Scan(
		&f.ID,
		&f.PublicID,
		&f.OriginalName,
		&f.MimeType,
		&f.SizeBytes,
		&f.StorageProvider,
		&f.StorageBucket,
		&f.StorageKey,
		&f.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FileRepository) Delete(ctx context.Context, publicID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE files SET deleted_at = NOW()
		WHERE public_id = $1 AND deleted_at IS NULL
	`, publicID)
	return err
}

func (r *FileRepository) List(ctx context.Context, limit, offset int) ([]domain.File, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, public_id, original_name, mime_type,
		       size_bytes, storage_provider, storage_bucket,
		       storage_key, created_at
		FROM files
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []domain.File
	for rows.Next() {
		var f domain.File
		if err := rows.Scan(
			&f.ID,
			&f.PublicID,
			&f.OriginalName,
			&f.MimeType,
			&f.SizeBytes,
			&f.StorageProvider,
			&f.StorageBucket,
			&f.StorageKey,
			&f.CreatedAt,
		); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, rows.Err()
}

func (r *FileRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM files WHERE deleted_at IS NULL
	`).Scan(&count)
	return count, err
}
