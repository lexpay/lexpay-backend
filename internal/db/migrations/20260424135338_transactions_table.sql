-- +goose Up
CREATE TYPE transaction_type AS ENUM(
    'swap'
);
CREATE TYPE transaction_status AS ENUM(
    'deposit_pending',
    'quote_locked',
    'confirming',
    'paying_out',
    'swapping',
    'slippage_review',
    'completed',
    'expired',
    'failed'
);


CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rate_lock_id UUID NOT NULL REFERENCES rate_locks(id) ON DELETE RESTRICT,
    
    --crypto info
    deposit_wallet_address TEXT UNIQUE NOT NULL,
    crypto_amount NUMERIC(36,18) NOT NULL,
    tx_hash TEXT UNIQUE,
    network_confirmations INT DEFAULT 0,
    crypto_currency deposit_currency DEFAULT 'USDT' NOT NULL,

    --settlement info(naira - for now)
    amount_to_receive_currency withdrawal_currency DEFAULT 'NGN' NOT NULL,
    amount_to_receive NUMERIC(20,2) NOT NULL,
    fee BIGINT, -- 1.5% spread
    platform_fee BIGINT, -- platform fee minus paystack and yellow card api cost
    
    --yellow card info
    yellow_card_order_id TEXT,
    yellow_card_status VARCHAR(30),
    
    -- Paystack info
    paystack_transfer_code TEXT,
    paystack_transfer_status VARCHAR(30),
    
    transaction_type transaction_type NOT NULL,
    status transaction_status NOT NULL,

    failure_reason TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);

-- +goose Down
DROP TABLE IF EXISTS transactions;
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS transaction_status;




