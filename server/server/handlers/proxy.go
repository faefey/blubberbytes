package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"server/database/models"
	"server/database/operations"
)

func UpdateProxyHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	decoder := json.NewDecoder(r.Body)
	var m models.Proxy
	err := decoder.Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = operations.UpdateProxy(db, m.IP, m.Port, m.Rate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RefreshProxiesHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	err := fmt.Errorf("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]models.Proxy{})
}
