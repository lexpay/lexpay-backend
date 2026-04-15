-- +goose Up
CREATE TYPE currency AS ENUM (
    'NGN',
    'USD',
    'GBP',
    'EUR'
);

CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    balance NUMERIC(20,2) DEFAULT 0,
    currency currency DEFAULT 'NGN',
    created_at TIMESTAMP DEFAULT NOW(),


    UNIQUE(user_id,currency)
);

-- +goose Down
DROP TYPE IF EXISTS currency;
DROP TABLE IF EXISTS wallets;
