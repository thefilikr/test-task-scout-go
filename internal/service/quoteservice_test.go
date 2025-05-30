package service_test

import (
	"errors"
	"test-task-scout-go/internal/domain"
	"test-task-scout-go/internal/service"
	"test-task-scout-go/internal/repository"
	"testing"
)

type MockQuoteRepository struct {
	CreateFunc      func(quote *domain.Quote) error
	GetAllFunc      func() ([]domain.Quote, error)
	GetByIDFunc     func(id string) (*domain.Quote, error)
	GetByAuthorFunc func(author string) ([]domain.Quote, error)
	DeleteFunc      func(id string) error
	GetRandomFunc   func() (*domain.Quote, error)
}

func (m *MockQuoteRepository) Create(quote *domain.Quote) error {
	return m.CreateFunc(quote)
}
func (m *MockQuoteRepository) GetAll() ([]domain.Quote, error) {
	return m.GetAllFunc()
}
func (m *MockQuoteRepository) GetByID(id string) (*domain.Quote, error) {
	return m.GetByIDFunc(id)
}
func (m *MockQuoteRepository) GetByAuthor(author string) ([]domain.Quote, error) {
	return m.GetByAuthorFunc(author)
}
func (m *MockQuoteRepository) Delete(id string) error {
	return m.DeleteFunc(id)
}
func (m *MockQuoteRepository) GetRandom() (*domain.Quote, error) {
	return m.GetRandomFunc()
}

func TestQuoteService_CreateQuote(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockQuoteRepository{
			CreateFunc: func(quote *domain.Quote) error {
				if quote.Text == "" || quote.Author == "" {
					return errors.New("validation error")
				}
				quote.ID = "generated-id"
				return nil
			},
		}
		quoteService := service.NewQuoteService(mockRepo)

		quote, err := quoteService.CreateQuote("Test Text", "Test Author")
		if err != nil {
			t.Fatalf("CreateQuote failed: %v", err)
		}
		if quote == nil {
			t.Fatal("CreateQuote returned nil quote")
		}
		if quote.ID == "" {
			t.Error("Created quote has empty ID")
		}
		if quote.Text != "Test Text" || quote.Author != "Test Author" {
			t.Errorf("Created quote has wrong text or author: %+v", quote)
		}
	})

	t.Run("ValidationError", func(t *testing.T) {
		mockRepo := &MockQuoteRepository{
			CreateFunc: func(quote *domain.Quote) error {
				if quote.Text == "" || quote.Author == "" {
					return errors.New("validation error")
				}
				return nil
			},
		}
		quoteService := service.NewQuoteService(mockRepo)

		_, err := quoteService.CreateQuote("", "Test Author")
		if err == nil {
			t.Error("CreateQuote did not return error for empty text")
		}
		if err != nil && err.Error() != "text and author cannot be empty" {
			t.Errorf("Expected validation error, got: %v", err)
		}

		_, err = quoteService.CreateQuote("Test Text", "")
		if err == nil {
			t.Error("CreateQuote did not return error for empty author")
		}
		if err != nil && err.Error() != "text and author cannot be empty" {
			t.Errorf("Expected validation error, got: %v", err)
		}
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := &MockQuoteRepository{
			CreateFunc: func(quote *domain.Quote) error {
				return errors.New("database error")
			},
		}
		quoteService := service.NewQuoteService(mockRepo)

		_, err := quoteService.CreateQuote("Test Text", "Test Author")
		if err == nil {
			t.Error("CreateQuote did not return repository error")
		}
		expectedErr := "failed to create quote in repository: database error"
		if err != nil && err.Error() != expectedErr {
			t.Errorf("Expected '%s' error, got: %v", expectedErr, err)
		}
	})
}

