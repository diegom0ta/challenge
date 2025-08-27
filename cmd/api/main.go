package main

import (
	"log"
	"os"

	"challenge/api"
)

func main() {
	port := getEnv("PORT", "8080")

	server := api.NewServer(port)

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
