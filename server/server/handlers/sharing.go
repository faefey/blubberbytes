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

	"github.com/libp2p/go-libp2p/core/host"
)

func SharingHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	sharingRecords, err := operations.GetAllSharing(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sharingRecords)
}

func AddSharingHandler(w http.ResponseWriter, r *http.Request, node host.Host, db *sql.DB) {
	decoder := json.NewDecoder(r.Body)
	var m models.Sharing
	err := decoder.Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	link, err := p2p.GenerateLink(db, node, m.Hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteSharingHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = operations.DeleteSharing(db, string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
