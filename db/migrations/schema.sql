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

-- Sessions
CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    refresh_token TEXT,
    refresh_token_expires_at DATETIME NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Categories
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(id)
);

-- Topics
CREATE TABLE IF NOT EXISTS topics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_path TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Comments
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_id INTEGER NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Votes
CREATE TABLE IF NOT EXISTS votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_id INTEGER REFERENCES topics(id) ON DELETE CASCADE,
    comment_id INTEGER REFERENCES comments(id) ON DELETE CASCADE,
    reaction_type INTEGER NOT NULL CHECK(reaction_type IN (-1, 1)),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, topic_id),
    UNIQUE (user_id, comment_id)
);

-- Notifications
CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    related_type TEXT,
    related_id INTEGER,
    is_read BOOLEAN DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Users table indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Sessions table indexes  
CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions(expires_at);

-- Categories table indexes
CREATE INDEX IF NOT EXISTS idx_categories_created_by ON categories(created_by);

-- Topics table indexes
CREATE INDEX IF NOT EXISTS idx_topics_user ON topics(user_id);
CREATE INDEX IF NOT EXISTS idx_topics_category ON topics(category_id);
CREATE INDEX IF NOT EXISTS idx_topics_created ON topics(created_at DESC);

-- Comments table indexes
CREATE INDEX IF NOT EXISTS idx_comments_topic ON comments(topic_id);
CREATE INDEX IF NOT EXISTS idx_comments_user ON comments(user_id);

-- Votes table indexes (CRITICAL for performance)
-- Prevent duplicate votes
CREATE UNIQUE INDEX IF NOT EXISTS idx_topic_votes 
ON votes(user_id, topic_id) WHERE comment_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_comment_votes 
ON votes(user_id, comment_id) WHERE comment_id IS NOT NULL;

-- Fast vote counting and lookups
CREATE INDEX IF NOT EXISTS idx_votes_topic_reaction 
ON votes(topic_id, reaction_type) WHERE comment_id IS NULL;

CREATE INDEX IF NOT EXISTS idx_votes_comment_reaction 
ON votes(comment_id, reaction_type) WHERE comment_id IS NOT NULL;

-- User activity lookup
CREATE INDEX IF NOT EXISTS idx_votes_user ON votes(user_id);

-- Notifications table indexes
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);