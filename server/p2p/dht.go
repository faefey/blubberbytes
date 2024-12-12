package p2p

import (
	"bufio"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"server/database/models"
	"server/database/operations"
	"strconv"
	"strings"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multihash"
)

// FileMetadata stores metadata about a file
type FileMetadata struct {
	FileSize      int64  `json:"file_size"`      // Size of the file
	Extension     string `json:"extension"`      // File extension
	DownloadTimes int    `json:"download_times"` // Number of times the file has been downloaded
}

// ProviderFileMetadata stores information specific to each provider of the file
type ProviderFileMetadata struct {
	PeerID    string  `json:"peer_id"`    // Peer ID of the provider
	FileName  string  `json:"file_name"`  // Name of the file provided by this peer
	FilePrice float64 `json:"file_price"` // Price of the file provided by this peer
}

// FileRecord stores metadata and a list of providers for a file in the DHT
type FileRecord struct {
	Metadata  FileMetadata           `json:"metadata"`
	Providers []ProviderFileMetadata `json:"providers"`
}

// Function to hash file content and return a base58-encoded multihash string
func hashFileContent(filePath string) (string, error) {
	log.Printf("Opening file: %s\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	log.Printf("Hashing file content...\n")
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file content: %w", err)
	}

	log.Printf("Encoding hash as multihash...\n")
	mh, err := multihash.EncodeName(hash.Sum(nil), "sha2-256")
	if err != nil {
		return "", fmt.Errorf("failed to encode multihash: %w", err)
	}

	log.Printf("Encoding multihash to base58 string...\n")
	encoded, err := multibase.Encode(multibase.Base58BTC, mh)
	if err != nil {
		return "", fmt.Errorf("failed to encode multihash to base58: %w", err)
	}

	log.Printf("File hash (multihash): %s\n", encoded)
	return encoded, nil
}

// Function to check if the file hash exists in the DHT
func fileHashExists(ctx context.Context, dht *dht.IpfsDHT, dhtKey string) (bool, error) {
	log.Printf("Checking if file hash exists under DHT key: %s\n", dhtKey)
	_, err := dht.GetValue(ctx, dhtKey)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("File hash does not exist in DHT.\n")
			return false, nil
		}
		return false, fmt.Errorf("failed to check file hash in DHT: %w", err)
	}
	log.Printf("File hash already exists in DHT.\n")
	return true, nil
}

// Function to get file metadata
func getFileMetadata(filePath string) (FileMetadata, error) {
	log.Printf("Getting metadata for file: %s\n", filePath)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return FileMetadata{}, fmt.Errorf("failed to get file info: %w", err)
	}

	fileSize := fileInfo.Size()
	fileExt := strings.ToLower(filepath.Ext(filePath))

	metadata := FileMetadata{
		FileSize:      fileSize,
		Extension:     fileExt,
		DownloadTimes: 0, // Initially, the download count is 0
	}

	log.Printf("File metadata: Size = %d, Extension = %s, DownloadTimes = %d\n", metadata.FileSize, metadata.Extension, metadata.DownloadTimes)
	return metadata, nil
}

