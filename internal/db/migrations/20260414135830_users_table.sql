-- +goose Up

CREATE TYPE account_status AS ENUM (
    'active',
    'suspended',
    'frozen',
    'pending_verification'
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    nationality TEXT NOT NULL,
    phone_number TEXT UNIQUE,
    name TEXT NOT NULL,

    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,

    date_of_birth DATE NOT NULL,

    bvn TEXT UNIQUE,

    phone_verified BOOLEAN DEFAULT FALSE,
    email_verified BOOLEAN DEFAULT FALSE,
    kyc_verified BOOLEAN DEFAULT FALSE,

    account_status account_status DEFAULT 'pending_verification',

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS account_status;