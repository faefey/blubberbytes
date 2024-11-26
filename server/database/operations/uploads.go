package operations

import (
	"database/sql"
	"fmt"
	"server/database/models"
)

func AddUploads(db *sql.DB, id int64, date, hash, name, extension string, size int64) error {
	query := `INSERT INTO Uploads (id, date, hash, name, extension, size) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, id, date, hash, name, extension, size)
	if err != nil {
		return fmt.Errorf("error adding record to Uploads: %v", err)
	}

	fmt.Printf("Record added to Uploads with id: %d\n", id)
	return nil
}

func GetAllUploads(db *sql.DB) ([]models.Uploads, error) {
	query := `SELECT id, date, hash, name, extension, size FROM Uploads`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying Uploads table: %v", err)
	}
	defer rows.Close()

	var uploadsRecords []models.Uploads
	for rows.Next() {
		var record models.Uploads
		err := rows.Scan(&record.Id, &record.Date, &record.Hash, &record.Name, &record.Extension, &record.Size)
		if err != nil {
			return nil, fmt.Errorf("error scanning Uploads record: %v", err)
		}
		uploadsRecords = append(uploadsRecords, record)
	}

	return uploadsRecords, nil
}

func GetAllTransactions(db *sql.DB) ([]models.Transactions, error) {
	query := `SELECT id, date, hash, wallet, amount, balance FROM Transactions`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying Transactions table: %v", err)
	}
	defer rows.Close()

	var transactionsRecords []models.Transactions
	for rows.Next() {
		var record models.Transactions
		err := rows.Scan(&record.Id, &record.Date, &record.Wallet, &record.Amount, &record.Balance)
		if err != nil {
			return nil, fmt.Errorf("error scanning Transactions record: %v", err)
		}
		transactionsRecords = append(transactionsRecords, record)
	}

	return transactionsRecords, nil
}
