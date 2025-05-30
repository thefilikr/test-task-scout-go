package repository_test

import (
	"os"
	"testing"

	"test-task-scout-go/internal/repository"
)

func newTestSQLiteRepository(t *testing.T) (*repository.SQLiteRepository, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "test_quotes_*.db")
	if err != nil {
		t.Fatalf("Failed to create temporary database file: %v", err)
	}
	dbPath := tmpfile.Name()
	tmpfile.Close()

	repo, err := repository.NewSQLiteRepository(dbPath)
	if err != nil {
		os.Remove(dbPath)
		t.Fatalf("Failed to initialize SQLite repository: %v", err)
	}

	cleanup := func() {
		repo.Close()
		os.Remove(dbPath)
	}

	return repo, cleanup
}

func TestSQLiteRepository(t *testing.T) {
	t.Run("CreateAndGet", func(t *testing.T) {
		repo, cleanup := newTestSQLiteRepository(t)
		defer cleanup()
		testRepositoryCreateAndGet(t, repo)
	})

	t.Run("CreateDuplicate", func(t *testing.T) {
		repo, cleanup := newTestSQLiteRepository(t)
		defer cleanup()
		testRepositoryCreateDuplicate(t, repo)
	})

	t.Run("GetAll", func(t *testing.T) {
		repo, cleanup := newTestSQLiteRepository(t)
		defer cleanup()
		testRepositoryGetAll(t, repo)
	})

	t.Run("GetByAuthor", func(t *testing.T) {
		repo, cleanup := newTestSQLiteRepository(t)
		defer cleanup()
		testRepositoryGetByAuthor(t, repo)
	})

	t.Run("Delete", func(t *testing.T) {
		repo, cleanup := newTestSQLiteRepository(t)
		defer cleanup()
		testRepositoryDelete(t, repo)
	})

	t.Run("GetRandom", func(t *testing.T) {
		repo, cleanup := newTestSQLiteRepository(t)
		defer cleanup()
		testRepositoryGetRandom(t, repo)
	})
} 
