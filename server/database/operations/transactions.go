package operations

import (
	"database/sql"
	"fmt"
	"server/database/models"
)

func AddTransactions(db *sql.DB, date, wallet string, amount, balance float64) error {
	query := `INSERT INTO Transactions (date, wallet, amount, balance) 
	          VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, date, wallet, amount, balance)
	if err != nil {
		return fmt.Errorf("error adding record to Transactions: %v", err)
	}

	fmt.Printf("Record added to Transactions with wallet: %s\n", wallet)
	return nil
}

func GetAllTransactions(db *sql.DB) ([]models.Transactions, error) {
	query := `SELECT id, date, wallet, amount, balance FROM Transactions`
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
