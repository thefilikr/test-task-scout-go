
package repository_test

import (
	"test-task-scout-go/internal/domain"
	"test-task-scout-go/internal/repository"

	"testing"
)

func testRepositoryCreateAndGet(t *testing.T, repo repository.QuoteRepository) {
	t.Helper()
	quote := &domain.Quote{
		ID:     "test-123",
		Text:   "Common Test Quote",
		Author: "Common Test Author",
	}

	err := repo.Create(quote)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	retrievedQuote, err := repo.GetByID("test-123")
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if retrievedQuote.ID != quote.ID || retrievedQuote.Text != quote.Text || retrievedQuote.Author != quote.Author {
		t.Errorf("Retrieved quote does not match created one. Expected %+v, got %+v", quote, retrievedQuote)
	}

	_, err = repo.GetByID("non-existent")
	if err == nil {
		t.Error("GetByID for non-existent ID did not return an error")
	}
	expectedErr := "quote not found"
	if err != nil && err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func testRepositoryCreateDuplicate(t *testing.T, repo repository.QuoteRepository) {
	t.Helper()

	quote1 := &domain.Quote{ID: "dup-1", Text: "Quote 1", Author: "Author 1"}
	quote2 := &domain.Quote{ID: "dup-1", Text: "Quote 2", Author: "Author 2"}

	err := repo.Create(quote1)
	if err != nil {
		t.Fatalf("Create failed for first quote: %v", err)
	}

	err = repo.Create(quote2)
	if err == nil {
		t.Error("Create did not return an error for duplicate ID")
	}
	expectedErr := "quote with this ID already exists"
	if err != nil && err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func testRepositoryGetAll(t *testing.T, repo repository.QuoteRepository) {
	t.Helper()

	emptyQuotes, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed on empty repo: %v", err)
	}
	if len(emptyQuotes) != 0 {
		t.Errorf("Expected 0 quotes from empty repo, got %d", len(emptyQuotes))
	}

	quote1 := &domain.Quote{ID: "all-1", Text: "Quote 1", Author: "Author 1"}
	quote2 := &domain.Quote{ID: "all-2", Text: "Quote 2", Author: "Author 2"}

	repo.Create(quote1)
	repo.Create(quote2)

	quotes, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(quotes) != 2 {
		t.Errorf("Expected 2 quotes, got %d", len(quotes))
	}

	found1 := false
	found2 := false
	for _, q := range quotes {
		if q.ID == "all-1" && q.Text == "Quote 1" && q.Author == "Author 1" {
			found1 = true
		}
		if q.ID == "all-2" && q.Text == "Quote 2" && q.Author == "Author 2" {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Errorf("GetAll did not return all created quotes. Found 1: %t, Found 2: %t", found1, found2)
	}
}

func testRepositoryGetByAuthor(t *testing.T, repo repository.QuoteRepository) {
	t.Helper()

	quote1 := &domain.Quote{ID: "author-1", Text: "Quote 1", Author: "Author A"}
	quote2 := &domain.Quote{ID: "author-2", Text: "Quote 2", Author: "Author B"}
	quote3 := &domain.Quote{ID: "author-3", Text: "Quote 3", Author: "Author A"}

	repo.Create(quote1)
	repo.Create(quote2)
	repo.Create(quote3)

	quotesA, err := repo.GetByAuthor("Author A")
	if err != nil {
		t.Fatalf("GetByAuthor('Author A') failed: %v", err)
	}
	if len(quotesA) != 2 {
		t.Errorf("Expected 2 quotes for 'Author A', got %d", len(quotesA))
	}
	found1 := false
	found3 := false
	for _, q := range quotesA {
		if q.ID == "author-1" {
			found1 = true
		}
		if q.ID == "author-3" {
			found3 = true
		}
	}
	if !found1 || !found3 {
		t.Errorf("GetByAuthor('Author A') did not return correct quotes")
	}

	quotesB, err := repo.GetByAuthor("Author B")
	if err != nil {
		t.Fatalf("GetByAuthor('Author B') failed: %v", err)
	}
	if len(quotesB) != 1 {
		t.Errorf("Expected 1 quote for 'Author B', got %d", len(quotesB))
	}
	if len(quotesB) > 0 && (quotesB[0].ID != "author-2" || quotesB[0].Text != "Quote 2" || quotesB[0].Author != "Author B") {
		t.Errorf("GetByAuthor('Author B') returned incorrect quote: %+v", quotesB[0])
	}

	quotesC, err := repo.GetByAuthor("Author C")
	if err != nil {
		t.Fatalf("GetByAuthor('Author C') failed: %v", err)
	}
	if len(quotesC) != 0 {
		t.Errorf("Expected 0 quotes for 'Author C', got %d", len(quotesC))
	}
}

func testRepositoryDelete(t *testing.T, repo repository.QuoteRepository) {
	t.Helper()

	quote1 := &domain.Quote{ID: "del-1", Text: "Quote 1", Author: "Author 1"}
	quote2 := &domain.Quote{ID: "del-2", Text: "Quote 2", Author: "Author 2"}

	repo.Create(quote1)
	repo.Create(quote2)

	err := repo.Delete("del-1")
	if err != nil {
		t.Fatalf("Delete failed for 'del-1': %v", err)
	}

	_, err = repo.GetByID("del-1")
	if err == nil {
		t.Error("GetByID for deleted quote 'del-1' did not return an error")
	}
	expectedErr := "quote not found"
	if err != nil && err.Error() != expectedErr {
		t.Errorf("Expected error '%s' for deleted quote, got '%s'", expectedErr, err.Error())
	}

	retrievedQuote2, err := repo.GetByID("del-2")
	if err != nil {
		t.Fatalf("GetByID failed for 'del-2' after deleting 'del-1': %v", err)
	}
	if retrievedQuote2.ID != quote2.ID {
		t.Errorf("Quote 'del-2' was unexpectedly deleted or modified")
	}

	err = repo.Delete("non-existent")
	if err == nil {
		t.Error("Delete for non-existent ID did not return an error")
	}
	if err != nil && err.Error() != expectedErr {
		t.Errorf("Expected error '%s' for non-existent delete, got '%s'", expectedErr, err.Error())
	}
}

func testRepositoryGetRandom(t *testing.T, repo repository.QuoteRepository) {
	t.Helper()

	_, err := repo.GetRandom()
	if err == nil {
		t.Error("GetRandom on empty repo did not return an error")
	}
	expectedErr := "no quotes available"
	if err != nil && err.Error() != expectedErr {
		t.Errorf("Expected error '%s' on empty repo, got '%s'", expectedErr, err.Error())
	}

	quote1 := &domain.Quote{ID: "rand-1", Text: "Random Quote 1", Author: "Author 1"}
	quote2 := &domain.Quote{ID: "rand-2", Text: "Random Quote 2", Author: "Author 2"}
	quote3 := &domain.Quote{ID: "rand-3", Text: "Random Quote 3", Author: "Author 3"}

	repo.Create(quote1)
	repo.Create(quote2)
	repo.Create(quote3)

	foundIDs := make(map[string]bool)
	for i := 0; i < 100; i++ {
		quote, err := repo.GetRandom()
		if err != nil {
			t.Fatalf("GetRandom failed after adding quotes: %v", err)
		}
		if quote == nil {
			t.Fatal("GetRandom returned nil quote after adding quotes")
		}
		if quote.ID != "rand-1" && quote.ID != "rand-2" && quote.ID != "rand-3" {
			t.Errorf("GetRandom returned quote with unexpected ID: %s", quote.ID)
		}
		foundIDs[quote.ID] = true
	}

	if len(foundIDs) < 2 {
		t.Logf("Warning: GetRandom returned only %d unique quotes in 10 tries. This might indicate an issue or just bad luck.", len(foundIDs))
	}
} 