package p2p

import (
	"bufio"
	"context" // for context usage
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"           // for logging
	"os"            // for file operations
	"path/filepath" // for file path manipulations
	"server/database/models"
	"server/database/operations"
	"strings"
	"sync"
	"time"

	// Add the necessary packages from libp2p, for example:
	"github.com/libp2p/go-libp2p/core/host"    // for host.Host
	"github.com/libp2p/go-libp2p/core/network" // for network.Stream
	"github.com/libp2p/go-libp2p/core/peer"
)

var (
	storing          []models.Storing // Global variable to hold Storing objects
	storingMutex     sync.Mutex       // Mutex to ensure thread-safe access to the global variable
	receivedFileData []byte
	receivedFileExt  string
	dataMutex        sync.Mutex
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

		// Log the header to help track the received type of data
		log.Printf("Received header: %s", header)

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
		} else if header == "requested_file" {
			// Handle requested file transfer

			log.Printf("Handling requested file transfer from peer %s", s.Conn().RemotePeer())

			_, err := receiveRequestedFile(s) // Discard file content since you don't need it
			if err != nil {
				log.Printf("Error receiving requested file from peer: %v", err)
				return
			}

			// Log the successful receipt of the file
			log.Printf("Requested file received successfully from peer: %s", s.Conn().RemotePeer())
		} else if header == "request_all" {
			log.Printf("Received 'send_all' request from peer: %s", s.Conn().RemotePeer())
			handleSendAllRequest(s, db, node, s.Conn().RemotePeer().String())
		} else if header == "requested_storings" {
			// Handle requested_storings
			log.Printf("Handling 'requested_storings' from peer: %s", s.Conn().RemotePeer())

			// Read JSON data
			data, err := io.ReadAll(reader)
			if err != nil {
				log.Printf("Error reading 'requested_storings' data: %v", err)
				return
			}

			// Parse JSON into a slice of Storing objects
			var receivedStorings []models.Storing
			err = json.Unmarshal(data, &receivedStorings)
			if err != nil {
				log.Printf("Error unmarshalling 'requested_storings' data: %v", err)
				return
			}

			// Safely add the received storings to the global storing list
			storingMutex.Lock()
			storing = append(storing, receivedStorings...)
			storingMutex.Unlock()

			log.Printf("Successfully added %d storings to the global list", len(receivedStorings))
		} else if header == "requested_file_ext" {
			// Handle file extension
			log.Printf("Handling file extension transfer from peer: %s", s.Conn().RemotePeer())

			ext, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error reading file extension from stream: %v", err)
				return
			}
			ext = strings.TrimSpace(ext)

			// Safely store the received extension
			dataMutex.Lock()
			receivedFileExt = ext
			dataMutex.Unlock()

			log.Printf("File extension received and stored: %s", receivedFileExt)
		} else {
			log.Printf("Unknown header type received: %s", header)
		}
	})
}
func handleSendAllRequest(s network.Stream, db *sql.DB, node host.Host, targetPeerID string) ([]models.Storing, error) {
	log.Printf("Handling 'send_all' request for peer: %s", targetPeerID)

	// Retrieve all storing records from the database
	storingRecords, err := operations.GetAllStoring(db)
	if err != nil {
		log.Printf("Error retrieving storing records: %v", err)
		return nil, err // Return the error if retrieval fails
	}

	// Decode the target peer ID
	targetPeerIDParsed, err := peer.Decode(targetPeerID)
	if err != nil {
		log.Printf("Failed to decode target peer ID: %v", err)
		return nil, err
	}

	// Open a stream to the target peer
	ctx := context.Background()
	stream, err := node.NewStream(network.WithAllowLimitedConn(ctx, "/senddata/p2p"), targetPeerIDParsed, "/senddata/p2p")
	if err != nil {
		log.Printf("Failed to open stream to peer %s: %v", targetPeerIDParsed, err)
		return nil, err
	}
	defer stream.Close()
	log.Printf("Stream opened successfully to peer %s", targetPeerIDParsed)

	// Send the header to indicate the type of data being sent
	header := "requested_storings\n"
	_, err = stream.Write([]byte(header))
	if err != nil {
		log.Printf("Error sending header to peer %s: %v", targetPeerIDParsed, err)
		return nil, err
	}

	// Serialize the storing records to JSON
	jsonData, err := json.Marshal(storingRecords)
	if err != nil {
		log.Printf("Error serializing storing records to JSON: %v", err)
		return nil, err // Return the error if serialization fails
	}

	// Send the JSON data back to the requesting peer
	_, err = stream.Write(jsonData)
	if err != nil {
		log.Printf("Error sending storing records to peer %s: %v", targetPeerIDParsed, err)
		return nil, err
	}

	log.Printf("All storing records sent successfully to peer: %s", targetPeerIDParsed)

	// Return the list of storing records
	return storingRecords, nil
}

