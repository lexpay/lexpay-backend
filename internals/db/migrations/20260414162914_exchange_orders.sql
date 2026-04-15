-- +goose Up
CREATE TYPE exchange_order_status AS ENUM (
    'pending',
    'completed',
    'cancelled'
);
CREATE TABLE IF NOT EXISTS exchange_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount_from NUMERIC(20,2) NOT NULL,
    amount_to NUMERIC(20,2) NOT NULL,
    rate NUMERIC(20,2) NOT NULL,
    currency_from currency NOT NULL,
    currency_to currency NOT NULL,
    status exchange_order_status NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS exchange_orders;
DROP TYPE IF EXISTS exchange_order_status;
