CREATE TABLE IF NOT EXISTS wallets (
    user_id TEXT PRIMARY KEY,
    balance BIGINT NOT NULL,
    currency TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

INSERT INTO wallets (user_id, balance, currency, created_at, updated_at) VALUES
('00000000-0000-0000-0000-000000000001', 10000, 'USD', NOW(), NOW()),
('00000000-0000-0000-0000-000000000002', 20000, 'USD', NOW(), NOW()),
('00000000-0000-0000-0000-000000000003', 30000, 'USD', NOW(), NOW()),
('00000000-0000-0000-0000-000000000004', 40000, 'USD', NOW(), NOW()),
('00000000-0000-0000-0000-000000000005', 0, 'USD', NOW(), NOW());
