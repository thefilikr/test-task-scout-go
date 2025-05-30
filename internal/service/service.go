package service

import "test-task-scout-go/internal/domain"

type QuoteService interface {
	CreateQuote(text, author string) (*domain.Quote, error)
	GetAllQuotes(authorFilter string) ([]domain.Quote, error)
	GetRandomQuote() (*domain.Quote, error)
	DeleteQuote(id string) error
	GetByID(id string) (*domain.Quote, error)
} 