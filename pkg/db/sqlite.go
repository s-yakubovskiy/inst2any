package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CheckAndInsert(id, table string, db *sql.DB) (bool, error) {
	var synced bool

	// check if table name is valid
	switch table {
	case "media", "stories":
		// valid table name, do nothing
	default:
		return false, fmt.Errorf("Invalid table name: %s", table)
	}

	query := fmt.Sprintf("SELECT synced FROM %s WHERE id = ?", table)

	err := db.QueryRow(query, id).Scan(&synced)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if err == sql.ErrNoRows {
		insertQuery := fmt.Sprintf("INSERT INTO %s (id, synced) VALUES (?, 0)", table)
		_, err = db.Exec(insertQuery, id)
		if err != nil {
			return false, err
		}
	}

	return synced, nil
}

func MarkAsSynced(id, table string, db *sql.DB) error {
	// check if table name is valid
	switch table {
	case "media", "stories":
		// valid table name, do nothing
	default:
		return fmt.Errorf("Invalid table name: %s", table)
	}

	// If media is successfully uploaded, update the media record as synced in the database
	query := fmt.Sprintf("UPDATE %s SET synced = 1 WHERE id = ?", table)
	_, err := db.Exec(query, id)
	if err != nil {
		log.Printf("Failed to update media as synced: %v", err)
		return err
	}
	return nil

}

func SetupDB(database string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil, err
	}

	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS media (id TEXT, synced BOOLEAN DEFAULT 0)")
	_, err = statement.Exec()
	if err != nil {
		return nil, err
	}
	statement, _ = db.Prepare("CREATE TABLE IF NOT EXISTS stories (id TEXT, synced BOOLEAN DEFAULT 0)")
	_, err = statement.Exec()
	if err != nil {
		return nil, err
	}

	return db, nil
}
