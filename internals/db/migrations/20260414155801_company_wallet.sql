-- +goose Up
CREATE TABLE IF NOT EXISTS liquidity_pool (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    balance NUMERIC(20,2) DEFAULT 0,
    currency currency DEFAULT 'NGN',

    updated_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS liquidity_pool;
DROP TYPE IF EXISTS currency;
