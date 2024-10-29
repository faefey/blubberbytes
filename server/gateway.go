package main

import (
	"fmt"
	"log"
	"net/http"

	//"os"
	"path/filepath"
)

// to test ---> go run gateway.go
// go to http://localhost:3000/hash098765

// mapping file hashes to file paths
var fileMapping = map[string]string{
	"hash098765": "files/dog.jpg",
}

// fileHandler handles requests to /hash and serves the file
func fileHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[len("/"):] // to get hash from URL path

	//  get the file path from the hash
	filePath, ok := fileMapping[hash]
	if !ok {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// get absolute path to prevent directory traversal attacks
	absPath, err := filepath.Abs(filePath)
	if err != nil || !isInFilesDir(absPath) {
		http.Error(w, "File not accessible", http.StatusForbidden)
		return
	}

	// serve the file
	http.ServeFile(w, r, absPath)
}

// isInFilesDir makes sure that the requested file is ni the files directory
func isInFilesDir(filePath string) bool {
	base, err := filepath.Abs("files")
	if err != nil {
		return false
	}
	return filepath.HasPrefix(filePath, base)
}

func gateway() {
	http.HandleFunc("/", fileHandler) // route all requests to fileHandler

	fmt.Println("Starting server on http://localhost:3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
