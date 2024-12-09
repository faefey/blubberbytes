package operations

import (
	"database/sql"
	"fmt"
	"server/database/models"
)

// UpdateProxy updates the only record in the Proxy table.
func UpdateProxy(db *sql.DB, ip, port string, rate float64) error {
	query := `UPDATE Proxy SET ip = ?, port = ?, rate = ?`
	_, err := db.Exec(query, ip, port, rate)
	if err != nil {
		return fmt.Errorf("error updating record from Proxy: %v", err)
	}

	fmt.Printf("Record updated successfully in Proxy.\n")
	return nil
}

// GetProxy retrieves the only record from the Proxy table.
func GetProxy(db *sql.DB) (*models.Proxy, error) {
	var proxy models.Proxy
	query := `SELECT ip, port, rate FROM Proxy`
	err := db.QueryRow(query).Scan(&proxy.IP, &proxy.Port, &proxy.Rate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("error finding record in Proxy: %v", err)
	}

	return &proxy, nil
}
