-- +goose Up
-- +goose StatementBegin
CREATE TYPE file_extension AS ENUM ('jpg', 'jpeg', 'png', 'gif', 'bmp', 'tiff', 'ico', 'webp', 'svg', 'heic', 'heif');

CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL,

    path VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    extension file_extension NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT files_path_unique UNIQUE (path)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS files CASCADE;
DROP TYPE IF EXISTS file_extension;
-- +goose StatementEnd
