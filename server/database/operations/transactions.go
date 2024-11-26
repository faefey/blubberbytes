package operations

import (
	"database/sql"
	"fmt"
	"server/database/models"
)

func AddTransactions(db *sql.DB, id int64, date, wallet string, amount, balance float64) error {
	query := `INSERT INTO Transactions (id, date, wallet, amount, balance) 
	          VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, id, date, wallet, amount, balance)
	if err != nil {
		return fmt.Errorf("error adding record to Transaction History: %v", err)
	}

	fmt.Printf("Record added to Transaction History with id: %d\n", id)
	return nil
}

func GetAllTransactions(db *sql.DB) ([]models.Transactions, error) {
	query := `SELECT id, date, wallet, amount, balance FROM Transactions`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying Transaction History table: %v", err)
	}
	defer rows.Close()

	var transactionsRecords []models.Transactions
	for rows.Next() {
		var record models.Transactions
		err := rows.Scan(&record.Id, &record.Date, &record.Wallet, &record.Amount, &record.Balance)
		if err != nil {
			return nil, fmt.Errorf("error scanning Transaction History record: %v", err)
		}
		transactionsRecords = append(transactionsRecords, record)
	}

	return transactionsRecords, nil
}
