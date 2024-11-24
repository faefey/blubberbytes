package database

import (
	"database/sql"
	"fmt"
)

// Join on file hash with Storing table
type Hosting struct {
	Hash  string  `json:"hash"`
	Price float64 `json:"price"`
}

// AddHosting inserts a new record into the Hosting table.
func AddHosting(db *sql.DB, hash string, price float64) error {
	query := `INSERT INTO Hosting (hash, price) VALUES (?, ?)`
	_, err := db.Exec(query, hash, price)
	if err != nil {
		return fmt.Errorf("error adding record to Hosting: %v", err)
	}

	fmt.Printf("Record added to Hosting with hash: %s\n", hash)
	return nil
}

// DeleteHosting removes a record from the Hosting table by its hash.
func DeleteHosting(db *sql.DB, hash string) error {
	query := `DELETE FROM Hosting WHERE hash = ?`
	_, err := db.Exec(query, hash)
	if err != nil {
		return fmt.Errorf("error deleting record from Hosting with hash %s: %v", hash, err)
	}

	fmt.Printf("Record with hash %s deleted successfully from Hosting.\n", hash)
	return nil
}

// FindHosting retrieves a record from the Hosting table by its hash.
func FindHosting(db *sql.DB, hash string) (*Hosting, error) {
	var hosting Hosting
	query := `SELECT hash, price FROM Hosting WHERE hash = ?`
	err := db.QueryRow(query, hash).Scan(&hosting.Hash, &hosting.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("error finding record in Hosting with hash %s: %v", hash, err)
	}

	return &hosting, nil
}
// GetAllHosting retrieves all records from the Hosting table.
func GetAllHosting(db *sql.DB) ([]Hosting, error) {
	query := `SELECT hash, price FROM Hosting`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying Hosting table: %v", err)
	}
	defer rows.Close()

	var hostingRecords []Hosting
	for rows.Next() {
		var record Hosting
		err := rows.Scan(&record.Hash, &record.Price)
		if err != nil {
			return nil, fmt.Errorf("error scanning Hosting record: %v", err)
		}
		hostingRecords = append(hostingRecords, record)
	}

	return hostingRecords, nil
}
