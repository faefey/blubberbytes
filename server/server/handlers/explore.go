package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"server/p2p"

	"github.com/libp2p/go-libp2p/core/host"
)

func ExploreHandler(w http.ResponseWriter, r *http.Request, node host.Host, db *sql.DB) {
	decoder := json.NewDecoder(r.Body)
	var request []string
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	explore, err := p2p.Explore(node, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(explore)
}
