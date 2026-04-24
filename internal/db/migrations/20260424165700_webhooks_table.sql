-- +goose Up
CREATE TYPE webhook_status AS ENUM('pending', 'processed', 'failed');

CREATE TABLE IF NOT EXISTS webhooks(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider VARCHAR(50) NOT NULL, -- e.g., 'yellowcard', 'paystack'
    event_type VARCHAR(100) NOT NULL, -- e.g., 'transfer.success'
    payload JSONB NOT NULL,
    status webhook_status DEFAULT 'pending' NOT NULL,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_webhooks_provider ON webhooks(provider);
CREATE INDEX IF NOT EXISTS idx_webhooks_status ON webhooks(status);

-- +goose Down
DROP TABLE IF EXISTS webhooks;
DROP TYPE IF EXISTS webhook_status;