func handleFileRequest(s network.Stream, db *sql.DB, node host.Host, targetPeerID string) {
	log.Printf("Handling file request from peer %s", targetPeerID)

	reader := bufio.NewReader(s)

	// Read the file hash
	fileHash, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading file hash from stream from peer %s: %v", targetPeerID, err)
		sendDataToPeer(node, targetPeerID, "", "Failed to read file hash", "", "", "")
		return
	}
	fileHash = strings.TrimSpace(fileHash)
	log.Printf("Received file hash: %s", fileHash)

	// Read the password
	password, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading password from stream from peer %s: %v", targetPeerID, err)
		sendDataToPeer(node, targetPeerID, "", "Failed to read password", "", "", "")
		return
	}
	password = strings.TrimSpace(password)
	log.Printf("Received password (masked): %s", password)

	// Retrieve file metadata from the database
	log.Printf("Searching for file metadata in the database for hash: %s", fileHash)
	storing, err := operations.FindStoring(db, fileHash)
	if err != nil || storing == nil {
		log.Printf("File not found or error occurred while fetching file metadata for hash %s: %v", fileHash, err)
		sendDataToPeer(node, targetPeerID, "", "File not found", "", "", "")
		return
	}

	log.Printf("Found file metadata for file hash: %s", fileHash)

	log.Printf("Checking password in the Sharing table for file hash: %s", fileHash)
	sharing, err := operations.FindSharing(db, fileHash)
	if err != nil || sharing == nil {
		log.Printf("No password found in the Sharing table for file hash %s: %v", fileHash, err)
		sendDataToPeer(node, targetPeerID, "", "Password not found", "", "", "")
		return
	}
	// Validate the password
	if sharing.Password != password {
		log.Printf("Invalid password provided for file hash: %s", fileHash)
		sendDataToPeer(node, targetPeerID, "", "Invalid password", "", "", "")
		return
	}

	log.Printf("Password validated successfully for file hash: %s", fileHash)

	// Use sendDataToPeer to send the requested file back
	log.Printf("Sending requested file back to peer %s from path: %s", targetPeerID, storing.Path)
	err = sendRequestedFileToPeer(node, targetPeerID, storing.Path)
	if err != nil {
		log.Printf("Error sending requested file to peer %s: %v", targetPeerID, err)
		return
	}

	log.Printf("File sent successfully to peer %s: %s", targetPeerID, storing.Path)

	// Send the file extension
	fileExt := storing.Extension
	if fileExt == "" {
		log.Printf("No extension found for file hash: %s", fileHash)
		fileExt = "unknown"
	}
	err = sendRequestedFileExtToPeer(node, targetPeerID, fileExt)
	if err != nil {
		log.Printf("Error sending file extension to peer %s: %v", targetPeerID, err)
		return
	}
	log.Printf("File extension sent successfully to peer %s: %s", targetPeerID, fileExt)
}

