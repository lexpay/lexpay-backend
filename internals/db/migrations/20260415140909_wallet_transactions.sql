-- +goose Up
CREATE TYPE transaction_type AS ENUM (
    'deposit',
    'withdrawal',
    'transfer',
    'exchange'
);
CREATE TYPE transaction_status AS ENUM (
    'pending',
    'completed',
    'failed'
);
CREATE TABLE IF NOT EXISTS wallet_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    amount NUMERIC(20,2) NOT NULL,
    currency currency NOT NULL,
    balance_before NUMERIC(20,2) NOT NULL,
    balance_after NUMERIC(20,2) NOT NULL,
    type transaction_type NOT NULL,
    reference_id UUID NOT NULL,
    description TEXT NOT NULL,
    status transaction_status NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS wallet_transactions;
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS transaction_status;
