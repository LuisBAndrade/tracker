-- +goose Up
UPDATE categories SET color = '#6B7280' WHERE color IS NULL;
ALTER TABLE categories ALTER COLUMN color SET NOT NULL;

-- +goose Down
ALTER TABLE categories ALTER COLUMN color DROP NOT NULL;