package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"crypto/sha256"
	"io"
	"os"
)

// file entry:
type fileEntry struct {
	path        string
	password    string
	expiration  time.Time
	accessLimit int
}

func hashFile(filePath string) (string, error) {
	fileContent, err := os.Open(filePath)
	if err != nil {
		return "Error hashing", err
	}
	defer fileContent.Close()

	h := sha256.New()
	if _, err := io.Copy(h, fileContent); err != nil {
		return "Error hashing", err
	}

	hash := hex.EncodeToString(h.Sum(nil))
	fmt.Println("Hash of file at " + filePath + ": " + hash)
	return hash, nil
}

var fileMapping = struct {
	sync.RWMutex
	m map[string]fileEntry
}{m: make(map[string]fileEntry)}

// shared files:
var sharedFiles = struct {
	sync.RWMutex
	m map[string]string // hash ->password
}{m: make(map[string]string)}

// hashFile from hash.go:
func hashFileWrapper(path string) string {
	hash, err := hashFile(path)
	if err != nil {
		fmt.Printf("Failed to hash file: %v\n", err)
		return "error_hashing_file"
	}
	return hash
}

// serves the files *
func fileHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[len("/hash/"):]
	fileMapping.RLock()
	entry, ok := fileMapping.m[hash]
	fileMapping.RUnlock()

	if !ok || time.Now().After(entry.expiration) || entry.accessLimit <= 0 {
		http.Error(w, "File not found or expired", http.StatusNotFound)
		return
	}

	absPath, err := filepath.Abs(entry.path)
	if err != nil || !isInFilesDir(absPath) {
		http.Error(w, "File not accessible", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, absPath)
	fileMapping.Lock()
	entry.accessLimit--
	fileMapping.m[hash] = entry
	fileMapping.Unlock()
}

func isInFilesDir(filePath string) bool {
	base, err := filepath.Abs("files")
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(base, filePath)
	return err == nil && !filepath.IsAbs(rel) && !filepath.HasPrefix(rel, "..")
}


// random password:
func generatePassword() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "defaultpassword"
	}
	return hex.EncodeToString(bytes)
}

// /sharefile route:
func shareFileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	hash, err := hashFile(filePath) // This should work if hashFile is properly defined
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	password := generatePassword()

	// to store hash and password in sharedFiles:
	sharedFiles.Lock()
	sharedFiles.m[hash] = password
	sharedFiles.Unlock()

	// the link to return:
	link := fmt.Sprintf("http://localhost:3000/viewfile?address=localhost&hash=%s&password=%s", hash, password)

 // to add fule:
	addFile(hash, filePath, time.Hour, 5)

	fmt.Fprintf(w, "Shareable link: %s\n", link)
}

// /viewfile route:
func viewFileHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	hash := r.URL.Query().Get("hash")
	password := r.URL.Query().Get("password")

	if address == "" || hash == "" || password == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// checks if the hahs and password are correct:
	sharedFiles.RLock()
	storedPassword, ok := sharedFiles.m[hash]
	sharedFiles.RUnlock()

	if !ok || storedPassword != password {
		http.Error(w, "Invalid hash or password", http.StatusForbidden)
		return
	}

	// get the file entry:
	fileMapping.RLock()
	entry, exists := fileMapping.m[hash]
	fileMapping.RUnlock()

	if !exists || time.Now().After(entry.expiration) || entry.accessLimit <= 0 {
		http.Error(w, "File not found or expired", http.StatusNotFound)
		return
	}

	// read file bytes:
	fileBytes, err := ioutil.ReadFile(entry.path)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// detect file extension:
	ext := filepath.Ext(entry.path)
	contentType := "application/octet-stream"
	if ext == ".jpg" || ext == ".jpeg" {
		contentType = "image/jpeg"
	} else if ext == ".png" {
		contentType = "image/png"
	} else if ext == ".txt" {
		contentType = "text/plain"
	}

	// serve the file content:
	w.Header().Set("Content-Type", contentType)
	w.Write(fileBytes)

	// decrease access limit
	fileMapping.Lock()
	entry.accessLimit--
	fileMapping.m[hash] = entry
	fileMapping.Unlock()
}

// add file to fileMapping:
func addFile(hash, path string, duration time.Duration, maxAccess int) {
	expiration := time.Now().Add(duration)
	fileMapping.Lock()
	fileMapping.m[hash] = fileEntry{path: path, expiration: expiration, accessLimit: maxAccess}
	fileMapping.Unlock()
	fmt.Printf("Added file: %s %s %v %d\n", hash, path, expiration, maxAccess)
}

// HTTP server
func gateway() {
	http.HandleFunc("/sharefile", shareFileHandler)
	http.HandleFunc("/viewfile", viewFileHandler)
	http.HandleFunc("/hash/", fileHandler)
	fmt.Println("Starting server on http://localhost:3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}

func main() {
	addFile("hash098765", "files/dog.jpg", time.Hour, 2) // Example file entry
	gateway()
}
