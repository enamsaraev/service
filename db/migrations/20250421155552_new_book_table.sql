-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS book (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    count INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS book;
-- +goose StatementEnd
