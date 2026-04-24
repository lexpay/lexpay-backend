-- +goose Up
CREATE TYPE dividend_status AS ENUM('pending', 'paid', 'failed');

CREATE TABLE IF NOT EXISTS referral_dividends(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    referrer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    referred_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    
    amount NUMERIC(20,2) NOT NULL,
    currency withdrawal_currency DEFAULT 'NGN' NOT NULL, -- Reusing the type from rate_locks
    
    status dividend_status DEFAULT 'pending' NOT NULL,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_referral_dividends_referrer_id ON referral_dividends(referrer_id);

-- +goose Down
DROP TABLE IF EXISTS referral_dividends;
DROP TYPE IF EXISTS dividend_status;
