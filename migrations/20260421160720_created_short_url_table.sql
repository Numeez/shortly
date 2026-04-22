-- +goose Up
-- +goose StatementBegin
CREATE TABLE short_url (
    id          BIGSERIAL PRIMARY KEY,
    long_url    TEXT NOT NULL,
    short_url   TEXT NOT NULL UNIQUE,
    expiry_time TIMESTAMPTZ,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE short_url;
-- +goose StatementEnd
