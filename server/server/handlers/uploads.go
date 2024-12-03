package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"server/database/operations"
)

func UploadsHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	uploadsRecords, err := operations.GetAllUploads(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploadsRecords)
}
