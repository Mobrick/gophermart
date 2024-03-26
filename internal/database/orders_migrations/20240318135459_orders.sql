-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders(
    number TEXT PRIMARY KEY,
    account_uuid TEXT NOT NULL, 
	status TEXT,
	accrual TEXT NOT NULL DEFAULT 0,
    uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    proceeded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
