-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255) NOT NULL,
    avatar VARCHAR(255),
    role VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_email ON users(email);

CREATE TABLE clients (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    ci VARCHAR(255),
    sex VARCHAR(255),
    birth_date VARCHAR(255),
    user_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_clients_deleted_at ON clients(deleted_at);
CREATE UNIQUE INDEX idx_clients_user_id_ci ON clients(user_id, ci) WHERE deleted_at IS NULL;

CREATE TABLE type_accounts (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_type_accounts_deleted_at ON type_accounts(deleted_at);
CREATE INDEX idx_type_accounts_name ON type_accounts(name);

CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    number VARCHAR(255),
    interest NUMERIC,
    balance NUMERIC,
    status VARCHAR(255),
    client_id UUID REFERENCES clients(id),
    type_account_id UUID REFERENCES type_accounts(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_accounts_deleted_at ON accounts(deleted_at);
CREATE INDEX idx_accounts_client_id ON accounts(client_id);
CREATE INDEX idx_accounts_type_account_id ON accounts(type_account_id);


CREATE TABLE type_operations (
    id UUID PRIMARY KEY,
    code VARCHAR(255),
    name VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_type_operations_deleted_at ON type_operations(deleted_at);
CREATE INDEX idx_type_operations_code ON type_operations(code);
CREATE INDEX idx_type_operations_name ON type_operations(name);

CREATE TABLE denominations (
    id UUID PRIMARY KEY,
    type VARCHAR(255),
    value NUMERIC,
    name VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_denominations_deleted_at ON denominations(deleted_at);
CREATE INDEX idx_denominations_value ON denominations(value);

CREATE TABLE cash_sessions (
    id UUID PRIMARY KEY,
    state VARCHAR(255),
    opening_date TIMESTAMPTZ,
    opening_amount NUMERIC,
    closing_date TIMESTAMPTZ,
    closing_amount NUMERIC,
    expected_amount NUMERIC,
    difference_amount NUMERIC,
    user_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_cash_sessions_deleted_at ON cash_sessions(deleted_at);
CREATE INDEX idx_cash_sessions_user_id ON cash_sessions(user_id);
CREATE INDEX idx_cash_sessions_opening_date ON cash_sessions(opening_date);
CREATE INDEX idx_cash_sessions_closing_date ON cash_sessions(closing_date);

CREATE TABLE bank_operations (
    id UUID PRIMARY KEY,
    code VARCHAR(255),
    date TIMESTAMPTZ,
    previous_balance NUMERIC,
    import NUMERIC,
    end_balance NUMERIC,
    type_operation_id UUID REFERENCES type_operations(id),
    cash_session_id UUID REFERENCES cash_sessions(id),
    account_id UUID REFERENCES accounts(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_bank_operations_deleted_at ON bank_operations(deleted_at);
CREATE INDEX idx_bank_operations_code ON bank_operations(code);
CREATE INDEX idx_bank_operations_date ON bank_operations(date);
CREATE INDEX idx_bank_operations_type_operation_id ON bank_operations(type_operation_id);
CREATE INDEX idx_bank_operations_cash_session_id ON bank_operations(cash_session_id);
CREATE INDEX idx_bank_operations_account_id ON bank_operations(account_id);

CREATE TABLE cash_counts (
    id UUID PRIMARY KEY,
    type VARCHAR(255),
    quantity INTEGER,
    subtotal NUMERIC,
    denomination_id UUID NOT NULL REFERENCES denominations(id),
    cash_session_id UUID NOT NULL REFERENCES cash_sessions(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_cash_counts_deleted_at ON cash_counts(deleted_at);
CREATE INDEX idx_cash_counts_type ON cash_counts(type);
CREATE INDEX idx_cash_counts_cash_session_id ON cash_counts(cash_session_id);

CREATE TABLE operation_informations (
    id UUID PRIMARY KEY,
    origin VARCHAR(255),
    reason VARCHAR(255),
    destination VARCHAR(255),
    details VARCHAR(255),
    bank_operation_id UUID REFERENCES bank_operations(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_operation_informations_deleted_at ON operation_informations(deleted_at);
CREATE INDEX idx_operation_informations_bank_operation_id ON operation_informations(bank_operation_id);

-- +goose Down
DROP TABLE IF EXISTS operation_informations;
DROP TABLE IF EXISTS cash_counts;
DROP TABLE IF EXISTS bank_operations;
DROP TABLE IF EXISTS cash_sessions;
DROP TABLE IF EXISTS denominations;
DROP TABLE IF EXISTS type_operations;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS type_accounts;
DROP TABLE IF EXISTS clients;
DROP TABLE IF EXISTS users;
