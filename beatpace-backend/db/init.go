package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

// InitDatabase initializes the database and runs migrations
func InitDatabase() (*sql.DB, error) {
	// First connect to the database
	db, err := InitFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Read and execute migrations
	migrationSQL, err := os.ReadFile(filepath.Join("db", "migrations.sql"))
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations file: %v", err)
	}

	// Execute migrations
	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		return nil, fmt.Errorf("failed to execute migrations: %v", err)
	}

	return db, nil
} 