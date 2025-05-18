package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool" // Using pgxpool for connection pooling
	"github.com/joho/godotenv"       // Package to load .env files
)

func main() {
	// Load environment variables from .env file
	// This should be one of the first things in your main function.
	// It will not override environment variables that are already set.
	err := godotenv.Load()
	if err != nil {
		// Log a warning, but don't make it fatal.
		// This allows the app to still run if .env is missing,
		// relying on OS-level environment variables.
		log.Println("Warning: Error loading .env file:", err)
	}


	dbUser := os.Getenv("user")
	dbPassword := os.Getenv("password")
	dbHost := os.Getenv("host")
	dbPort := os.Getenv("port")
	dbName := os.Getenv("dbname")

	// Check if all required environment variables are set
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Error: One or more database environment variables (DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME) are not set. Ensure they are in your .env file or OS environment.")
	}

	databaseURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Create a new connection pool.
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\nRaw connection string used (password redacted): %s\n", err, fmt.Sprintf("postgresql://%s:PASSWORD_REDACTED@%s:%s/%s", dbUser, dbHost, dbPort, dbName))
	}
	// Defer closing the pool to ensure resources are released when main exits.
	defer pool.Close()

	// Ping the database to verify the connection.
	// Acquiring a connection from the pool and then pinging it.
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Unable to acquire a connection from the pool: %v\n", err)
	}
	defer conn.Release() // Release the connection back to the pool

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("Error pinging the database: %v\nMake sure your IP is whitelisted in Supabase network restrictions if you have any, and that your .env file is correctly configured.\n", err)
	}

	fmt.Println("Successfully connected to Supabase database using pgx/v5!")

	// Example query to test the connection further and get PostgreSQL version
	var version string
	err = conn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}

	log.Println("PostgreSQL version:", version)
}
