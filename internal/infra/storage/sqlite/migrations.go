package sqlite

import (
	"fmt"
	"log"
)

func (r Repo) CreateUserTable() error {
	_, err := r.DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id CHAR(36) PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255),
		role VARCHAR(50) NOT NULL DEFAULT 'user',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		avatar_url VARCHAR(512)
	);
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = r.DB.Exec(`
	CREATE TABLE IF NOT EXISTS sessions (
		token CHAR(36) PRIMARY KEY,
		user_id CHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		expiry DATETIME NOT NULL,
		user_agent TEXT,
		ip_address TEXT
	);
	`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Tables created successfully.")

	return nil
}

func (r Repo) CreateSessionsTable() error {
	if _, err := r.DB.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	query := `
CREATE TABLE IF NOT EXISTS sessions (
	token CHAR(36) PRIMARY KEY,
	user_id CHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	expiry DATETIME NOT NULL,
	user_agent TEXT,
	ip_address TEXT
);
`
	_, err := r.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create sessions table: %w", err)
	}
	return nil
}
