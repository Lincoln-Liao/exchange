CREATE TABLE IF NOT EXISTS transactions (
    id TEXT PRIMARY KEY,
    from_user_id TEXT,
    to_user_id TEXT,
    amount BIGINT NOT NULL,
    currency TEXT NOT NULL,
    transaction_type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_transactions_from_user_id ON transactions (from_user_id);

CREATE INDEX IF NOT EXISTS idx_transactions_to_user_id ON transactions (to_user_id);

CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions (created_at);