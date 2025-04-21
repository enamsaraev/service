-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS log (
    id SERIAL PRIMARY KEY,
    created_at DATE NOT NULL,
    reader_id INTEGER REFERENCES reader (id),
    book_id INTEGER REFERENCES book (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS log;
-- +goose StatementEnd
