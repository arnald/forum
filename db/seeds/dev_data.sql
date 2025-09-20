-- Users
INSERT OR IGNORE INTO users (id, email, username, password_hash) VALUES
('df16d238-e4dd-4645-9101-54aed9c0fbf4','dev1@forum.test', 'dev_user1', '150000$ZGV2c2FsdDEyMw==$bXzDzL8hQN1qV7z6X0Xj3a8l6y1wY0s3J7xKt8fHfE4='),
('000dec3a-51af-4e7c-ae0c-21436a0a2395','dev2@forum.test', 'dev_user2', '150000$ZGV2c2FsdDEyMw==$bXzDzL8hQN1qV7z6X0Xj3a8l6y1wY0s3J7xKt8fHfE4='),
('f1433622-9c10-44e5-94b1-1f6a148c9131','admin@forum.test', 'forum_admin', '150000$YWRtaW5zYWx0$c2VjcmV0YWRtaW5oYXNo');

-- Sessions
INSERT OR IGNORE INTO sessions (token, user_id, expires_at, refresh_token, refresh_token_expires_at) VALUES
('dev_session_token_1', 'df16d238-e4dd-4645-9101-54aed9c0fbf4', DATETIME('now', '+7 days'), 'dev_refresh_token_1', DATETIME('now', '+30 days')),
('dev_session_token_2', '000dec3a-51af-4e7c-ae0c-21436a0a2395', DATETIME('now', '+7 days'), 'dev_refresh_token_2', DATETIME('now', '+30 days')),
('dev_session_token_3', 'f1433622-9c10-44e5-94b1-1f6a148c9131', DATETIME('now', '+7 days'), 'dev_refresh_token_3', DATETIME('now', '+30 days'));

-- Categories
INSERT OR IGNORE INTO categories (id, name, description) VALUES
('cat1', 'General', 'General discussions'),
('cat2', 'Tech', 'Technology related topics'),
('cat3', 'Sports', 'Sports discussions');

-- Sample Posts
INSERT OR IGNORE INTO posts (id, title, content, user_id, created_at, updated_at) VALUES
('post1', 'Welcome to the Forum!', 'This is the first post on our new forum. Feel free to discuss anything here!', 'df16d238-e4dd-4645-9101-54aed9c0fbf4', DATETIME('now', '-2 days'), DATETIME('now', '-2 days')),
('post2', 'Technology Trends 2025', 'What are your thoughts on the latest technology trends? AI, blockchain, quantum computing...', '000dec3a-51af-4e7c-ae0c-21436a0a2395', DATETIME('now', '-1 day'), DATETIME('now', '-1 day'));

-- Post Categories (many-to-many)
INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES
('post1', 'cat1'),
('post2', 'cat2');

-- Sample Comments
INSERT OR IGNORE INTO comments (id, content, post_id, user_id, parent_id, level, created_at, updated_at) VALUES
('comment1', 'Great post! Thanks for sharing this with the community.', 'post1', '000dec3a-51af-4e7c-ae0c-21436a0a2395', NULL, 0, DATETIME('now', '-1 day'), DATETIME('now', '-1 day')),
('comment2', 'I completely agree with your points here.', 'post1', 'f1433622-9c10-44e5-94b1-1f6a148c9131', 'comment1', 1, DATETIME('now', '-12 hours'), DATETIME('now', '-12 hours')),
('comment3', 'Very interesting perspective on technology trends!', 'post2', 'df16d238-e4dd-4645-9101-54aed9c0fbf4', NULL, 0, DATETIME('now', '-6 hours'), DATETIME('now', '-6 hours')),
('comment4', 'What do you think about AI development in 2025?', 'post2', '000dec3a-51af-4e7c-ae0c-21436a0a2395', NULL, 0, DATETIME('now', '-3 hours'), DATETIME('now', '-3 hours')),
('comment5', 'AI will definitely change everything. Great question!', 'post2', 'f1433622-9c10-44e5-94b1-1f6a148c9131', 'comment4', 1, DATETIME('now', '-2 hours'), DATETIME('now', '-2 hours'));
