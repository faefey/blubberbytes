package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"server/database/operations"
)

func StatisticsHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	statistics, err := operations.CalcStatistics(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statistics)
}
