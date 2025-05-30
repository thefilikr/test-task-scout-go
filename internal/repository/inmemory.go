package repository

import (
	"errors"
	"math/rand"
	"test-task-scout-go/internal/domain"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
} 

type InMemoryRepository struct {
	mu    sync.RWMutex
	quotes map[string]domain.Quote
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		quotes: make(map[string]domain.Quote),
	}
}

func (r *InMemoryRepository) Create(quote *domain.Quote) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.quotes[quote.ID]; exists {
		return errors.New("quote with this ID already exists")
	}
	r.quotes[quote.ID] = *quote
	return nil
}

func (r *InMemoryRepository) GetAll() ([]domain.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	quotes := make([]domain.Quote, 0, len(r.quotes))
	for _, quote := range r.quotes {
		quotes = append(quotes, quote)
	}
	return quotes, nil
}

func (r *InMemoryRepository) GetByID(id string) (*domain.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	quote, exists := r.quotes[id]
	if !exists {
		return nil, errors.New("quote not found")
	}
	return &quote, nil
}

func (r *InMemoryRepository) GetByAuthor(author string) ([]domain.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filteredQuotes []domain.Quote
	for _, quote := range r.quotes {
		if quote.Author == author {
			filteredQuotes = append(filteredQuotes, quote)
		}
	}
	return filteredQuotes, nil
}

func (r *InMemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.quotes[id]; !exists {
		return errors.New("quote not found")
	}
	delete(r.quotes, id)
	return nil
}

func (r *InMemoryRepository) GetRandom() (*domain.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.quotes) == 0 {
		return nil, errors.New("no quotes available")
	}
	quotes := make([]domain.Quote, 0, len(r.quotes))
	for _, quote := range r.quotes {
		quotes = append(quotes, quote)
	}
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(quotes))
	return &quotes[randomIndex], nil
} 