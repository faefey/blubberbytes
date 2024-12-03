package gateway

import (
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
	fileBytes, ext, err := p2p.SendRequest(node, address, hash, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		// detect file extension:
		contentType := "application/octet-stream"
		if ext == ".jpg" || ext == ".jpeg" {
			contentType = "image/jpeg"
		} else if ext == ".png" {
			contentType = "image/png"
		} else if ext == ".txt" {
			contentType = "text/plain"
		} else if ext == ".pdf" {
			contentType = "application/pdf"
		}

		// serve the file content:
		w.Header().Set("Content-Type", contentType)
		w.Write(fileBytes)
	}
}
