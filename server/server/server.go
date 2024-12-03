package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"server/server/handlers"

	"github.com/libp2p/go-libp2p/core/host"
)

// handler for downloading a file by hash:
func downloadFileByHash(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	fmt.Printf("Received download request for hash: %s\n", hash)
	w.Write([]byte("File download reqeust received")) // to send a response back to the frontend
}

// handler for HTTP proxy setup:
func setupHTTPProxy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Recieved request to setup HTTP proxy")
	w.Write([]byte("HTTP proxy setup request received"))
}

// handler for viewing a random neighbor's files:
func viewRandomNeighborFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request to view random neighbors files")
	w.Write([]byte("Random neighbors files displayed"))
}

func cors(w http.ResponseWriter, r *http.Request, db *sql.DB, handler func(w http.ResponseWriter, r *http.Request, db *sql.DB)) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	} else {
		handler(w, r, db)
	}
}

func corsWithNode(w http.ResponseWriter, r *http.Request, node host.Host, db *sql.DB, handler func(w http.ResponseWriter, r *http.Request, node host.Host, db *sql.DB)) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	} else {
		handler(w, r, node, db)
	}
}

func Server(node host.Host, db *sql.DB) {
	http.HandleFunc("/downloadFileByHash", downloadFileByHash)
	http.HandleFunc("/setupHTTPProxy", setupHTTPProxy)
	http.HandleFunc("/viewRandomNeighborFiles", viewRandomNeighborFiles)

	// GET routes
	http.HandleFunc("/storing", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, handlers.StoringHandler)
	})

	http.HandleFunc("/hosting", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, handlers.HostingHandler)
	})

	http.HandleFunc("/sharing", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, handlers.SharingHandler)
	})

	http.HandleFunc("/saved", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, handlers.SavedHandler)
	})

	http.HandleFunc("/downloads", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, handlers.DownloadsHandler)
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, handlers.TransactionsHandler)
	})

	http.HandleFunc("/uploads", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, handlers.UploadsHandler)
	})

	// POST routes
	http.HandleFunc("/addstoring", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, handlers.AddStoringHandler)
	})

	http.HandleFunc("/deletestoring", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, handlers.DeleteStoringHandler)
	})

	http.HandleFunc("/addhosting", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, handlers.AddHostingHandler)
	})

	http.HandleFunc("/deletehosting", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, handlers.DeleteHostingHandler)
	})

	http.HandleFunc("/addsharing", func(w http.ResponseWriter, r *http.Request) {
		corsWithNode(w, r, node, db, handlers.AddSharingHandler)
	})

	http.HandleFunc("/deletesharing", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, handlers.DeleteSharingHandler)
	})

	http.HandleFunc("/addsaved", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, handlers.AddSavedHandler)
	})

	http.HandleFunc("/deletesaved", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, handlers.DeleteSavedHandler)
	})

	// Run the server
	fmt.Println("Server is running on port 3001...")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
