package router

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"test-task-scout-go/internal/service"
)

type Router struct {
	service service.QuoteService
	mux     *http.ServeMux
}

func NewRouter(service service.QuoteService) *Router {
	r := &Router{
		service: service,
		mux:     http.NewServeMux(),
	}

	r.mux.HandleFunc("/quotes", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			r.createQuoteHandler(w, req)
		case http.MethodGet:
			r.getAllQuotesHandler(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.mux.HandleFunc("/quotes/", func(w http.ResponseWriter, req *http.Request) {
		id := strings.TrimPrefix(req.URL.Path, "/quotes/")
		if id == "" {
			http.Error(w, "ID is required", http.StatusNotFound)
			return
		}

		switch req.Method {
		case http.MethodGet:
			r.getQuoteByIDHandler(w, req, id)
		case http.MethodDelete:
			r.deleteQuoteHandler(w, req, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.mux.HandleFunc("/quotes/random", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.getRandomQuoteHandler(w, req)
	})

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) createQuoteHandler(w http.ResponseWriter, req *http.Request) {
	if req.Body == nil || req.ContentLength == 0 {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	var quoteData struct {
		Text   string `json:"text"`
		Author string `json:"author"`
	}

	decoder := json.NewDecoder(req.Body)
	req.Body = http.MaxBytesReader(w, req.Body, 1048576)
	err := decoder.Decode(&quoteData)
	if err != nil {
		if err == io.EOF {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
		} else if _, ok := err.(*json.SyntaxError); ok {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		} else if _, ok := err.(*json.UnmarshalTypeError); ok {
			http.Error(w, "Invalid JSON data types", http.StatusBadRequest)
		} else if err.Error() == "http: request body too large" {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		} else {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		}
		return
	}

	quote, err := r.service.CreateQuote(quoteData.Text, quoteData.Author)
	if err != nil {
		if strings.Contains(err.Error(), "cannot be empty") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			log.Printf("Error creating quote: %v", err)
			http.Error(w, "Failed to create quote", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(quote)
}

func (r *Router) getAllQuotesHandler(w http.ResponseWriter, req *http.Request) {
	authorFilter := req.URL.Query().Get("author")

	quotes, err := r.service.GetAllQuotes(authorFilter)
	if err != nil {
		log.Printf("Error getting all quotes: %v", err)
		http.Error(w, "Failed to retrieve quotes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

func (r *Router) getQuoteByIDHandler(w http.ResponseWriter, req *http.Request, id string) {
	quote, err := r.service.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Quote not found", http.StatusNotFound)
		} else {
			log.Printf("Error getting quote by ID %s: %v", id, err)
			http.Error(w, "Failed to retrieve quote", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quote)
}

func (r *Router) getRandomQuoteHandler(w http.ResponseWriter, req *http.Request) {
	quote, err := r.service.GetRandomQuote()
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "No quotes found", http.StatusNotFound)
		} else {
			log.Printf("Error getting random quote: %v", err)
			http.Error(w, "Failed to retrieve random quote", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quote)
}

func (r *Router) deleteQuoteHandler(w http.ResponseWriter, req *http.Request, id string) {
	err := r.service.DeleteQuote(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Quote not found", http.StatusNotFound)
		} else {
			log.Printf("Error deleting quote: %v", err)
			http.Error(w, "Failed to delete quote", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
} 