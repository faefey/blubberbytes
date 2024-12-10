package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"server/database/models"
)

func ExploreHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	decoder := json.NewDecoder(r.Body)
	var request []string
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]models.JoinedHosting{})
}
