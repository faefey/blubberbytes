package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"server/database/operations"
)

func TransactionsHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	transactionsRecords, err := operations.GetAllTransactions(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactionsRecords)
}
