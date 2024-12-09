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

func GetProvidersHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	providers, err := p2p.GetProviderIDs(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providers)
}

func RequestMetadataHandler(w http.ResponseWriter, r *http.Request, node host.Host, db *sql.DB) {
	decoder := json.NewDecoder(r.Body)
	var request struct {
		Peer string `json:"peer"`
		Hash string `json:"hash"`
	}
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metadata, err := p2p.RequestFileInfo(node, request.Peer, request.Hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metadata)
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	decoder := json.NewDecoder(r.Body)
	var request struct {
		Peer string `json:"peer"`
		Hash string `json:"hash"`
	}
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = operations.AddDownloads(db, "", "", "", "", 0, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "walletaddress")
}
