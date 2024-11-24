package database

import (
	"database/sql"
	"fmt"
)

// Table for saved files
type Saved struct {
	Hash      string `json:"hash"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
}

// AddSaved inserts a new record into the Saved table.
func AddSaved(db *sql.DB, hash, name, extension string, size int64) error {
	query := `INSERT INTO Saved (hash, name, extension, size) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, hash, name, extension, size)
	if err != nil {
		return fmt.Errorf("error adding record to Saved: %v", err)
	}

	fmt.Printf("Record added to Saved with hash: %s\n", hash)
	return nil
}

// DeleteSaved removes a record from the Saved table by its hash.
func DeleteSaved(db *sql.DB, hash string) error {
	query := `DELETE FROM Saved WHERE hash = ?`
	_, err := db.Exec(query, hash)
	if err != nil {
		return fmt.Errorf("error deleting record from Saved with hash %s: %v", hash, err)
	}

	fmt.Printf("Record with hash %s deleted successfully from Saved.\n", hash)
	return nil
}

// FindSaved retrieves a record from the Saved table by its hash.
func FindSaved(db *sql.DB, hash string) (*Saved, error) {
	var saved Saved
	query := `SELECT hash, name, extension, size FROM Saved WHERE hash = ?`
	err := db.QueryRow(query, hash).Scan(&saved.Hash, &saved.Name, &saved.Extension, &saved.Size)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("error finding record in Saved with hash %s: %v", hash, err)
	}

	return &saved, nil
}