func sendDataToPeer(node host.Host, targetPeerID, filePath, message, dataType string, hash string, password string) error {
	ctx := context.Background()
	targetPeerIDParsed, err := peer.Decode(targetPeerID)
	if err != nil {
		log.Printf("Failed to decode target peer ID: %v", err)
		return err
	}

	// Open a stream to the target peer
	s, err := node.NewStream(network.WithAllowLimitedConn(ctx, "/senddata/p2p"), targetPeerIDParsed, "/senddata/p2p")
	if err != nil {
		log.Printf("Failed to open stream to %s: %v", targetPeerIDParsed, err)
		return err
	}
	defer func() {
		log.Printf("Closing stream to peer %s", targetPeerIDParsed)
		s.Close()
	}()

	// Handle request, message, or file
	if dataType == "request" {
		// Send a file request
		log.Printf("Sending file request to peer %s with hash: %s", targetPeerIDParsed, hash)
		_, err = s.Write([]byte("request\n"))
		if err != nil {
			log.Printf("Failed to send request header to peer %s: %v", targetPeerIDParsed, err)
			return err
		}

		// Write hash and password
		_, err = s.Write([]byte(hash + "\n" + password + "\n"))
		if err != nil {
			log.Printf("Failed to send hash or password to peer %s: %v", targetPeerIDParsed, err)
			return err
		}

		log.Printf("File request sent successfully to peer %s", targetPeerIDParsed)

		// Wait for half a second to let the file be received
		time.Sleep(500 * time.Millisecond)

	} else if dataType == "request_all" {
		log.Printf("Sending 'request_all' signal to peer %s", targetPeerIDParsed)
		_, err = s.Write([]byte("request_all\n"))
		if err != nil {
			log.Printf("Failed to send 'request_all' signal to peer %s: %v", targetPeerIDParsed, err)
			return err
		}
		log.Printf("'Request all files' signal sent successfully to peer %s", targetPeerIDParsed)

		// Wait for half a second to let the file be received
	} else if message != "" {
		// Send a message
		log.Printf("Sending message to peer %s: %s", targetPeerIDParsed, message)
		_, err = s.Write([]byte("message\n" + message + "\n"))
		if err != nil {
			log.Printf("Failed to send message to peer %s: %v", targetPeerIDParsed, err)
			return err
		}
		log.Printf("Message sent successfully to peer %s", targetPeerIDParsed)

	} else if filePath != "" {
		// Send a file
		log.Printf("Sending file to peer %s: %s", targetPeerIDParsed, filePath)
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Failed to open file: %v", err)
			return err
		}
		defer file.Close()

		// Write the "file" header
		_, err = s.Write([]byte("file\n"))
		if err != nil {
			log.Printf("Failed to send file header to peer %s: %v", targetPeerIDParsed, err)
			return err
		}

		// Write the file content
		fileContent, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Error reading file content: %v", err)
			return err
		}

		n, err := s.Write(fileContent)
		if err != nil {
			log.Printf("Failed to send file content to peer %s: %v", targetPeerIDParsed, err)
			return err
		}

		log.Printf("File sent successfully. Total bytes sent: %d to peer %s", n, targetPeerIDParsed)
	} else {
		log.Println("No file, message, or request provided to send.")
		return fmt.Errorf("no data to send")
	}

	return nil
}

func SendRequest(node host.Host, targetPeerID, hash, password string) ([]byte, string, error) {
	// Call sendDataToPeer to send the request
	err := sendDataToPeer(node, targetPeerID, "", "", "request", hash, password)
	if err != nil {
		return nil, "", err
	}

	// Wait and retrieve the global variables
	time.Sleep(500 * time.Millisecond) // Wait for data to be received

	dataMutex.Lock() // Lock the mutex to safely access the global variables
	defer dataMutex.Unlock()

	if receivedFileData == nil || receivedFileExt == "" {
		log.Printf("No data or file extension received after waiting for the requested file")
		return nil, "", nil
	}

	data := receivedFileData // Copy the data
	ext := receivedFileExt   // Copy the file extension
	receivedFileData = nil   // Clear the global variables
	receivedFileExt = ""

	log.Printf("Returning data and extension: %d bytes, ext: %s", len(data), ext)
	return data, ext, nil
}

// Function to receive a requested file and store it in the global variable
func receiveRequestedFile(s network.Stream) ([]byte, error) {
	reader := bufio.NewReader(s)

	// Directly read the file content
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("Error reading requested file data from stream: %v", err)
		return nil, err
	}

	log.Printf("Requested file received successfully with %d bytes", len(data))

	// Store data in the global variable
	dataMutex.Lock()
	receivedFileData = data
	dataMutex.Unlock()

	return data, nil
}

