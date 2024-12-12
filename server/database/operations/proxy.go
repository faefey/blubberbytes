package operations

import (
	"database/sql"
	"fmt"
	"server/database/models"
)

// UpdateProxy updates the only record in the Proxy table.
func UpdateProxy(db *sql.DB, ip string, rate float64, node string, wallet string) error {
	query := `UPDATE Proxy SET ip = ?, rate = ?, node = ?, wallet = ?`
	_, err := db.Exec(query, ip, rate, node, wallet)
	if err != nil {
		return fmt.Errorf("error updating record from Proxy: %v", err)
	}

	fmt.Printf("Record updated successfully in Proxy.\n")
	return nil
}

// GetProxy retrieves the only record from the Proxy table.
func GetProxy(db *sql.DB) (*models.Proxy, error) {
	var proxy models.Proxy
	query := `SELECT ip, rate, node, wallet FROM Proxy`
	err := db.QueryRow(query).Scan(&proxy.IP, &proxy.Rate, &proxy.Node, &proxy.Wallet)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("error finding record in Proxy: %v", err)
	}

	return &proxy, nil
}

// AddSaved inserts a new record into the Saved table.
func AddProxyLogs(db *sql.DB, ip string, bytes, time int64) error {
	query := `INSERT INTO ProxyLogs (ip, bytes, time) VALUES (?, ?, ?)`
	_, err := db.Exec(query, ip, bytes, time)
	if err != nil {
		return fmt.Errorf("error adding record to ProxyLogs: %v", err)
	}

	fmt.Printf("Record added to ProxyLogs\n")
	return nil
}

func AddIPtoNode(db *sql.DB, ip, node string) error {
	query := `INSERT INTO IPtoNode (ip, node) VALUES (?, ?)`
	_, err := db.Exec(query, ip, node)
	if err != nil {
		return fmt.Errorf("error adding record to IPtoNode: %v", err)
	}

	fmt.Printf("Record added to IPtoNode\n")
	return nil
}

// func CalcProxyBill(db *sql.DB) error {
// 	query := `SELECT SUM(bytes) FROM `
// }
