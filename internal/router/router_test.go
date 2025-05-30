package router_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"test-task-scout-go/internal/domain"
	"test-task-scout-go/internal/router"
	"strings"
	"testing"
)

type MockQuoteService struct {
	CreateQuoteFunc   func(text, author string) (*domain.Quote, error)
	GetAllQuotesFunc  func(authorFilter string) ([]domain.Quote, error)
	GetRandomQuoteFunc func() (*domain.Quote, error)
	DeleteQuoteFunc   func(id string) error
}

func (m *MockQuoteService) CreateQuote(text, author string) (*domain.Quote, error) {
	return m.CreateQuoteFunc(text, author)
}
func (m *MockQuoteService) GetAllQuotes(authorFilter string) ([]domain.Quote, error) {
	return m.GetAllQuotesFunc(authorFilter)
}
func (m *MockQuoteService) GetRandomQuote() (*domain.Quote, error) {
	return m.GetRandomQuoteFunc()
}
func (m *MockQuoteService) DeleteQuote(id string) error {
	return m.DeleteQuoteFunc(id)
}

func TestRouter_CreateQuote(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockServiceSuccess := &MockQuoteService{
			CreateQuoteFunc: func(text, author string) (*domain.Quote, error) {
				return &domain.Quote{ID: "new-id", Text: text, Author: author}, nil
			},
		}
		r := router.NewRouter(mockServiceSuccess)

		quoteData := `{"text": "New Quote", "author": "New Author"}`
		req, _ := http.NewRequest("POST", "/quotes", bytes.NewBufferString(quoteData))
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		var createdQuote domain.Quote
		if err := json.NewDecoder(rr.Body).Decode(&createdQuote); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}
		if createdQuote.Text != "New Quote" || createdQuote.Author != "New Author" || createdQuote.ID != "new-id" {
			t.Errorf("handler returned unexpected body: got %+v", createdQuote)
		}
	})

	t.Run("EmptyBody", func(t *testing.T) {
		t.Parallel()
		mockService := &MockQuoteService{}
		r := router.NewRouter(mockService)

		req, _ := http.NewRequest("POST", "/quotes", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code for empty body: got %v want %v",
				status, http.StatusBadRequest)
		}
		if !strings.Contains(rr.Body.String(), "Request body is empty") {
			t.Errorf("handler returned unexpected body for empty body: got %v", rr.Body.String())
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		t.Parallel()
		mockService := &MockQuoteService{}
		r := router.NewRouter(mockService)

		invalidJsonData := `{"text": "New Quote", "author": "New Author"`
		req, _ := http.NewRequest("POST", "/quotes", bytes.NewBufferString(invalidJsonData))
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code for invalid JSON: got %v want %v",
				status, http.StatusBadRequest)
		}
		if !strings.Contains(rr.Body.String(), "Invalid request body") {
			t.Errorf("handler returned unexpected body for invalid JSON: got %v", rr.Body.String())
		}
	})

	t.Run("ServiceValidationError", func(t *testing.T) {
		t.Parallel()
		mockServiceValidationError := &MockQuoteService{
			CreateQuoteFunc: func(text, author string) (*domain.Quote, error) {
				return nil, errors.New("text and author cannot be empty")
			},
		}
		rValidationError := router.NewRouter(mockServiceValidationError)

		quoteDataEmpty := `{"text": "", "author": ""}`
		req, _ := http.NewRequest("POST", "/quotes", bytes.NewBufferString(quoteDataEmpty))
		rr := httptest.NewRecorder()
		rValidationError.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code for service validation error: got %v want %v",
				status, http.StatusBadRequest)
		}
		if !strings.Contains(rr.Body.String(), "text and author cannot be empty") {
			t.Errorf("handler returned unexpected body for service validation error: got %v", rr.Body.String())
		}
	})

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		mockServiceError := &MockQuoteService{
			CreateQuoteFunc: func(text, author string) (*domain.Quote, error) {
				return nil, errors.New("some service error")
			},
		}
		rServiceError := router.NewRouter(mockServiceError)

		quoteData := `{"text": "New Quote", "author": "New Author"}`
		req, _ := http.NewRequest("POST", "/quotes", bytes.NewBufferString(quoteData))
		rr := httptest.NewRecorder()
		rServiceError.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code for other service error: got %v want %v",
				status, http.StatusInternalServerError)
		}
		if !strings.Contains(rr.Body.String(), "Failed to create quote") {
			t.Errorf("handler returned unexpected body for other service error: got %v", rr.Body.String())
		}
	})
}

