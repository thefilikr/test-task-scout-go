package repository

import "test-task-scout-go/internal/domain"

type QuoteRepository interface {
	Create(quote *domain.Quote) error
	GetAll() ([]domain.Quote, error)
	GetByID(id string) (*domain.Quote, error)
	GetByAuthor(author string) ([]domain.Quote, error)
	Delete(id string) error
	GetRandom() (*domain.Quote, error)
} 