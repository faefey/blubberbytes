package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

var fileMapping = struct {
	sync.RWMutex
	m map[string]fileEntry
}{m: make(map[string]fileEntry)}

type fileEntry struct {
	path        string
	expiration  time.Time
	accessLimit int
}

// fileHandler serves files based on hash
func fileHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[len("/hash/"):]
	//fmt.Println("received hash:", hash) // Debugging line

	fileMapping.RLock()
	entry, ok := fileMapping.m[hash]
	fileMapping.RUnlock()

	//fmt.Println("file exists in map:", ok)
	//fmt.Println("file entry:", entry)

	if !ok || time.Now().After(entry.expiration) || entry.accessLimit <= 0 {
		http.Error(w, "File not found or expired", http.StatusNotFound)
		//fmt.Println("expiration:", entry.expiration)
		//fmt.Println("acccess Limit:", entry.accessLimit)
		return
	}

	absPath, err := filepath.Abs(entry.path)
	//fmt.Println("Absolute Path:", absPath)

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

func addFile(hash, path string, duration time.Duration, maxAccess int) {
	expiration := time.Now().Add(duration)
	fileMapping.Lock()
	fileMapping.m[hash] = fileEntry{path: path, expiration: expiration, accessLimit: maxAccess}
	fileMapping.Unlock()
	fmt.Printf("Added file: %s %s %v %d\n", hash, path, expiration, maxAccess)
}

//  HTTP server
func gateway() {
	http.HandleFunc("/hash/", fileHandler)
	fmt.Println("Starting server on http://localhost:3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}

func main() {
	addFile("hash098765", "files/dog.jpg", time.Hour, 2) // 2 --> access limit 
	fileMapping.RLock()
	fmt.Println("File mapping after adding file:", fileMapping.m)
	fileMapping.RUnlock()

	gateway()
}