func TestQuoteService_GetAllQuotes(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expectedQuotes := []domain.Quote{
			{ID: "1", Text: "Quote 1", Author: "Author 1"},
			{ID: "2", Text: "Quote 2", Author: "Author 2"},
		}
		mockRepo := &MockQuoteRepository{
			GetAllFunc: func() ([]domain.Quote, error) {
				return expectedQuotes, nil
			},
			GetByAuthorFunc: func(author string) ([]domain.Quote, error) {
				if author == "Author 1" {
					return []domain.Quote{{ID: "1", Text: "Quote 1", Author: "Author 1"}}, nil
				}
				return nil, nil
			},
		}
		quoteService := service.NewQuoteService(mockRepo)

		quotes, err := quoteService.GetAllQuotes("")
		if err != nil {
			t.Fatalf("GetAllQuotes without filter failed: %v", err)
		}
		if len(quotes) != len(expectedQuotes) {
			t.Errorf("Expected %d quotes, got %d", len(expectedQuotes), len(quotes))
		}

		quotesByAuthor, err := quoteService.GetAllQuotes("Author 1")
		if err != nil {
			t.Fatalf("GetAllQuotes with filter failed: %v", err)
		}
		if len(quotesByAuthor) != 1 {
			t.Errorf("Expected 1 quote for author 'Author 1', got %d", len(quotesByAuthor))
		}
		if len(quotesByAuthor) > 0 && (quotesByAuthor[0].Text != "Quote 1" || quotesByAuthor[0].Author != "Author 1") {
			t.Errorf("Unexpected quote returned for author filter: %+v", quotesByAuthor[0])
		}

		quotesByNonExistentAuthor, err := quoteService.GetAllQuotes("NonExistent Author")
		if err != nil {
			t.Fatalf("GetAllQuotes with non-existent author filter failed: %v", err)
		}
		if len(quotesByNonExistentAuthor) != 0 {
			t.Errorf("Expected 0 quotes for non-existent author, got %d", len(quotesByNonExistentAuthor))
		}
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := &MockQuoteRepository{
			GetAllFunc: func() ([]domain.Quote, error) {
				return nil, errors.New("database error")
			},
			GetByAuthorFunc: func(author string) ([]domain.Quote, error) {
				return nil, errors.New("database error")
			},
		}
		quoteService := service.NewQuoteService(mockRepo)

		_, err := quoteService.GetAllQuotes("")
		if err == nil {
			t.Error("GetAllQuotes without filter did not return repository error")
		}
		if err != nil && err.Error() != "failed to get all quotes from repository: database error" {
			t.Errorf("Expected database error, got: %v", err)
		}

		_, err = quoteService.GetAllQuotes("Some Author")
		if err == nil {
			t.Error("GetAllQuotes with filter did not return repository error")
		}
		if err != nil && err.Error() != "failed to get quotes by author from repository: database error" {
			t.Errorf("Expected database error, got: %v", err)
		}
	})
}

func TestQuoteService_GetRandomQuote(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expectedQuote := &domain.Quote{ID: "random-1", Text: "Random Quote", Author: "Random Author"}
		mockRepo := &MockQuoteRepository{
			GetRandomFunc: func() (*domain.Quote, error) {
				return expectedQuote, nil
			},
		}
		quoteService := service.NewQuoteService(mockRepo)

		quote, err := quoteService.GetRandomQuote()
		if err != nil {
			t.Fatalf("GetRandomQuote failed: %v", err)
		}
		if quote == nil {
			t.Fatal("GetRandomQuote returned nil quote")
		}
		if quote.ID != expectedQuote.ID || quote.Text != expectedQuote.Text || quote.Author != expectedQuote.Author {
			t.Errorf("GetRandomQuote returned unexpected quote: %+v", quote)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo := &MockQuoteRepository{
			GetRandomFunc: func() (*domain.Quote, error) {
				return nil, errors.New("not found")
			},
		}
		quoteService := service.NewQuoteService(mockRepo)

		quote, err := quoteService.GetRandomQuote()
		if err == nil {
			t.Error("GetRandomQuote did not return error when not found")
		}
		if err != nil && err.Error() != "not found" {
			t.Errorf("Expected 'not found' error, got: %v", err)
		}
		if quote != nil {
			t.Error("GetRandomQuote returned non-nil quote when not found")
		}
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := &MockQuoteRepository{
			GetRandomFunc: func() (*domain.Quote, error) {
				return nil, errors.New("database error")
			},
		}
		quoteService := service.NewQuoteService(mockRepo)

		_, err := quoteService.GetRandomQuote()
		if err == nil {
			t.Error("GetRandomQuote did not return repository error")
		}
		if err != nil && err.Error() != "failed to get random quote from repository: database error" {
			t.Errorf("Expected database error, got: %v", err)
		}
	})
}

func TestQuoteService_DeleteQuote(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockQuoteRepository{
			DeleteFunc: func(id string) error {
				if id == "123" {
					return nil
				}
				return errors.New("not found")
			},
		}
		quoteService := service.NewQuoteService(mockRepo)

		err := quoteService.DeleteQuote("123")
		if err != nil {
			t.Fatalf("DeleteQuote failed: %v", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo := &MockQuoteRepository{
			DeleteFunc: func(id string) error {
				return errors.New("not found")
			},
		}
		quoteService := service.NewQuoteService(mockRepo)

		err := quoteService.DeleteQuote("non-existent")
		if err == nil {
			t.Error("DeleteQuote did not return error when not found")
		}
		if err != nil && err.Error() != "not found" {
			t.Errorf("Expected 'not found' error, got: %v", err)
		}
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := &MockQuoteRepository{
			DeleteFunc: func(id string) error {
				return errors.New("database error")
			},
		}
		quoteService := service.NewQuoteService(mockRepo)

		err := quoteService.DeleteQuote("some-id")
		if err == nil {
			t.Error("DeleteQuote did not return repository error")
		}
		if err != nil && err.Error() != "failed to delete quote from repository: database error" {
			t.Errorf("Expected database error, got: %v", err)
		}
	})
} 