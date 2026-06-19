package domain

import "time"

type File struct {
	ID string

	PublicID string

	OriginalName string
	MimeType     string

	SizeBytes int64

	StorageProvider string
	StorageBucket   string
	StorageKey      string

	CreatedAt   time.Time
	DeletedAt   *time.Time
}
