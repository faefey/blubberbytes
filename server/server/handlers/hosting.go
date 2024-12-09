package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/database/models"
	"server/database/operations"
	"server/p2p"
)

func HostingHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	hostingRecords, err := operations.GetAllHosting(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hostingRecords)
}

func AddHostingHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	decoder := json.NewDecoder(r.Body)
	var m models.Hosting
	err := decoder.Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	record, err := operations.FindHosting(db, m.Hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if record != nil {
		fmt.Fprint(w, "The file is already being hosted.")
		return
	}

	err = p2p.ProvideKey(m.Hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = operations.AddHosting(db, m.Hash, m.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteHostingHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = operations.DeleteHosting(db, string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
