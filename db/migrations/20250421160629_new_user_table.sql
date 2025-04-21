-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reader (
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(255) NOT NULL,
    lastname VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reader;
-- +goose StatementEnd
