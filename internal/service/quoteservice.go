package service

import (
	"errors"
	"fmt"
	"strconv"
	"test-task-scout-go/internal/domain"
	"test-task-scout-go/internal/repository"
	"time"
)

type QuoteServiceImpl struct {
	repo repository.QuoteRepository
}

func NewQuoteService(repo repository.QuoteRepository) *QuoteServiceImpl {
	return &QuoteServiceImpl{repo: repo}
}

func (s *QuoteServiceImpl) CreateQuote(text, author string) (*domain.Quote, error) {
	if text == "" || author == "" {
		return nil, errors.New("text and author cannot be empty")
	}

	// Простая генерация ID на основе времени. В ТЗ сказанно, не использовать сторонние библеотеки,
	// хотя я бы здесь генерировал UUID используя github.com/google/uuid
	id := strconv.FormatInt(time.Now().UnixNano(), 10)

	quote := &domain.Quote{
		ID:     id,
		Text:   text,
		Author: author,
	}

	err := s.repo.Create(quote)
	if err != nil {
		return nil, fmt.Errorf("failed to create quote in repository: %w", err)
	}

	return quote, nil
}

func (s *QuoteServiceImpl) GetAllQuotes(author string) ([]domain.Quote, error) {
	if author == "" {
		quotes, err := s.repo.GetAll()
		if err != nil {
			return nil, fmt.Errorf("failed to get all quotes from repository: %w", err)
		}
		return quotes, nil
	}

	quotes, err := s.repo.GetByAuthor(author)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes by author from repository: %w", err)
	}
	return quotes, nil
}

func (s *QuoteServiceImpl) GetRandomQuote() (*domain.Quote, error) {
	quote, err := s.repo.GetRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to get random quote from repository: %w", err)
	}
	
	return quote, nil
}

func (s *QuoteServiceImpl) DeleteQuote(id string) error {
	if id == "" {
		return errors.New("ID cannot be empty")
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete quote from repository: %w", err)
	}
	return nil	
}

func (s *QuoteServiceImpl) GetByID(id string) (*domain.Quote, error) {
	if id == "" {
		return nil, errors.New("ID cannot be empty")
	}
	quote, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote by ID from repository: %w", err)
	}
	return quote, nil
} 