func TestRouter_GetAllQuotes(t *testing.T) {
	t.Parallel()
	quotes := []domain.Quote{
		{ID: "1", Text: "Q1", Author: "A1"},
		{ID: "2", Text: "Q2", Author: "A2"},
	}
	mockServiceSuccess := &MockQuoteService{
		GetAllQuotesFunc: func(authorFilter string) ([]domain.Quote, error) {
			if authorFilter == "" {
				return quotes, nil
			}
			filtered := []domain.Quote{}
			for _, q := range quotes {
				if q.Author == authorFilter {
					filtered = append(filtered, q)
				}
			}
			return filtered, nil
		},
	}
	r := router.NewRouter(mockServiceSuccess)

	req, _ := http.NewRequest("GET", "/quotes", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for GetAll: got %v want %v",
			status, http.StatusOK)
	}

	var returnedQuotes []domain.Quote
	if err := json.NewDecoder(rr.Body).Decode(&returnedQuotes); err != nil {
		t.Fatalf("Failed to decode response body for GetAll: %v", err)
	}
	if len(returnedQuotes) != 2 {
		t.Errorf("Expected 2 quotes from GetAll, got %d", len(returnedQuotes))
	}

	req, _ = http.NewRequest("GET", "/quotes?author=A1", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for GetByAuthor: got %v want %v",
			status, http.StatusOK)
	}

	if err := json.NewDecoder(rr.Body).Decode(&returnedQuotes); err != nil {
		t.Fatalf("Failed to decode response body for GetByAuthor: %v", err)
	}
	if len(returnedQuotes) != 1 || returnedQuotes[0].ID != "1" {
		t.Errorf("Expected 1 quote from GetByAuthor, got %+v", returnedQuotes)
	}

	req, _ = http.NewRequest("GET", "/quotes?author=NonExistent", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for GetByAuthor (non-existent): got %v want %v",
			status, http.StatusOK)
	}

	if err := json.NewDecoder(rr.Body).Decode(&returnedQuotes); err != nil {
		t.Fatalf("Failed to decode response body for GetByAuthor (non-existent): %v", err)
	}
	if len(returnedQuotes) != 0 {
		t.Errorf("Expected 0 quotes from GetByAuthor (non-existent), got %d", len(returnedQuotes))
	}

	req, _ = http.NewRequest("POST", "/quotes", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	mockServiceError := &MockQuoteService{
		GetAllQuotesFunc: func(authorFilter string) ([]domain.Quote, error) {
			return nil, errors.New("service error")
		},
	}
	rServiceError := router.NewRouter(mockServiceError)

	req, _ = http.NewRequest("GET", "/quotes", nil)
	rr = httptest.NewRecorder()
	rServiceError.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code for service error: got %v want %v",
			status, http.StatusInternalServerError)
	}
	if !strings.Contains(rr.Body.String(), "Failed to retrieve quotes") {
		t.Errorf("handler returned unexpected body for service error: got %v", rr.Body.String())
	}
}

func TestRouter_GetRandomQuote(t *testing.T) {
	t.Parallel()
	mockServiceSuccess := &MockQuoteService{
		GetRandomQuoteFunc: func() (*domain.Quote, error) {
			return &domain.Quote{ID: "random", Text: "Random Quote", Author: "Random Author"}, nil
		},
	}
	r := router.NewRouter(mockServiceSuccess)

	req, _ := http.NewRequest("GET", "/quotes/random", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var returnedQuote domain.Quote
	if err := json.NewDecoder(rr.Body).Decode(&returnedQuote); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	if returnedQuote.ID != "random" || returnedQuote.Text != "Random Quote" || returnedQuote.Author != "Random Author" {
		t.Errorf("handler returned unexpected body: got %+v", returnedQuote)
	}

	req, _ = http.NewRequest("POST", "/quotes/random", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code for wrong method: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	mockServiceNoQuotes := &MockQuoteService{
		GetRandomQuoteFunc: func() (*domain.Quote, error) {
			return nil, errors.New("no quotes available")
		},
	}
	rNoQuotes := router.NewRouter(mockServiceNoQuotes)

	req, _ = http.NewRequest("GET", "/quotes/random", nil)
	rr = httptest.NewRecorder()
	rNoQuotes.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for no quotes: got %v want %v",
			status, http.StatusNotFound)
	}
	if !strings.Contains(rr.Body.String(), "no quotes available") {
		t.Errorf("handler returned unexpected body for no quotes: got %v", rr.Body.String())
	}

	mockServiceError := &MockQuoteService{
		GetRandomQuoteFunc: func() (*domain.Quote, error) {
			return nil, errors.New("service error")
		},
	}
	rServiceError := router.NewRouter(mockServiceError)

	req, _ = http.NewRequest("GET", "/quotes/random", nil)
	rr = httptest.NewRecorder()
	rServiceError.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code for service error: got %v want %v",
			status, http.StatusInternalServerError)
	}
	if !strings.Contains(rr.Body.String(), "Failed to retrieve random quote") {
		t.Errorf("handler returned unexpected body for service error: got %v", rr.Body.String())
	}
}

func TestRouter_DeleteQuote(t *testing.T) {
	t.Parallel()
	mockServiceSuccess := &MockQuoteService{
		DeleteQuoteFunc: func(id string) error {
			if id == "123" {
				return nil
			}
			return errors.New("quote not found")
		},
	}
	r := router.NewRouter(mockServiceSuccess)

	req, _ := http.NewRequest("DELETE", "/quotes/123", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
	if rr.Body.String() != "" {
		t.Errorf("handler returned non-empty body for 204 status: got %v", rr.Body.String())
	}

	req, _ = http.NewRequest("DELETE", "/quotes/non-existent", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent quote: got %v want %v",
			status, http.StatusNotFound)
	}
	if !strings.Contains(rr.Body.String(), "quote not found") {
		t.Errorf("handler returned unexpected body for non-existent quote: got %v", rr.Body.String())
	}

	req, _ = http.NewRequest("GET", "/quotes/123", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code for wrong method: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	req, _ = http.NewRequest("DELETE", "/quotes/123/extra", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for invalid path: got %v want %v",
			status, http.StatusNotFound)
	}

	req, _ = http.NewRequest("DELETE", "/quotes/", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for empty ID: got %v want %v",
			status, http.StatusNotFound)
	}

	mockServiceError := &MockQuoteService{
		DeleteQuoteFunc: func(id string) error {
			return errors.New("service error")
		},
	}
	rServiceError := router.NewRouter(mockServiceError)

	req, _ = http.NewRequest("DELETE", "/quotes/some-id", nil)
	rr = httptest.NewRecorder()
	rServiceError.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code for service error: got %v want %v",
			status, http.StatusInternalServerError)
	}
	if !strings.Contains(rr.Body.String(), "Failed to delete quote") {
		t.Errorf("handler returned unexpected body for service error: got %v", rr.Body.String())
	}
} 