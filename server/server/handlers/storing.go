package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"server/database/models"
	"server/database/operations"
)

func StoringHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	storingRecords, err := operations.GetAllStoring(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(storingRecords)
}

func AddStoringHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	decoder := json.NewDecoder(r.Body)
	var m models.Storing
	err := decoder.Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hash, err := operations.HashFile(m.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = operations.AddStoring(db, hash, m.Name, m.Extension, m.Path, m.Date, m.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteStoringHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = operations.DeleteStoring(db, string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
