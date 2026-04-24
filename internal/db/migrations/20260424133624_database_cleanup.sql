-- +goose Up
DROP TABLE IF EXISTS exchange_orders;
DROP TABLE IF EXISTS wallet_transactions;
DROP TABLE IF EXISTS liquidity_pool;
DROP TABLE IF EXISTS wallets;

-- +goose Down

