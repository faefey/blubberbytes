package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"server/server/handlers"

	"github.com/libp2p/go-libp2p/core/host"
)

// Actions:
// -- download a file by hash
// -- host a file
// -- HTTP proxy setup
// -- view a random neighbor's file
// -- generate a public HtTP getaway link
// -- manage active public links
// --

// To test: http://localhost:3000/downloadFileByHash?hash=hash
// run server: go run server.go

// handler for downloading a file by hash:
// func downloadFileByHash(w http.ResponseWriter, r *http.Request) {
// 	hash := r.URL.Query().Get("hash")
// 	fmt.Printf("Received download request for hash: %s\n", hash)
// 	w.Write([]byte("File download reqeust received")) // to send a response back to the frontend
// }

// // handler for hosting a file:
// func hostFile(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Received request to host a  file")
// 	w.Write([]byte("File hosting request received"))
// }

// // handler for HTTP proxy setup:
// func setupHTTPProxy(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Recieved request to setup HTTP proxy")
// 	w.Write([]byte("HTTP proxy setup request received"))
// }

// //handler for viewing a random neighbor's files:
// func viewRandomNeighborFiles(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Request to view random neighbors files")
// 	w.Write([]byte("Random neighbors files displayed"))
// }

// //  handler for generating a public HTTP gateway link:
// func generatePublicLink(w http.ResponseWriter, r *http.Request) {
// 	fileHash := r.URL.Query().Get("hash")
// 	fmt.Printf("Generating public link for file hash: %s\n", fileHash)
// 	w.Write([]byte("Public link generated"))
// }

// //  handler for managing active public links:
// func managePublicLinks(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Request sent for managing public linsk")
// 	w.Write([]byte("Managing public links"))
// }

func Cors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
}

func Server(node host.Host, db *sql.DB) {
	// mux.HandleFunc("/downloadFileByHash", downloadFileByHash)
	// mux.HandleFunc("/hostFile", hostFile)
	// mux.HandleFunc("/setupHTTPProxy", setupHTTPProxy)
	// mux.HandleFunc("/viewRandomNeighborFiles", viewRandomNeighborFiles)
	// mux.HandleFunc("/generatePublicLink", generatePublicLink)
	// mux.HandleFunc("/managePublicLinks", managePublicLinks)

	mux := http.NewServeMux()

	// GET routes
	mux.HandleFunc("GET /storing", func(w http.ResponseWriter, r *http.Request) {
		Cors(w)
		handlers.StoringHandler(w, r, db)
	})

	mux.HandleFunc("GET /hosting", func(w http.ResponseWriter, r *http.Request) {
		Cors(w)
		handlers.HostingHandler(w, r, db)
	})

	mux.HandleFunc("GET /sharing", func(w http.ResponseWriter, r *http.Request) {
		Cors(w)
		handlers.SharingHandler(w, r, db)
	})

	mux.HandleFunc("GET /saved", func(w http.ResponseWriter, r *http.Request) {
		Cors(w)
		handlers.SavedHandler(w, r, db)
	})

	mux.HandleFunc("GET /downloads", func(w http.ResponseWriter, r *http.Request) {
		Cors(w)
		handlers.DownloadsHandler(w, r, db)
	})

	mux.HandleFunc("GET /transactions", func(w http.ResponseWriter, r *http.Request) {
		Cors(w)
		handlers.TransactionsHandler(w, r, db)
	})

	mux.HandleFunc("GET /uploads", func(w http.ResponseWriter, r *http.Request) {
		Cors(w)
		handlers.UploadsHandler(w, r, db)
	})

	// POST routes
	mux.HandleFunc("POST /addstoring", func(w http.ResponseWriter, r *http.Request) {
		Cors(w)
		handlers.AddStoringHandler(w, r, db)
	})

	mux.HandleFunc("POST /deletestoring", func(w http.ResponseWriter, r *http.Request) {
		Cors(w)
		handlers.DeleteStoringHandler(w, r, db)
	})

	mux.HandleFunc("POST /addhosting", func(w http.ResponseWriter, r *http.Request) {
		Cors(w)
		handlers.AddHostingHandler(w, r, db)
	})

	mux.HandleFunc("POST /deletehosting", func(w http.ResponseWriter, r *http.Request) {
		Cors(w)
		handlers.DeleteHostingHandler(w, r, db)
	})

	mux.HandleFunc("POST /addsharing", func(w http.ResponseWriter, r *http.Request) {
		Cors(w)
		handlers.AddSharingHandler(w, r, node, db)
	})

	mux.HandleFunc("POST /deletesharing", func(w http.ResponseWriter, r *http.Request) {
		Cors(w)
		handlers.DeleteSharingHandler(w, r, db)
	})

	mux.HandleFunc("POST /addsaved", func(w http.ResponseWriter, r *http.Request) {
		Cors(w)
		handlers.AddSavedHandler(w, r, db)
	})

	mux.HandleFunc("POST /deletesaved", func(w http.ResponseWriter, r *http.Request) {
		Cors(w)
		handlers.DeleteSavedHandler(w, r, db)
	})

	// Run the server
	fmt.Println("Server is running on port 3001...")
	log.Fatal(http.ListenAndServe(":3001", mux))
}
