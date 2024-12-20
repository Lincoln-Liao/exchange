CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY,
    "name" TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

INSERT INTO users (user_id, "name", email, created_at, updated_at) VALUES
('00000000-0000-0000-0000-000000000001', 'Alice Johnson', 'alice.johnson@example.com', NOW(), NOW()),
('00000000-0000-0000-0000-000000000002', 'Bob Smith', 'bob.smith@example.com', NOW(), NOW()),
('00000000-0000-0000-0000-000000000003', 'Charlie Brown', 'charlie.brown@example.com', NOW(), NOW()),
('00000000-0000-0000-0000-000000000004', 'Diana Prince', 'diana.prince@example.com', NOW(), NOW()),
('00000000-0000-0000-0000-000000000005', 'Ethan Hunt', 'ethan.hunt@example.com', NOW(), NOW());
