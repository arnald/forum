-- Users table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE CHECK(email LIKE '%_@__%.__%'),
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    avatar_url TEXT
);

-- CREATE TABLE IF NOT EXISTS users (
--     id CHAR(36) PRIMARY KEY,
--     username VARCHAR(255) UNIQUE NOT NULL,
--     email VARCHAR(255) UNIQUE NOT NULL,
--     password VARCHAR(255),
--     role VARCHAR(50) NOT NULL DEFAULT 'user',
--     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
--     avatar_url VARCHAR(512)
-- );

-- Sessions
CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- CREATE TABLE IF NOT EXISTS sessions (
--     token CHAR(36) PRIMARY KEY,
--     user_id CHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
--     expiry DATETIME NOT NULL,
--     user_agent TEXT,
--     ip_address TEXT
-- );
