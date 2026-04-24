-- +goose Up
CREATE TYPE rate_lock_status AS ENUM(
    'active',
    'expired',
    'used',
    'cancelled'
);

CREATE TYPE deposit_currency AS ENUM('USDT');
CREATE TYPE withdrawal_currency AS ENUM('NGN');

CREATE TABLE IF NOT EXISTS rate_locks(
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
   
   --crypto info
   crypto_amount NUMERIC(36,18) NOT NULL,
   crypto_currency deposit_currency DEFAULT 'USDT' NOT NULL,
   
   --settlement info
   amount NUMERIC(20,2) NOT NULL,
   currency withdrawal_currency DEFAULT 'NGN' NOT NULL,

   rate NUMERIC(36,18) NOT NULL,

   status rate_lock_status DEFAULT 'active' NOT NULL,
   expires_at TIMESTAMPTZ NOT NULL,

   created_at TIMESTAMPTZ DEFAULT NOW(),
   updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS rate_locks;
DROP TYPE IF EXISTS rate_lock_status;
DROP TYPE IF EXISTS withdrawal_currency;
DROP TYPE IF EXISTS deposit_currency;