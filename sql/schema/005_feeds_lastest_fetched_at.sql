-- +goose Up
ALTER TABLE feeds ADD COLUMN latest_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE feeds DROP COLUMN latest_fetched_at;
