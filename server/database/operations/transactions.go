package operations

import (
	"database/sql"
	"fmt"
)

func AddTransactions(db *sql.DB, id int64, date, wallet string, amount, balance float64) error {
	query := `INSERT INTO Transactions (id, date, wallet, amount, balance) 
	          VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, id, date, wallet, amount, balance)
	if err != nil {
		return fmt.Errorf("error adding record to Transactions: %v", err)
	}

	fmt.Printf("Record added to Transactions with id: %d\n", id)
	return nil
}
