package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(1) // SQLite works best with single connection
	DB.SetMaxIdleConns(1)

	// Test the connection
	if err = DB.Ping(); err != nil {
		return err
	}

	// Enable foreign keys
	_, err = DB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return err
	}

	log.Println("Database connection established")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
