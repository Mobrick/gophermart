-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS accounts(
    uuid TEXT PRIMARY KEY, 
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS accounts;
-- +goose StatementEnd
