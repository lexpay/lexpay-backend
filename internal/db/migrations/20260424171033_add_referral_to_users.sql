-- +goose Up
ALTER TABLE users 
ADD COLUMN referral_code VARCHAR(20) UNIQUE,
ADD COLUMN referred_by UUID REFERENCES users(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE users 
DROP COLUMN IF EXISTS referred_by,
DROP COLUMN IF EXISTS referral_code;