// Function to store file metadata in the DHT with the new structure
func storeFileInDHT(ctx context.Context, dht *dht.IpfsDHT, filePath string, filePrice float64) error {
	// Step 1: Hash the file content
	log.Printf("Hashing file content for: %s\n", filePath)
	fileHash, err := hashFileContent(filePath)
	if err != nil {
		return fmt.Errorf("failed to hash file: %w", err)
	}
	fmt.Printf("File hash (key): %s\n", fileHash)

	// Step 2: Retrieve or create the DHT key for the file hash
	dhtKey := "/orcanet/" + fileHash
	var fileRecord FileRecord
	exists, err := fileHashExists(ctx, dht, dhtKey)
	if err != nil {
		return err
	}

	if exists {
		// Fetch the existing record and decode it
		log.Printf("Retrieving existing record from DHT for key: %s\n", dhtKey)
		value, err := dht.GetValue(ctx, dhtKey)
		if err != nil {
			return fmt.Errorf("failed to get existing file record from DHT: %w", err)
		}
		if err := json.Unmarshal(value, &fileRecord); err != nil {
			return fmt.Errorf("failed to unmarshal existing file record: %w", err)
		}
	} else {
		// If not existing, initialize a new file record with metadata
		fileMetadata, err := getFileMetadata(filePath)
		if err != nil {
			return fmt.Errorf("failed to get file metadata: %w", err)
		}
		fileRecord = FileRecord{
			Metadata:  fileMetadata,
			Providers: []ProviderFileMetadata{},
		}
	}

	// Step 3: Add provider information
	peerID := dht.Host().ID().String()
	fileName := filepath.Base(filePath)
	provider := ProviderFileMetadata{
		PeerID:    peerID,
		FileName:  fileName,
		FilePrice: filePrice,
	}
	fileRecord.Providers = append(fileRecord.Providers, provider)

	// Step 4: Serialize and store the updated file record in the DHT
	log.Printf("Serializing and storing updated file record under DHT key: %s\n", dhtKey)
	fileRecordJSON, err := json.Marshal(fileRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal file record: %w", err)
	}

	err = dht.PutValue(ctx, dhtKey, fileRecordJSON)
	if err != nil {
		return fmt.Errorf("failed to store file record in DHT: %w", err)
	}
	log.Println("File record with metadata and providers successfully stored in DHT.")

	return nil
}

// Helper function to perform periodic tasks
func periodicTaskHelper(interval time.Duration, db *sql.DB) {
	// Create a ticker
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := provideAllKeys(db)
			if err != nil {
				log.Printf("Error in periodic task: %v\n", err)
			}
		}
	}
}

// Function to provide all keys from the Hosting table
func provideAllKeys(db *sql.DB) error {
	// Retrieve all hosting records
	hostingRecords, err := operations.GetAllHosting(db)
	if err != nil {
		return fmt.Errorf("error retrieving hosting records: %v", err)
	}

	// Provide each hosting record's key to the DHT
	for _, record := range hostingRecords {
		err := ProvideKey(record.Hash)
		if err != nil {
			log.Printf("Error providing key for hash %s: %v\n", record.Hash, err)
		}
	}
	return nil
}

