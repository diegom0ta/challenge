package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"challenge/db"
	"challenge/service"
)

type Server struct {
	httpServer *http.Server
	db         *db.DB
	b3Service  *service.B3Service
}

func NewServer(port string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + port,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	database, err := db.NewConnection()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	s.db = database

	s.b3Service = service.NewB3Service(s.db.GetConnection())

	s.setupRoutes()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	if err := s.db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	log.Println("Server gracefully stopped")
	return nil
}

func (s *Server) setupRoutes() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.healthHandler)

	mux.HandleFunc("/api/v1/b3/aggregated", s.handleB3Aggregated)

	s.httpServer.Handler = mux
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := s.db.GetConnection().Ping(); err != nil {
		http.Error(w, "Database connection failed", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().UTC().Format(time.RFC3339))
}

func (s *Server) handleB3Aggregated(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		http.Error(w, "ticker parameter is required", http.StatusBadRequest)
		return
	}

	var startDate *time.Time
	dataInicioStr := r.URL.Query().Get("data_inicio")
	if dataInicioStr != "" {
		parsed, err := time.Parse("2006-01-02", dataInicioStr)
		if err != nil {
			http.Error(w, "Invalid date format. Use ISO-8601 format (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		startDate = &parsed
	}

	aggregatedData, err := s.b3Service.GetAggregatedData(ticker, startDate)
	if err != nil {
		log.Printf("Error getting aggregated data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(aggregatedData); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
