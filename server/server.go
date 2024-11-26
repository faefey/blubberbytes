// server.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/database/operations"
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

func cors(w http.ResponseWriter, r *http.Request, db *sql.DB, handler func(w http.ResponseWriter, r *http.Request, db *sql.DB)) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
	handler(w, r, db)
}

func server(db *sql.DB) {
	// http.HandleFunc("/downloadFileByHash", downloadFileByHash)
	// http.HandleFunc("/hostFile", hostFile)
	// http.HandleFunc("/setupHTTPProxy", setupHTTPProxy)
	// http.HandleFunc("/viewRandomNeighborFiles", viewRandomNeighborFiles)
	// http.HandleFunc("/generatePublicLink", generatePublicLink)
	// http.HandleFunc("/managePublicLinks", managePublicLinks)

	// fmt.Println("server is running on port 3000...")

	// http.ListenAndServe(":3000", nil)

	http.HandleFunc("/storing", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, storingHandler)
	})

	http.HandleFunc("/hosting", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, hostingHandler)
	})

	http.HandleFunc("/sharing", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, sharingHandler)
	})

	http.HandleFunc("/saved", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, savedHandler)
	})

	http.HandleFunc("/downloads", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, downloadsHandler)
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, transactionsHandler)
	})

	http.HandleFunc("/uploads", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, db, uploadsHandler)
	})

	fmt.Println("Server is running on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func storingHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	storingRecords, err := operations.GetAllStoring(db)
	if err != nil {
		http.Error(w, "Error fetching storing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(storingRecords)
}

func hostingHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	hostingRecords, err := operations.GetAllHosting(db)
	if err != nil {
		http.Error(w, "Error fetching hosting data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hostingRecords)
}

func sharingHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	sharingRecords, err := operations.GetAllSharing(db)
	if err != nil {
		http.Error(w, "Error fetching sharing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sharingRecords)
}

func savedHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	savedRecords, err := operations.GetAllSaved(db)
	if err != nil {
		http.Error(w, "Error fetching saved data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(savedRecords)
}

func downloadsHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	downloadsRecords, err := operations.GetAllDownloads(db)
	if err != nil {
		http.Error(w, "Error fetching download data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(downloadsRecords)
}

func transactionsHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	transactionsRecords, err := operations.GetAllTransactions(db)
	if err != nil {
		http.Error(w, "Error fetching transaction data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactionsRecords)
}

func uploadsHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	uploadsRecords, err := operations.GetAllUploads(db)
	if err != nil {
		http.Error(w, "Error fetching upload data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploadsRecords)
}
