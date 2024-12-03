package operations

import (
	"database/sql"
	"fmt"
	"server/database/models"
)

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
func FindHosting(db *sql.DB, hash string) (*models.JoinedHosting, error) {
	var hosting models.JoinedHosting
	query := `SELECT Storing.*, price FROM Hosting JOIN Storing ON Hosting.hash == Storing.hash WHERE Hosting.hash = ?`
	err := db.QueryRow(query, hash).Scan(&hosting.Hash, &hosting.Name, &hosting.Extension, &hosting.Size, &hosting.Path, &hosting.Date, &hosting.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("error finding record in Hosting with hash %s: %v", hash, err)
	}

	return &hosting, nil
}

// GetAllHosting retrieves all records from the Hosting table.
func GetAllHosting(db *sql.DB) ([]models.JoinedHosting, error) {
	query := `SELECT Storing.*, price FROM Hosting JOIN Storing ON Hosting.hash == Storing.hash`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying Hosting table: %v", err)
	}
	defer rows.Close()

	hostingRecords := []models.JoinedHosting{}
	for rows.Next() {
		var record models.JoinedHosting
		err := rows.Scan(&record.Hash, &record.Name, &record.Extension, &record.Size, &record.Path, &record.Date, &record.Price)
		if err != nil {
			return nil, fmt.Errorf("error scanning Hosting record: %v", err)
		}
		hostingRecords = append(hostingRecords, record)
	}

	return hostingRecords, nil
}
