package config

import (
	"database/sql"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./task_management.db")
	if err != nil {
		return err
	}

	// Run migrations
	if err := goose.Up(DB, "./migrations"); err != nil {
		return err
	}

	return nil
}

// InitTestDB initializes an in-memory database for testing
func InitTestDB() *sql.DB {
	testDB, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}

	// Create the goose version table first and initialize it
	if _, err := testDB.Exec(`
		CREATE TABLE IF NOT EXISTS goose_db_version (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version_id INTEGER NOT NULL,
			is_applied INTEGER NOT NULL,
			tstamp TIMESTAMP DEFAULT (datetime('now'))
		);
	`); err != nil {
		panic(err)
	}

	// Insert initial version record
	if _, err := testDB.Exec(`
		INSERT INTO goose_db_version (version_id, is_applied)
		VALUES (0, 1)
	`); err != nil {
		panic(err)
	}

	// Get the absolute path to the migrations directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get caller information")
	}

	// Get the directory of this file
	configDir := filepath.Dir(filename)
	// Get the parent directory (tasks directory)
	tasksDir := filepath.Dir(configDir)
	// Construct the migrations path
	migrationsPath := filepath.Join(tasksDir, "migrations")

	// Run migrations
	if err := goose.Up(testDB, migrationsPath); err != nil {
		panic(err)
	}

	return testDB
}
