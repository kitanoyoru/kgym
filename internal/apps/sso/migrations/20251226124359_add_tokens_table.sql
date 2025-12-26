-- +goose Up
-- +goose StatementBegin
CREATE TYPE token_type AS ENUM ('refresh');

CREATE TABLE IF NOT EXISTS tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    token_type token_type NOT NULL,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    CONSTRAINT tokens_user_id_unique UNIQUE (user_id),
    CONSTRAINT tokens_token_unique UNIQUE (token)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS token_type;
DROP TABLE IF EXISTS tokens CASCADE;
-- +goose StatementEnd
