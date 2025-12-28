-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255) NOT NULL;
ALTER TABLE users ADD COLUMN mobile VARCHAR(20) NOT NULL;
ALTER TABLE users ADD COLUMN first_name VARCHAR(255) NOT NULL;
ALTER TABLE users ADD COLUMN last_name VARCHAR(255) NOT NULL;
ALTER TABLE users ADD COLUMN birth_date DATE NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
ALTER TABLE users DROP COLUMN IF EXISTS mobile;
ALTER TABLE users DROP COLUMN IF EXISTS first_name;
ALTER TABLE users DROP COLUMN IF EXISTS last_name;
ALTER TABLE users DROP COLUMN IF EXISTS birth_date;
-- +goose StatementEnd
