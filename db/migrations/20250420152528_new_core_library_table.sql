-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS library (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS library;
-- +goose StatementEnd
