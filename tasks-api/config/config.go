package config

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

var DB *sql.DB

// initializeDB initializes a database connection and runs migrations
func initializeDB(db *sql.DB) error {
	// Create goose version tracking table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS goose_db_version (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version_id INTEGER NOT NULL,
			is_applied INTEGER NOT NULL,
			tstamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create version tracking table: %v", err)
	}

	// Insert initial version record if it doesn't exist
	if _, err := db.Exec(`
		INSERT OR IGNORE INTO goose_db_version (version_id, is_applied)
		VALUES (0, 1)
	`); err != nil {
		return fmt.Errorf("failed to insert initial version record: %v", err)
	}

	// Get the absolute path to the migrations directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get caller information")
	}

	// Get the directory of this file
	configDir := filepath.Dir(filename)
	// Get the parent directory (tasks directory)
	tasksDir := filepath.Dir(configDir)
	// Construct the migrations path
	migrationsPath := filepath.Join(tasksDir, "migrations")

	// Run migrations
	if err := goose.Up(db, migrationsPath); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	return nil
}

// InitDB initializes the database connection and runs migrations
func InitDB() error {
	var err error

	// Check if database file exists
	dbPath := "./task.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Create empty database file
		if _, err := os.Create(dbPath); err != nil {
			return fmt.Errorf("failed to create database file: %v", err)
		}
	}

	// Open the database
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	return initializeDB(DB)
}

// InitTestDB initializes an in-memory database for testing
func InitTestDB() *sql.DB {
	testDB, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}

	if err := initializeDB(testDB); err != nil {
		panic(err)
	}

	return testDB
}
