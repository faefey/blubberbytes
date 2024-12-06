package gateway

import (
	"fmt"
	"net/http"
	"server/p2p"

	"github.com/libp2p/go-libp2p/core/host"
)

// /viewfile route:
func viewFileHandler(w http.ResponseWriter, r *http.Request, node host.Host) {
	address := r.URL.Query().Get("address")
	hash := r.URL.Query().Get("hash")
	password := r.URL.Query().Get("password")

	if address == "" || hash == "" || password == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// read file bytes:
	name, data, ext, err := p2p.SendRequest(node, address, hash, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		// detect file type:
		var contentType string
		if ext == "" {
			contentType = "application/octet-stream"
		} else {
			contentType = ext
		}

		// serve the file content:
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", name))
		w.Write(data)
	}
}
