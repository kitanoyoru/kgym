-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_users_role ON users USING btree (role);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_role;
-- +goose StatementEnd
