package tests

import (
	"forum/internal/auth"
	"forum/internal/database"
	"os"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) *database.DB {
	dbPath := "test_forum.db"

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	if err := db.Init(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		os.Remove(dbPath)
	})

	return db
}

func TestHashPassword(t *testing.T) {
	db := setupTestDB(t)
	authService := auth.NewAuth(db)

	password := "testpassword123"
	hash, err := authService.HashPassword(password)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatal("Expected hash to be non-empty")
	}

	if hash == password {
		t.Fatal("Expected hash to be different from password")
	}
}

func TestCheckPassword(t *testing.T) {
	db := setupTestDB(t)
	authService := auth.NewAuth(db)

	password := "testpassword123"
	hash, err := authService.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !authService.CheckPassword(password, hash) {
		t.Fatal("Expected password check to pass")
	}

	if authService.CheckPassword("wrongpassword", hash) {
		t.Fatal("Expected password check to fail for wrong password")
	}
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	authService := auth.NewAuth(db)

	email := "test@example.com"
	username := "testuser"
	password := "testpassword123"

	user, err := authService.CreateUser(email, username, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Fatal("Expected user ID to be set")
	}

	if user.Email != email {
		t.Fatalf("Expected email %s, got %s", email, user.Email)
	}

	if user.Username != username {
		t.Fatalf("Expected username %s, got %s", username, user.Username)
	}

	if user.Password == password {
		t.Fatal("Expected password to be hashed")
	}
}

func TestGetUserByEmail(t *testing.T) {
	db := setupTestDB(t)
	authService := auth.NewAuth(db)

	email := "test@example.com"
	username := "testuser"
	password := "testpassword123"

	createdUser, err := authService.CreateUser(email, username, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrievedUser, err := authService.GetUserByEmail(email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}

	if retrievedUser.ID != createdUser.ID {
		t.Fatalf("Expected user ID %d, got %d", createdUser.ID, retrievedUser.ID)
	}

	if retrievedUser.Email != email {
		t.Fatalf("Expected email %s, got %s", email, retrievedUser.Email)
	}
}

func TestCreateSession(t *testing.T) {
	db := setupTestDB(t)
	authService := auth.NewAuth(db)

	user, err := authService.CreateUser("test@example.com", "testuser", "testpassword123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	session, err := authService.CreateSession(user.ID)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if session.ID == "" {
		t.Fatal("Expected session ID to be set")
	}

	if session.UserID != user.ID {
		t.Fatalf("Expected user ID %d, got %d", user.ID, session.UserID)
	}

	if session.ExpiresAt.Before(time.Now()) {
		t.Fatal("Expected session to expire in the future")
	}
}

func TestGetSession(t *testing.T) {
	db := setupTestDB(t)
	authService := auth.NewAuth(db)

	user, err := authService.CreateUser("test@example.com", "testuser", "testpassword123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	createdSession, err := authService.CreateSession(user.ID)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	retrievedSession, err := authService.GetSession(createdSession.ID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if retrievedSession.ID != createdSession.ID {
		t.Fatalf("Expected session ID %s, got %s", createdSession.ID, retrievedSession.ID)
	}

	if retrievedSession.UserID != user.ID {
		t.Fatalf("Expected user ID %d, got %d", user.ID, retrievedSession.UserID)
	}
}

func TestEmailExists(t *testing.T) {
	db := setupTestDB(t)
	authService := auth.NewAuth(db)

	email := "test@example.com"

	exists, err := authService.EmailExists(email)
	if err != nil {
		t.Fatalf("Failed to check email existence: %v", err)
	}

	if exists {
		t.Fatal("Expected email to not exist")
	}

	_, err = authService.CreateUser(email, "testuser", "testpassword123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	exists, err = authService.EmailExists(email)
	if err != nil {
		t.Fatalf("Failed to check email existence: %v", err)
	}

	if !exists {
		t.Fatal("Expected email to exist")
	}
}
