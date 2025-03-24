package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/yourusername/pkb/internal/api"
	"github.com/yourusername/pkb/internal/db"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Parse command line flags
	var (
		dbPath   = flag.String("db", "", "Path to SQLite database file")
		port     = flag.String("port", "8080", "Port to run the server on")
		dataDir  = flag.String("data", "", "Directory to store data files")
	)
	flag.Parse()

	// Set default paths if not provided
	if *dbPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to get user home directory: %v", err)
		}
		defaultDataDir := filepath.Join(homeDir, ".pkb")
		*dbPath = filepath.Join(defaultDataDir, "pkb.db")
		
		if *dataDir == "" {
			*dataDir = defaultDataDir
		}
	}

	// Ensure data directory exists
	if *dataDir != "" {
		if err := os.MkdirAll(*dataDir, 0755); err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}
	}

	// Initialize database
	database, err := db.New(*dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize and run API server
	server, err := api.NewServer(database, *dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	log.Printf("Starting server on port %s", *port)
	if err := server.Run(":" + *port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
