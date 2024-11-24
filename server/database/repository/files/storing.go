package database

import (
	"database/sql"
	"fmt"
)

// Table for added files
type Storing struct {
	Hash      string `json:"hash"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
	Path      string `json:"path"`
	Date      string `json:"date"`
}

// AddStoring inserts a new record into the Storing table.
func AddStoring(db *sql.DB, hash, name, extension, path, date string, size int64) error {
	query := `INSERT INTO Storing (hash, name, extension, size, path, date) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, hash, name, extension, size, path, date)
	if err != nil {
		return fmt.Errorf("error adding record to Storing: %v", err)
	}

	fmt.Printf("Record added to Storing with hash: %s\n", hash)
	return nil
}

// DeleteStoring removes a record from the Storing table by its hash.
func DeleteStoring(db *sql.DB, hash string) error {
	query := `DELETE FROM Storing WHERE hash = ?`
	_, err := db.Exec(query, hash)
	if err != nil {
		return fmt.Errorf("error deleting record from Storing with hash %s: %v", hash, err)
	}

	fmt.Printf("Record with hash %s deleted successfully.\n", hash)
	return nil
}

// FindStoring retrieves a record from the Storing table by its hash.
func FindStoring(db *sql.DB, hash string) (*Storing, error) {
	var storing Storing
	query := `SELECT hash, name, extension, size, path, date FROM Storing WHERE hash = ?`
	err := db.QueryRow(query, hash).Scan(
		&storing.Hash,
		&storing.Name,
		&storing.Extension,
		&storing.Size,
		&storing.Path,
		&storing.Date,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("error finding record in Storing with hash %s: %v", hash, err)
	}

	return &storing, nil
}

// GetAllStoring retrieves all records from the Storing table.
func GetAllStoring(db *sql.DB) ([]Storing, error) {
	query := `SELECT hash, name, extension, size, path, date FROM Storing`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying Storing table: %v", err)
	}
	defer rows.Close()

	var storingRecords []Storing
	for rows.Next() {
		var record Storing
		err := rows.Scan(&record.Hash, &record.Name, &record.Extension, &record.Size, &record.Path, &record.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning Storing record: %v", err)
		}
		storingRecords = append(storingRecords, record)
	}

	return storingRecords, nil
}
