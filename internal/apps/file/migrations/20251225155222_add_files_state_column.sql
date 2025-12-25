-- +goose Up
-- +goose StatementBegin
CREATE TYPE file_state AS ENUM ('pending', 'completed', 'failed');

ALTER TABLE files ADD COLUMN state file_state NOT NULL DEFAULT 'pending';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS file_state;

ALTER TABLE files DROP COLUMN IF EXISTS state;
-- +goose StatementEnd
