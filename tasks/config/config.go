package config

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./task_management.db")
	if err != nil {
		return err
	}

	// Create users table
	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err = DB.Exec(createUsersTable)
	return err
}
