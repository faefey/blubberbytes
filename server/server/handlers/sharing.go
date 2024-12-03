package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/database/operations"
	"server/p2p"

	"github.com/libp2p/go-libp2p/core/host"
)

// All other sharing handlers are in the gateway folder

func SharingHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	sharingRecords, err := operations.GetAllSharing(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sharingRecords)
}

// /sharefile route:
func AddSharingHandler(w http.ResponseWriter, r *http.Request, node host.Host, db *sql.DB) {
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	hash, err := operations.HashFile(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// the link to return:
	link, err := p2p.GenerateLink(db, node, hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Shareable link: %s\n", link)
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