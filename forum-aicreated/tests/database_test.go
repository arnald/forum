package tests

import (
	"forum/internal/database"
	"os"
	"testing"
)

func TestNewDB(t *testing.T) {
	dbPath := "test_db.db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	if db.DB == nil {
		t.Fatal("Expected database connection to be established")
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestInit(t *testing.T) {
	dbPath := "test_init.db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	err = db.Init()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	tables := []string{"users", "sessions", "categories", "posts", "post_categories", "comments", "post_likes", "comment_likes"}

	for _, table := range tables {
		var count int
		query := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?"
		err := db.QueryRow(query, table).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to query table existence for %s: %v", table, err)
		}

		if count != 1 {
			t.Fatalf("Expected table %s to exist", table)
		}
	}

	var categoryCount int
	err = db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&categoryCount)
	if err != nil {
		t.Fatalf("Failed to count categories: %v", err)
	}

	if categoryCount == 0 {
		t.Fatal("Expected default categories to be inserted")
	}
}

func TestGetDBPath(t *testing.T) {
	originalPath := os.Getenv("DB_PATH")
	defer os.Setenv("DB_PATH", originalPath)

	os.Unsetenv("DB_PATH")
	path := database.GetDBPath()
	expected := "./data/forum.db"
	if path != expected {
		t.Fatalf("Expected default path %s, got %s", expected, path)
	}

	customPath := "/custom/path/forum.db"
	os.Setenv("DB_PATH", customPath)
	path = database.GetDBPath()
	if path != customPath {
		t.Fatalf("Expected custom path %s, got %s", customPath, path)
	}
}
