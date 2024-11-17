// server.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"server/database"
	"strings"
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

func server(db *sql.DB) {
	// http.HandleFunc("/downloadFileByHash", downloadFileByHash)
	// http.HandleFunc("/hostFile", hostFile)
	// http.HandleFunc("/setupHTTPProxy", setupHTTPProxy)
	// http.HandleFunc("/viewRandomNeighborFiles", viewRandomNeighborFiles)
	// http.HandleFunc("/generatePublicLink", generatePublicLink)
	// http.HandleFunc("/managePublicLinks", managePublicLinks)

	// fmt.Println("server is running on port 3000...")

	// http.ListenAndServe(":3000", nil)

	http.HandleFunc("/hosting", func(w http.ResponseWriter, r *http.Request) {
		serveDataByCategory(w, db, "hosting")
	})

	http.HandleFunc("/sharing", func(w http.ResponseWriter, r *http.Request) {
		serveDataByCategory(w, db, "sharing")
	})

	http.HandleFunc("/purchased", func(w http.ResponseWriter, r *http.Request) {
		serveDataByCategory(w, db, "purchased")
	})

	http.HandleFunc("/explore", func(w http.ResponseWriter, r *http.Request) {
		serveDataByCategory(w, db, "explore")
	})

	http.HandleFunc("/delete/hosting", func(w http.ResponseWriter, r *http.Request) {
		deleteHostingHandler(w, r, db)
	})

	http.HandleFunc("/delete/sharing", func(w http.ResponseWriter, r *http.Request) {
		deleteSharingHandler(w, r, db)
	})

	http.HandleFunc("/downloadFile", func(w http.ResponseWriter, r *http.Request) {
		downloadFileHandler(w, r, db)
	})

	http.HandleFunc("/hostingFile", func(w http.ResponseWriter, r *http.Request) {
		hostingFileHandler(w, r, db)
	})

	fmt.Println("Server is running on port 3005...")
	http.ListenAndServe(":3005", nil)
}

func serveDataByCategory(w http.ResponseWriter, db *sql.DB, category string) {
	files, err := database.GetFileDataFromTable(db, category)
	if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func deleteHostingHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		FileHash string `json:"file_hash"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data !!!!!", http.StatusBadRequest)
		return
	}

	if requestData.FileHash == "" {
		http.Error(w, "File hash is required !!!!!!!", http.StatusBadRequest)
		return
	}

	// Delete from the hosting table
	err = database.DeleteHosting(db, requestData.FileHash)
	if err != nil {
		if strings.Contains(err.Error(), "no file found :( ") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "File deleted from hosting successfully :) ",
	})
}

func deleteSharingHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		FileHash string `json:"file_hash"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data !!!", http.StatusBadRequest)
		return
	}

	if requestData.FileHash == "" {
		http.Error(w, "File hash is required!!!!", http.StatusBadRequest)
		return
	}

	// Delete from sharing table
	err = database.DeleteSharing(db, requestData.FileHash)
	if err != nil {
		if strings.Contains(err.Error(), "no file found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "File deleted from sharing successfully :)",
	})
}

func downloadFileHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body to get the file information
	var fileData database.FileData
	err := json.NewDecoder(r.Body).Decode(&fileData)
	if err != nil {
		http.Error(w, "Invalid JSON data !!!!", http.StatusBadRequest)
		return
	}

	if fileData.Hash == "" {
		http.Error(w, "File hash is required", http.StatusBadRequest)
		return
	}

	// Check if the file already exists in the 'purchased' table
	exists, err := database.FileExistsInTable(db, "purchased", fileData.Hash)
	if err != nil {
		http.Error(w, "Error checking file existence", http.StatusInternalServerError)
		return
	}

	if exists {
		fmt.Println("File already exists in 'purchased' table")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "File already exists in 'purchased' table",
		})
		return
	}

	// File does not exist; insert it into the 'purchased' table
	err = database.AddFileDataToTable(db, "purchased", fileData)
	if err != nil {
		http.Error(w, "Error inserting file into 'purchased' table", http.StatusInternalServerError)
		return
	}

	fmt.Println("File added to 'purchased' table successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "File added to 'purchased' table successfully",
	})
}

func hostingFileHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body to get the file information
	var fileData database.FileData
	err := json.NewDecoder(r.Body).Decode(&fileData)
	if err != nil {
		http.Error(w, "Invalid JSON data !!!!", http.StatusBadRequest)
		return
	}

	if fileData.Hash == "" {
		http.Error(w, "File hash is required", http.StatusBadRequest)
		return
	}

	// Check if the file already exists in the 'hosting' table
	exists, err := database.FileExistsInTable(db, "hosting", fileData.Hash)
	if err != nil {
		http.Error(w, "Error checking file existence", http.StatusInternalServerError)
		return
	}

	if exists {
		fmt.Println("File already exists in 'hosting' table")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "File already exists in 'hosting' table",
		})
		return
	}

	// File does not exist; insert it into the 'hosting' table
	err = database.AddFileDataToTable(db, "hosting", fileData)
	if err != nil {
		http.Error(w, "Error inserting file into 'hosting' table", http.StatusInternalServerError)
		return
	}

	fmt.Println("File added to 'hosting' table successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "File added to 'hosting' table successfully",
	})
}