func handleInput(ctx context.Context, dht *dht.IpfsDHT, node host.Host, db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("User Input \n ")
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n') // Read input from keyboard
		input = strings.TrimSpace(input)    // Trim any trailing newline or spaces
		args := strings.Split(input, " ")
		if len(args) < 1 {
			fmt.Println("No command provided")
			continue
		}
		command := args[0]
		command = strings.ToUpper(command)
		switch command {
		case "SEND_BILL":
			if len(args) < 2 {
				fmt.Println("PeerID required for SEND_BILL command")
				continue
			}
			peerID := args[1]

			// Create a ProxyBill instance with sample data
			proxyBill := models.ProxyBill{
				IP:     "192.168.1.100",
				Rate:   0.02,
				Bytes:  2048,
				Amount: 0.04,
				Wallet: "0xABCDEF1234567890",
			}

			// Call the new function to send the ProxyBill and wait for confirmation
			err := SendProxyBillWithConfirmation(node, peerID, proxyBill)
			if err != nil {
				fmt.Printf("Error during ProxyBill transaction: %v\n", err)
			} else {
				fmt.Println("ProxyBill transaction completed successfully")
			}

		case "ACA":
			if len(args) < 2 {
				fmt.Println("PeerID required for ACA command")
				continue
			}
			peerID := args[1]
			err := sendDataToPeer(node, peerID, "", "", "request_all", "", "")
			if err != nil {
				fmt.Printf("Failed to send data to peer: %v\n", err)
			} else {
				fmt.Println("Data sent to peer successfully")
			}

		case "UPDATE_WALLET_INFO":
			// Random test data for the WalletInfo table
			address := "1BitcoinWalletAddress"
			pubPassphrase := "publicPass123"
			privPassphrase := "privatePass456"

			// Update the wallet address
			err := operations.UpdateWalletAddress(db, address)
			if err != nil {
				fmt.Printf("Error updating wallet address: %v\n", err)
			} else {
				fmt.Println("Wallet address updated successfully with test data!")
			}

			// Update the wallet passphrases
			err = operations.UpdateWalletPassphrases(db, pubPassphrase, privPassphrase)
			if err != nil {
				fmt.Printf("Error updating wallet passphrases: %v\n", err)
			} else {
				fmt.Println("Wallet passphrases updated successfully with test data!")
			}

		case "UPDATE_PROXY":
			// Random test data for the Proxy table
			ip := "192.168.0.100"
			rate := 50.0
			address := "123"

			// Call the UpdateProxy function with the random test data
			err := operations.UpdateProxy(db, ip, rate, node_id, address)
			if err != nil {
				fmt.Printf("Error updating proxy: %v\n", err)
			} else {
				fmt.Println("Proxy updated successfully with test data!")
			}

		case "PROXY":
			// Call the handleProxyRequest function
			proxies, err := RandomProxiesInfo(node)
			if err != nil {
				log.Fatalf("Error handling proxy request: %v", err)
			}

			log.Printf("Received proxies: %+v", proxies)

		case "REQUEST_FILE_INFO":
			if len(args) < 3 {
				fmt.Println("Usage: REQUEST_FILE_INFO <peerID> <file_hash>")
				continue
			}
			targetPeerID := args[1]
			hash := args[2]

			// Call requestFileInfo function
			RequestFileInfo(node, targetPeerID, hash)

		case "FIND_SHARING":
			if len(args) < 2 {
				fmt.Println("Usage: FIND <hash>")
				continue
			}
			hash := args[1] // Extract the hash from the user input
			sharing, err := operations.FindSharing(db, hash)
			if err != nil {
				fmt.Printf("Error finding sharing record for hash %s: %v\n", hash, err)
				continue
			}

			if sharing == nil {
				fmt.Printf("No record found for hash %s\n", hash)
			} else {
				fmt.Printf("Record found:\nHash: %s\nName: %s\nExtension: %s\nSize: %d bytes\nPath: %s\nDate: %s\nPassword: %s\n",
					sharing.Hash, sharing.Name, sharing.Extension, sharing.Size, sharing.Path, sharing.Date, sharing.Password)
			}

		case "CONNECT":
			if len(args) < 2 {
				fmt.Println("Usage: CONNECT <peerID>")
				continue
			}
			peerID := args[1] // Extract the peerID from user input
			connectToPeerUsingRelay(node, peerID)
		case "EXPLORE":
			// Define a test list of peer IDs (placeholder for your two nodes)
			peerIDs := []string{
				"12D3KooWBe5oQUZbpupTUeEaw4PgZ9DMtnPQ3ebnuU1TXqxN3poP", // Replace with the actual first peer ID
				"12D3KooWHbbocnbPcrckuofGodoSBeyDE67ixLm9g45MUF8LApFv", // Replace with the actual second peer ID
			}

			// Call the explore function with the test list of peer IDs
			collectedHostings, err := Explore(node, peerIDs)
			if err != nil {
				log.Printf("Error during explore: %v", err)
			} else {
				log.Printf("Explore completed. Total collected hostings: %d", len(collectedHostings))
				for i, hosting := range collectedHostings {
					log.Printf("Hosting %d: %+v", i+1, hosting)
				}
			}

		case "DELETE_SHARING":
			if len(args) < 2 {
				fmt.Println("Expected file hash")
				continue
			}
			fileHash := args[1]

			// Delete file metadata from the Storing table
			err := operations.DeleteSharing(db, fileHash)
			if err != nil {
				fmt.Printf("Error deleting file with hash %s: %v\n", fileHash, err)
				continue
			}
			fmt.Printf("File with hash %s deleted successfully.\n", fileHash)

		case "PRINT":
			printPeerList()
		case "GENERATE_LINK":
			if len(args) < 2 {
				fmt.Println("Usage: GENERATE_LINK <file_hash>")
				continue
			}

			fileHash := args[1]

			// Generate a shareable link for the file using the hash
			link, err := GenerateLink(db, node, fileHash)
			if err != nil {
				fmt.Printf("Error generating link for file hash %s: %v\n", fileHash, err)
				continue
			}

			fmt.Printf("Generated Link: %s\n", link)
		case "FIND_FILE_BY_HASH":
			if len(args) < 2 {
				fmt.Println("Usage: FIND_FILE_BY_HASH <file_hash>")
				continue
			}

			fileHash := args[1]

			// Retrieve file metadata from the database by file hash
			file, err := operations.FindStoring(db, fileHash)
			if err != nil {
				fmt.Printf("Error finding file by hash: %v\n", err)
				continue
			}

			if file == nil {
				fmt.Println("No file found with the given hash.")
			} else {
				fmt.Printf("File found:\n")
				fmt.Printf("Hash: %s\n", file.Hash)
				fmt.Printf("Name: %s\n", file.Name)
				fmt.Printf("Extension: %s\n", file.Extension)
				fmt.Printf("Size: %d bytes\n", file.Size)
				fmt.Printf("Path: %s\n", file.Path)
				fmt.Printf("Date: %s\n", file.Date)
			}

		case "ADD_FILE":
			if len(args) < 2 {
				fmt.Println("Expected file path")
				continue
			}
			filePath := args[1]

			// Get file information
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				fmt.Printf("Error accessing file: %v\n", err)
				continue
			}

			// Generate a unique hash for the file content
			fileHash, err := hashFileContent(filePath)
			if err != nil {
				fmt.Printf("Error generating file hash: %v\n", err)
				continue
			}

			fileSize := fileInfo.Size()
			fileName := fileInfo.Name()
			fileDate := time.Now().Format("2006-01-02 15:04:05") // Example: Current timestamp

			// Store file metadata in the Storing table
			err = operations.AddStoring(db, fileHash, fileName, filepath.Ext(filePath), filePath, fileDate, fileSize)
			if err != nil {
				fmt.Printf("Error adding file metadata: %v\n", err)
				continue
			}
			fmt.Printf("File metadata added successfully with hash: %s\n", fileHash)

		case "DELETE_FILE":
			if len(args) < 2 {
				fmt.Println("Expected file hash")
				continue
			}
			fileHash := args[1]

			// Delete file metadata from the Storing table
			err := operations.DeleteStoring(db, fileHash)
			if err != nil {
				fmt.Printf("Error deleting file with hash %s: %v\n", fileHash, err)
				continue
			}
			fmt.Printf("File with hash %s deleted successfully.\n", fileHash)

		case "SEND_MESSAGE":
			if len(args) < 3 {
				fmt.Println("Expected target peer ID and message")
				continue
			}
			targetPeerID := args[1]
			message := strings.Join(args[2:], " ")
			fmt.Printf("Sending message to peer %s: %s\n", targetPeerID, message)
			sendDataToPeer(node, targetPeerID, "", message, "", "", "")

		case "SEND_FILE":
			if len(args) < 3 {
				fmt.Println("Expected target peer ID and file path")
				continue
			}
			targetPeerID := args[1]
			filePath := args[2]
			fmt.Printf("Sending file to peer %s: %s\n", targetPeerID, filePath)
			sendDataToPeer(node, targetPeerID, filePath, "", "", "", "")
		case "SEND_DOWNLOAD_REQUEST":
			if len(args) < 3 {
				fmt.Println("Expected target peer ID and file hash")
				continue
			}
			targetPeerID := args[1]
			hash := args[2]

			// Log the test command
			fmt.Printf("Testing SEND_DOWNLOAD_REQUEST with target peer: %s and hash: %s\n", targetPeerID, hash)

			// Call the SimplyDownload function
			name, data, ext, walletAddress, err := SimplyDownload(node, targetPeerID, hash)

			if err != nil {
				fmt.Printf("Failed to send download request: %v\n", err)
				continue
			}

			// Display file information
			fmt.Printf("Download request successful:\n")
			fmt.Printf("File Name: %s\n", name)
			fmt.Printf("File Extension: %s\n", ext)
			fmt.Printf("File Data Size: %d bytes\n", len(data))

			// Display wallet address if available
			fmt.Println("Wallet Address:")
			if walletAddress != "" {
				fmt.Printf(" - Address: %s\n", walletAddress)
			} else {
				fmt.Println(" - No wallet address received.")
			}

			// Test debug: Verify global variable clearing
			fmt.Println("Testing global variable clearing...")
			dataMutex.Lock()
			defer dataMutex.Unlock()
			if receivedFileData != nil || receivedFileExt != "" || receivedFileName != "" || receivedWalletAddress != "" {
				fmt.Println("Error: Global variables were not cleared properly after the request!")
			} else {
				fmt.Println("Global variables cleared successfully.")
			}

		case "SEND_REQUEST":
			if len(args) < 4 {
				fmt.Println("Expected target peer ID, file hash, and password")
				continue
			}
			targetPeerID := args[1]
			hash := args[2]
			password := args[3]

			// Call the SendRequest function
			SendRequest(node, targetPeerID, hash, password)

		case "GET":
			if len(args) < 2 {
				fmt.Println("Expected key")
				continue
			}
			key := args[1]
			dhtKey := "/orcanet/" + key
			res, err := dht.GetValue(ctx, dhtKey)
			if err != nil {
				fmt.Printf("Failed to get record: %v\n", err)
				continue
			}
			fmt.Printf("Record: %s\n", res)

		case "GET_PROVIDERS":
			if len(args) < 2 {
				fmt.Println("Expected key")
				continue
			}
			key := args[1]

			fmt.Println("Searching for providers...")
			providerIDs, err := GetProviderIDs(node, key)
			if err != nil {
				fmt.Printf("Error getting providers: %v\n", err)
				continue
			}

			// Print the list of provider IDs
			if len(providerIDs) == 0 {
				fmt.Println("No providers found")
			} else {
				fmt.Println("Providers found:")
				for _, id := range providerIDs {
					fmt.Println(id)
				}
			}

		case "PUT":
			if len(args) < 3 {
				fmt.Println("Expected key and value")
				continue
			}
			key := args[1]
			value := args[2]
			dhtKey := "/orcanet/" + key
			log.Println(dhtKey)
			err := dht.PutValue(ctx, dhtKey, []byte(value))
			if err != nil {
				fmt.Printf("Failed to put record: %v\n", err)
				continue
			}
			// provideKey(ctx, dht, key)
			fmt.Println("Record stored successfully")

		case "PUT_PROVIDER":
			if len(args) < 2 {
				fmt.Println("Expected key")
				continue
			}
			key := args[1]
			ProvideKey(key)

			// New command handling for storing file metadata in DHT with a specified price
		case "HOST_FILE":
			if len(args) < 3 {
				fmt.Println("Expected file path and price")
				continue
			}

			filePath := args[1]
			filePrice, err := strconv.ParseFloat(args[2], 64) // Parse the price argument as a float
			if err != nil {
				fmt.Println("Invalid price, please provide a valid number")
				continue
			}

			fmt.Printf("Storing file metadata for file: %s with price: %.2f\n", filePath, filePrice)

			// Call the updated storeFileInDHT function to hash the file, get metadata, and store it
			err = storeFileInDHT(ctx, dht, filePath, filePrice)
			if err != nil {
				fmt.Printf("Failed to store file metadata: %v\n", err)
				continue
			}

			fmt.Println("File metadata stored successfully in DHT with the new provider structure.")

		default:
			fmt.Println("Expected GET, GET_PROVIDERS, PUT or PUT_PROVIDER")
		}
	}
}
