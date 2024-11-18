package p2p

import (
	"bufio"
	"context" // for context usage
	"database/sql"
	"io"
	"log"           // for logging
	"os"            // for file operations
	"path/filepath" // for file path manipulations
	"server/database"
	"strings"

	// Add the necessary packages from libp2p, for example:
	"github.com/libp2p/go-libp2p/core/host"    // for host.Host
	"github.com/libp2p/go-libp2p/core/network" // for network.Stream
	"github.com/libp2p/go-libp2p/core/peer"
)

func receiveDataFromPeer(node host.Host, db *sql.DB, folderPath string) {
	node.SetStreamHandler("/senddata/p2p", func(s network.Stream) {
		log.Printf("New stream opened from peer: %s", s.Conn().RemotePeer())
		defer func() {
			log.Printf("Stream closed by peer: %s", s.Conn().RemotePeer())
			s.Close()
		}()

		// Read the header to determine the type of data (file or message)
		reader := bufio.NewReader(s)
		header, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading header from peer %s: %v", s.Conn().RemotePeer(), err)
			return
		}
		header = strings.TrimSpace(header)

		if header == "file" {
			// Handle file transfer
			fileName := "node_file.pdf"
			filePath := filepath.Join(folderPath, fileName)
			log.Printf("Receiving file. Saving to path: %s", filePath)

			file, err := os.Create(filePath)
			if err != nil {
				log.Printf("Failed to create file in folder %s: %v", folderPath, err)
				return
			}
			defer file.Close()

			data, err := io.ReadAll(reader)
			if err != nil {
				log.Printf("Error reading file data from stream: %v", err)
				return
			}

			n, err := file.Write(data)
			if err != nil {
				log.Printf("Error writing to file %s: %v", filePath, err)
				return
			}

			log.Printf("File received successfully. Total bytes written: %d to file: %s", n, filePath)

		} else if header == "message" {
			// Handle message transfer
			message, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error reading message from stream: %v", err)
				return
			}
			log.Printf("Message received from peer %s: %s", s.Conn().RemotePeer(), strings.TrimSpace(message))
		} else if header == "request" {
			// Handle file request
			handleFileRequest(s, db, node, s.Conn().RemotePeer().String())
		} else {
			log.Printf("Unknown header type received: %s", header)
		}
	})
}

func handleFileRequest(s network.Stream, db *sql.DB, node host.Host, targetPeerID string) {
	reader := bufio.NewReader(s)

	// Read the file hash
	fileHash, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading file hash from stream: %v", err)
		sendDataToPeer(node, targetPeerID, "", "Failed to read file hash", "", "", "")
		return
	}
	fileHash = strings.TrimSpace(fileHash)

	// Read the password
	password, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading password from stream: %v", err)
		sendDataToPeer(node, targetPeerID, "", "Failed to read password", "", "", "")
		return
	}
	password = strings.TrimSpace(password)

	// Retrieve file metadata from the database
	fileMetadata, err := database.FindFileMetadataByHash(db, fileHash)
	if err != nil || fileMetadata == nil {
		log.Printf("File not found or error occurred: %v", err)
		sendDataToPeer(node, targetPeerID, "", "File not found", "", "", "")
		return
	}

	// Validate the password
	isPasswordValid := false
	for _, p := range fileMetadata.Passwords {
		if p == password {
			isPasswordValid = true
			break
		}
	}
	if !isPasswordValid {
		log.Printf("Invalid password for file hash: %s", fileHash)
		sendDataToPeer(node, targetPeerID, "", "Invalid password", "", "", "")
		return
	}

	// Use sendDataToPeer to send the file back
	sendDataToPeer(node, targetPeerID, fileMetadata.Path, "", "", "", "")

	log.Printf("File sent successfully to peer %s: %s", targetPeerID, fileMetadata.Path)
}

func sendDataToPeer(node host.Host, targetPeerID, filePath, message, request string, hash string, password string) {
	ctx := context.Background()
	targetPeerIDParsed, err := peer.Decode(targetPeerID)
	if err != nil {
		log.Printf("Failed to decode target peer ID: %v", err)
		return
	}

	// Open a stream to the target peer
	s, err := node.NewStream(network.WithAllowLimitedConn(ctx, "/senddata/p2p"), targetPeerIDParsed, "/senddata/p2p")
	if err != nil {
		log.Printf("Failed to open stream to %s: %v", targetPeerIDParsed, err)
		return
	}
	defer func() {
		log.Printf("Closing stream to peer %s", targetPeerIDParsed)
		s.Close()
	}()

	// Handle request, message, or file
	if request == "request" {
		// Send a file request
		log.Printf("Sending file request to peer %s with hash: %s", targetPeerIDParsed, hash)
		_, err = s.Write([]byte("request\n"))
		if err != nil {
			log.Printf("Failed to send request header to peer %s: %v", targetPeerIDParsed, err)
			return
		}

		// Write hash and password
		_, err = s.Write([]byte(hash + "\n" + password + "\n"))
		if err != nil {
			log.Printf("Failed to send hash or password to peer %s: %v", targetPeerIDParsed, err)
			return
		}

		log.Printf("File request sent successfully to peer %s", targetPeerIDParsed)
	} else if message != "" {
		// Send a message
		log.Printf("Sending message to peer %s: %s", targetPeerIDParsed, message)
		_, err = s.Write([]byte("message\n" + message + "\n"))
		if err != nil {
			log.Printf("Failed to send message to peer %s: %v", targetPeerIDParsed, err)
			return
		}
		log.Printf("Message sent successfully to peer %s", targetPeerIDParsed)

	} else if filePath != "" {
		// Send a file
		log.Printf("Sending file to peer %s: %s", targetPeerIDParsed, filePath)
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Failed to open file: %v", err)
			return
		}
		defer file.Close()

		// Write the "file" header
		_, err = s.Write([]byte("file\n"))
		if err != nil {
			log.Printf("Failed to send file header to peer %s: %v", targetPeerIDParsed, err)
			return
		}

		// Write the file content
		fileContent, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Error reading file content: %v", err)
			return
		}

		n, err := s.Write(fileContent)
		if err != nil {
			log.Printf("Failed to send file content to peer %s: %v", targetPeerIDParsed, err)
			return
		}

		log.Printf("File sent successfully. Total bytes sent: %d to peer %s", n, targetPeerIDParsed)
	} else {
		log.Println("No file, message, or request provided to send.")
	}
}
