package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"server/database/models"
	"server/database/operations"
)

func GetProvidersHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	providers := []string{"123", "abc", "2024", string(body)}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providers)
}

func RequestMetadataHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metadata := models.JoinedHosting{Name: string(body)}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metadata)
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file := body

	operations.AddDownloads(db, "", "", "", "", 0, 0)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(file)
}
