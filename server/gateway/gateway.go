package gateway

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/libp2p/go-libp2p/core/host"
)

// HTTP server
func Gateway(node host.Host, db *sql.DB) {
	http.HandleFunc("/viewfile", func(w http.ResponseWriter, r *http.Request) {
		viewFileHandler(w, r, node)
	})

	fmt.Println("Starting server on http://localhost:3002")
	if err := http.ListenAndServe(":3002", nil); err != nil {
		panic(fmt.Sprintf("Server failed: %s", err))
	}
}
