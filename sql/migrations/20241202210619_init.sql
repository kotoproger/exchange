-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS general;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS general;
-- +goose StatementEnd
