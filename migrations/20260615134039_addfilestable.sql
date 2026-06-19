-- +goose Up
SELECT 'up SQL query';

CREATE TABLE files (
                       id UUID PRIMARY KEY,

                       public_id VARCHAR(64) UNIQUE NOT NULL,

                       original_name TEXT NOT NULL,
                       mime_type TEXT NOT NULL,

                       size_bytes BIGINT NOT NULL,


                       storage_provider VARCHAR(50) NOT NULL,
                       storage_bucket TEXT NOT NULL,
                       storage_key TEXT NOT NULL,

                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       deleted_at TIMESTAMPTZ
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE files;
