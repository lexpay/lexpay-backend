-- +goose Up
ALTER TABLE users
ALTER COLUMN date_of_birth TYPE DATE USING date_of_birth::DATE;

-- +goose Down
ALTER TABLE users
ALTER COLUMN date_of_birth TYPE TIMESTAMP USING date_of_birth::TIMESTAMP;
