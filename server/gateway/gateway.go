package gateway

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"server/server"

	"github.com/libp2p/go-libp2p/core/host"
)

// HTTP server
func Gateway(node host.Host, db *sql.DB) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /viewfile", func(w http.ResponseWriter, r *http.Request) {
		server.Cors(w)
		viewFileHandler(w, r, node)
	})

	fmt.Println("Starting server on http://localhost:3002")
	if err := http.ListenAndServe(":3002", nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
