package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"test-task-scout-go/internal/config"
	"test-task-scout-go/internal/repository"
	"test-task-scout-go/internal/router"
	"test-task-scout-go/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	quoteRepo, repoCloser, err := initRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	quoteService := service.NewQuoteService(quoteRepo)

	httpHandler := router.NewRouter(quoteService)

	server := startServer(cfg.Port, httpHandler, cfg.RepositoryType, cfg.DatabasePath)

	shutdownServer(server, repoCloser)

	log.Println("Application stopped.")
}

func initRepository(cfg *config.Config) (repository.QuoteRepository, func() error, error) {
	var quoteRepo repository.QuoteRepository
	var repoCloser func() error

	switch cfg.RepositoryType {
	case "inmemory":
		log.Println("Using In-Memory Repository")
		quoteRepo = repository.NewInMemoryRepository()
		repoCloser = func() error { return nil }
	case "sqlite":
		log.Println("Using SQLite Repository")
		sqliteRepo, err := repository.NewSQLiteRepository(cfg.DatabasePath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to initialize SQLite repository: %w", err)
		}
		quoteRepo = sqliteRepo
		repoCloser = sqliteRepo.Close
	default:
		return nil, nil, fmt.Errorf("unknown repository type: %s", cfg.RepositoryType)
	}

	return quoteRepo, repoCloser, nil
}

func startServer(port string, handler http.Handler, repoType, dbPath string) *http.Server {
	addr := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Starting server on %s", addr)
		log.Printf("Repository type: %s", repoType)
		if repoType == "sqlite" {
			log.Printf("Database path: %s", dbPath)
		}

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	return server
}

func shutdownServer(server *http.Server, repoCloser func() error) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server graceful shutdown failed: %v", err)
	}

	log.Println("HTTP server stopped.")

	if repoCloser != nil {
		log.Println("Closing repository...")
		if err := repoCloser(); err != nil {
			log.Printf("Error closing repository: %v", err)
		} else {
			log.Println("Repository closed.")
		}
	}
} 