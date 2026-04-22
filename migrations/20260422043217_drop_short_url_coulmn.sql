-- +goose Up
-- +goose StatementBegin
ALTER TABLE short_url RENAME TO urls;
ALTER TABLE urls DROP COLUMN short_url;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE urls ADD COLUMN short_url TEXT NOT NULL UNIQUE;
ALTER TABLE urls RENAME TO short_url;
-- +goose StatementEnd