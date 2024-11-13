package main

import (
	"fmt"
	"net/http"
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
func downloadFileByHash(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	fmt.Printf("Received download request for hash: %s\n", hash)
	w.Write([]byte("File download reqeust received")) // to send a response back to the frontend
}

// handler for hosting a file:
func hostFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request to host a  file")
	w.Write([]byte("File hosting request received"))
}

// handler for HTTP proxy setup:
func setupHTTPProxy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Recieved request to setup HTTP proxy")
	w.Write([]byte("HTTP proxy setup request received"))
}

//handler for viewing a random neighbor's files:
func viewRandomNeighborFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request to view random neighbors files")
	w.Write([]byte("Random neighbors files displayed"))
}

//  handler for generating a public HTTP gateway link:
func generatePublicLink(w http.ResponseWriter, r *http.Request) {
	fileHash := r.URL.Query().Get("hash")
	fmt.Printf("Generating public link for file hash: %s\n", fileHash)
	w.Write([]byte("Public link generated"))
}

//  handler for managing active public links:
func managePublicLinks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request sent for managing public linsk")
	w.Write([]byte("Managing public links"))
}

// server:
func server() {
	http.HandleFunc("/downloadFileByHash", downloadFileByHash)
	http.HandleFunc("/hostFile", hostFile)
	http.HandleFunc("/setupHTTPProxy", setupHTTPProxy)
	http.HandleFunc("/viewRandomNeighborFiles", viewRandomNeighborFiles)
	http.HandleFunc("/generatePublicLink", generatePublicLink)
	http.HandleFunc("/managePublicLinks", managePublicLinks)

	fmt.Println("server is running on port 3000...")

	http.ListenAndServe(":3000", nil)
}
