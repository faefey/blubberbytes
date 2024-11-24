package database

import (
	"database/sql"
	"fmt"
)

// Join on file hash with Storing table
type Sharing struct {
	Hash     string `json:"hash"`
	Password string `json:"password"`
}

// AddSharing inserts a new record into the Sharing table.
func AddSharing(db *sql.DB, hash, password string) error {
	query := `INSERT INTO Sharing (hash, password) VALUES (?, ?)`
	_, err := db.Exec(query, hash, password)
	if err != nil {
		return fmt.Errorf("error adding record to Sharing: %v", err)
	}

	fmt.Printf("Record added to Sharing with hash: %s\n", hash)
	return nil
}

// DeleteSharing removes a record from the Sharing table by its hash.
func DeleteSharing(db *sql.DB, hash string) error {
	query := `DELETE FROM Sharing WHERE hash = ?`
	_, err := db.Exec(query, hash)
	if err != nil {
		return fmt.Errorf("error deleting record from Sharing with hash %s: %v", hash, err)
	}

	fmt.Printf("Record with hash %s deleted successfully from Sharing.\n", hash)
	return nil
}

// FindSharing retrieves a record from the Sharing table by its hash.
func FindSharing(db *sql.DB, hash string) (*Sharing, error) {
	var sharing Sharing
	query := `SELECT hash, password FROM Sharing WHERE hash = ?`
	err := db.QueryRow(query, hash).Scan(&sharing.Hash, &sharing.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("error finding record in Sharing with hash %s: %v", hash, err)
	}

	return &sharing, nil
}
