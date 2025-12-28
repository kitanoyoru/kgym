-- +goose Up
-- +goose StatementBegin
CREATE TYPE token_type AS ENUM ('refresh');

CREATE TABLE IF NOT EXISTS tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject VARCHAR(255) NOT NULL,
    client_id VARCHAR(255) NOT NULL,
    token_type token_type NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT tokens_token_hash_unique UNIQUE (token_hash)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS token_type;
DROP TABLE IF EXISTS tokens CASCADE;
-- +goose StatementEnd
