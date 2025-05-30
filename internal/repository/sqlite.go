package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"test-task-scout-go/internal/domain"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	repo := &SQLiteRepository{db: db}

	if err = repo.initDB(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return repo, nil
}

func (r *SQLiteRepository) initDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS quotes (
		id TEXT PRIMARY KEY,
		text TEXT NOT NULL,
		author TEXT NOT NULL
	);`
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

func (r *SQLiteRepository) Create(quote *domain.Quote) error {
	query := "INSERT INTO quotes (id, text, author) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, quote.ID, quote.Text, quote.Author)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: quotes.id" {
			return errors.New("quote with this ID already exists")
		}
		return fmt.Errorf("failed to create quote: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) GetAll() ([]domain.Quote, error) {
	query := "SELECT id, text, author FROM quotes"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all quotes: %w", err)
	}
	defer rows.Close()

	var quotes []domain.Quote
	for rows.Next() {
		var quote domain.Quote
		if err := rows.Scan(&quote.ID, &quote.Text, &quote.Author); err != nil {
			return nil, fmt.Errorf("failed to scan quote row: %w", err)
		}
		quotes = append(quotes, quote)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return quotes, nil
}

func (r *SQLiteRepository) GetByID(id string) (*domain.Quote, error) {
	query := "SELECT id, text, author FROM quotes WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var quote domain.Quote
	err := row.Scan(&quote.ID, &quote.Text, &quote.Author)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("quote not found")
		}
		return nil, fmt.Errorf("failed to get quote by ID: %w", err)
	}

	return &quote, nil
}

func (r *SQLiteRepository) GetByAuthor(author string) ([]domain.Quote, error) {
	query := "SELECT id, text, author FROM quotes WHERE author = ?"
	rows, err := r.db.Query(query, author)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes by author: %w", err)
	}
	defer rows.Close()

	var quotes []domain.Quote
	for rows.Next() {
		var quote domain.Quote
		if err := rows.Scan(&quote.ID, &quote.Text, &quote.Author); err != nil {
			return nil, fmt.Errorf("failed to scan quote row: %w", err)
		}
		quotes = append(quotes, quote)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return quotes, nil
}

func (r *SQLiteRepository) Delete(id string) error {
	query := "DELETE FROM quotes WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete quote: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("quote not found")
	}

	return nil
}

func (r *SQLiteRepository) GetRandom() (*domain.Quote, error) {
	query := "SELECT id, text, author FROM quotes ORDER BY RANDOM() LIMIT 1"
	row := r.db.QueryRow(query)

	var quote domain.Quote
	err := row.Scan(&quote.ID, &quote.Text, &quote.Author)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no quotes available")
		}
		return nil, fmt.Errorf("failed to get random quote: %w", err)
	}

	return &quote, nil
}
