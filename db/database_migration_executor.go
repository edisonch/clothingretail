package db

import (
	"database/sql"
	"io/ioutil"
	"log"
)

func RunMigration(db *sql.DB, migrationFile string) error {
	sqlContent, err := ioutil.ReadFile(migrationFile)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(sqlContent))
	if err != nil {
		return err
	}

	log.Printf("Migration %s executed successfully", migrationFile)
	return nil
}
