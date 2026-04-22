-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX unique_long_url ON urls (long_url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX unique_long_url;
-- +goose StatementEnd
