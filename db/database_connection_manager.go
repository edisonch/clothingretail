package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	// Use a proper migrations directory path instead of the db file path
	migrationsPath := "db/migrate-sqlite"
	// Run migrations automatically
	if err = RunMigration(DB, migrationsPath); err != nil {
		return err
	}

	return nil
}

func RunMigration(db *sql.DB, migrationsPath string) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}
	log.Println("Migrations path: ", migrationsPath)
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"sqlite", driver)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	fmt.Println("Migrations executed successfully")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
