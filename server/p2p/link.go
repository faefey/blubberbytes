package p2p

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"

	"server/database/operations"

	"github.com/libp2p/go-libp2p/core/host"
)

// GenerateLink generates a link for sharing a file using a hash and stores the details in the database.
func GenerateLink(db *sql.DB, node host.Host, fileHash string) (string, error) {
	// Step 1: Get the node's peer address
	nodeAddress := node.ID()

	// Step 2: Generate a secure random password
	password, err := generateSecurePassword(16) // Length: 16 characters
	if err != nil {
		return "", fmt.Errorf("failed to generate password: %v", err)
	}

	// Step 3: Add the generated password to the file's password list
	err = operations.AddSharing(db, fileHash, password)
	if err != nil {
		return "", fmt.Errorf("failed to add password to file hash: %v", err)
	}

	// Step 4: Generate the shareable link
	link := fmt.Sprintf("http://localhost:3002/viewfile?address=%s&hash=%s&password=%s", nodeAddress, fileHash, password)

	log.Printf("Generated link: %s", link)
	return link, nil
}

// generateSecurePassword generates a secure random password of the specified length.
func generateSecurePassword(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