func sendRequestedFileToPeer(node host.Host, targetPeerID, filePath string) error {
	log.Printf("Preparing to send requested file to peer %s, file: %s", targetPeerID, filePath)

	// Decode the target peer ID
	targetPeerIDParsed, err := peer.Decode(targetPeerID)
	if err != nil {
		log.Printf("Failed to decode target peer ID: %v", err)
		return err
	}
	log.Printf("Successfully decoded target peer ID: %s", targetPeerID)

	// Open a stream to the target peer first
	ctx := context.Background()
	s, err := node.NewStream(network.WithAllowLimitedConn(ctx, "/senddata/p2p"), targetPeerIDParsed, "/senddata/p2p")
	if err != nil {
		log.Printf("Failed to open stream to peer %s: %v", targetPeerIDParsed, err)
		return err
	}
	defer s.Close()
	log.Printf("Stream opened successfully to peer %s", targetPeerID)

	// Open the file to send its content
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file %s: %v", filePath, err)
		return err
	}
	defer file.Close()
	log.Printf("File %s opened successfully", filePath)

	// Write the "requested_file" header
	_, err = s.Write([]byte("requested_file\n"))
	if err != nil {
		log.Printf("Failed to send requested_file header to peer %s: %v", targetPeerIDParsed, err)
		return err
	}
	log.Printf("Sent 'requested_file' header to peer %s", targetPeerID)

	// Write the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file content: %v", err)
		return err
	}

	n, err := s.Write(fileContent)
	if err != nil {
		log.Printf("Failed to send file content to peer %s: %v", targetPeerIDParsed, err)
		return err
	}
	log.Printf("Sent %d bytes of requested file content to peer %s", n, targetPeerID)

	log.Printf("Requested file sent successfully to peer %s: %s", targetPeerID, filePath)
	return nil
}

func sendRequestedFileExtToPeer(node host.Host, targetPeerID, fileExt string) error {
	log.Printf("Preparing to send file extension to peer %s, extension: %s", targetPeerID, fileExt)

	// Decode the target peer ID
	targetPeerIDParsed, err := peer.Decode(targetPeerID)
	if err != nil {
		log.Printf("Failed to decode target peer ID: %v", err)
		return err
	}
	log.Printf("Successfully decoded target peer ID: %s", targetPeerID)

	// Open a stream to the target peer
	ctx := context.Background()
	s, err := node.NewStream(network.WithAllowLimitedConn(ctx, "/senddata/p2p"), targetPeerIDParsed, "/senddata/p2p")
	if err != nil {
		log.Printf("Failed to open stream to peer %s: %v", targetPeerIDParsed, err)
		return err
	}
	defer s.Close()
	log.Printf("Stream opened successfully to peer %s", targetPeerID)

	// Write the "requested_file_ext" header
	_, err = s.Write([]byte("requested_file_ext\n"))
	if err != nil {
		log.Printf("Failed to send requested_file_ext header to peer %s: %v", targetPeerIDParsed, err)
		return err
	}
	log.Printf("Sent 'requested_file_ext' header to peer %s", targetPeerID)

	// Write the file extension
	_, err = s.Write([]byte(fileExt + "\n"))
	if err != nil {
		log.Printf("Failed to send file extension to peer %s: %v", targetPeerIDParsed, err)
		return err
	}
	log.Printf("Sent file extension to peer %s: %s", targetPeerID, fileExt)

	return nil
}

func explore(node host.Host) ([]models.Storing, error) {
	listMutex.Lock()                           // Lock the mutex for safe access to peerIDList
	peers := append([]string{}, peerIDList...) // Make a copy to avoid issues with concurrent modifications
	listMutex.Unlock()

	// Request all files from each peer
	for _, peerID := range peers {
		log.Printf("Requesting all files from peer: %s", peerID)

		// Send a generic "request all files" signal to the peer
		err := sendDataToPeer(node, peerID, "", "", "request_all", "", "")
		if err != nil {
			log.Printf("Error requesting all files from peer %s: %v", peerID, err)
			continue
		}

		log.Printf("Request sent to peer %s for all files", peerID)
	}

	// Wait to allow responses to be processed (if needed)
	time.Sleep(1 * time.Second) // Adjust this delay based on your network latency

	// Lock the global storing list to safely access it
	storingMutex.Lock()
	defer storingMutex.Unlock()

	// Make a copy of the storing list
	storedFiles := append([]models.Storing{}, storing...)

	// Print each file's details
	log.Printf("Number of stored files: %d", len(storedFiles))
	for i, file := range storedFiles {
		log.Printf("File %d: Hash=%s, Name=%s, Extension=%s, Size=%d, Path=%s, Date=%s",
			i+1, file.Hash, file.Name, file.Extension, file.Size, file.Path, file.Date)
	}

	// Clear the global storing list
	storing = []models.Storing{}

	// Return the collected storing records
	return storedFiles, nil
}
