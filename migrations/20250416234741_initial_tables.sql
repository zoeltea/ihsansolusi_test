-- +goose Up
-- Create accounts table
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    nik VARCHAR(20) UNIQUE NOT NULL,
    no_hp VARCHAR(15) NOT NULL,
    no_rekening VARCHAR(20) NOT NULL,
    saldo DECIMAL(15, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create mutations table
CREATE TABLE mutations (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id),
    nominal DECIMAL(15, 2) NOT NULL,
    type VARCHAR(15) NOT NULL, -- 'credit/tabung' or 'debit/tarik'
    reference VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster lookups
CREATE INDEX idx_mutations_account_id ON mutations(account_id);
CREATE INDEX idx_accounts_no_rekening ON accounts(no_rekening);
CREATE INDEX idx_accounts_no_hp ON accounts(no_hp);
CREATE INDEX idx_accounts_nik ON accounts(nik);

-- +goose Down
DROP INDEX IF EXISTS idx_accounts_no_rekening;
DROP INDEX IF EXISTS idx_mutations_account_id;
DROP INDEX IF EXISTS idx_accounts_no_hp;
DROP INDEX IF EXISTS idx_accounts_nik;
DROP TABLE IF EXISTS mutations;
DROP TABLE IF EXISTS accounts;
