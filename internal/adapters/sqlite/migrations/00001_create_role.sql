-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS role (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS role;
-- +goose StatementEnd
