package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Add an argument worker to specify which worker you are checking (VK or TG)
func CheckAndInsert(id, table, worker string, db *sql.DB) (bool, error) {
	var synced bool
	syncColumn := "synced_vk"
	if worker == "TG" {
		syncColumn = "synced_tg"
	}

	switch table {
	case "media", "stories":
	default:
		return false, fmt.Errorf("Invalid table name: %s", table)
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = ?", syncColumn, table)
	err := db.QueryRow(query, id).Scan(&synced)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if err == sql.ErrNoRows {
		insertQuery := fmt.Sprintf("INSERT INTO %s (id, %s) VALUES (?, 0)", table, syncColumn)
		_, err = db.Exec(insertQuery, id)
		if err != nil {
			return false, err
		}
	}

	return synced, nil
}

// Add an argument worker to specify which worker you are marking as synced
func MarkAsSynced(id, table, worker string, db *sql.DB) error {
	syncColumn := "synced_vk"
	if worker == "TG" {
		syncColumn = "synced_tg"
	}

	switch table {
	case "media", "stories":
	default:
		return fmt.Errorf("Invalid table name: %s", table)
	}

	query := fmt.Sprintf("UPDATE %s SET %s = 1 WHERE id = ?", table, syncColumn)
	_, err := db.Exec(query, id)
	if err != nil {
		log.Printf("Failed to update %s as synced: %v", table, err)
		return err
	}
	return nil
}

func SetupDB(database string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil, err
	}

	mediaTableSQL := "CREATE TABLE IF NOT EXISTS media (id TEXT, synced_vk BOOLEAN DEFAULT 0, synced_tg BOOLEAN DEFAULT 0)"
	statement, _ := db.Prepare(mediaTableSQL)
	_, err = statement.Exec()
	if err != nil {
		return nil, err
	}

	storiesTableSQL := "CREATE TABLE IF NOT EXISTS stories (id TEXT, synced_vk BOOLEAN DEFAULT 0, synced_tg BOOLEAN DEFAULT 0)"
	statement, _ = db.Prepare(storiesTableSQL)
	_, err = statement.Exec()
	if err != nil {
		return nil, err
	}

	return db, nil
}
