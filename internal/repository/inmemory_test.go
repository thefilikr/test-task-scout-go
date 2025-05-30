package repository_test

import (
	"testing"

	"test-task-scout-go/internal/repository"
)

func TestInMemoryRepository(t *testing.T) {
	t.Run("CreateAndGet", func(t *testing.T) {
		testRepositoryCreateAndGet(t, repository.NewInMemoryRepository())
	})

	t.Run("CreateDuplicate", func(t *testing.T) {
		testRepositoryCreateDuplicate(t, repository.NewInMemoryRepository())
	})

	t.Run("GetAll", func(t *testing.T) {
		testRepositoryGetAll(t, repository.NewInMemoryRepository())
	})

	t.Run("GetByAuthor", func(t *testing.T) {
		testRepositoryGetByAuthor(t, repository.NewInMemoryRepository())
	})

	t.Run("Delete", func(t *testing.T) {
		testRepositoryDelete(t, repository.NewInMemoryRepository())
	})

	t.Run("GetRandom", func(t *testing.T) {
		testRepositoryGetRandom(t, repository.NewInMemoryRepository())
	})
} 