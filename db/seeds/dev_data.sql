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
INSERT OR IGNORE INTO categories (name, description, created_by) VALUES
('General Discussion', 'Topics that don''t fit into other categories', 'df16d238-e4dd-4645-9101-54aed9c0fbf4'),
('Feedback', 'Suggestions and feedback about the forum', '000dec3a-51af-4e7c-ae0c-21436a0a2395'),
('Off-Topic', 'Anything not related to the forum', 'f1433622-9c10-44e5-94b1-1f6a148c9131');

-- Topics
INSERT OR IGNORE INTO topics (user_id, category_id, title, content) VALUES
('df16d238-e4dd-4645-9101-54aed9c0fbf4', 1, 'Welcome to the Forum', 'This is a sample topic created for testing purposes.'),
('000dec3a-51af-4e7c-ae0c-21436a0a2395', 2, 'Feedback on New Features', 'What do you think about the new features?'),
('f1433622-9c10-44e5-94b1-1f6a148c9131', 3, 'Forum Guidelines', 'Please read the forum guidelines before posting.');

-- Comments
INSERT OR IGNORE INTO comments (user_id, topic_id, content) VALUES
('df16d238-e4dd-4645-9101-54aed9c0fbf4', 1, 'This is a comment on the welcome topic.'),
('000dec3a-51af-4e7c-ae0c-21436a0a2395', 2, 'I really like the new features!'),
('f1433622-9c10-44e5-94b1-1f6a148c9131', 3, 'These guidelines are very helpful.');

-- Votes
INSERT OR IGNORE INTO votes (user_id, topic_id, comment_id, reaction_type) VALUES
('df16d238-e4dd-4645-9101-54aed9c0fbf4', 1, NULL, 1),
('000dec3a-51af-4e7c-ae0c-21436a0a2395', 2, NULL, 1),
('f1433622-9c10-44e5-94b1-1f6a148c9131', 3, NULL, -1);