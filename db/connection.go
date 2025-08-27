package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" 
)

type DB struct {
	conn *sql.DB
}

func NewConnection() (*DB, error) {
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "postgres")
	dbname := getEnv("POSTGRES_DB", "challenge")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	fmt.Println("Database connection established successfully")
	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

func (db *DB) GetConnection() *sql.DB {
	return db.conn
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
