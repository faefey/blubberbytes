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
	"server/database/operations"
	"strconv"
	"strings"
	"time"

	"github.com/ipfs/go-cid"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
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

func provideKey(ctx context.Context, dht *dht.IpfsDHT, key string) error {
	data := []byte(key)
	hash := sha256.Sum256(data)
	mh, err := multihash.EncodeName(hash[:], "sha2-256")
	if err != nil {
		return fmt.Errorf("error encoding multihash: %v", err)
	}
	c := cid.NewCidV1(cid.Raw, mh)

	// Start providing the key
	err = dht.Provide(ctx, c, true)
	if err != nil {
		return fmt.Errorf("failed to start providing key: %v", err)
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
		case "EXPLORE":
			explore(node)
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

		case "SEND_REQUEST":
			if len(args) < 4 {
				fmt.Println("Expected target peer ID, file hash, and password")
				continue
			}
			targetPeerID := args[1]
			hash := args[2]
			password := args[3]
			fmt.Printf("Sending request to peer %s with hash: %s\n", targetPeerID, hash)
			data, err := sendRequest(node, targetPeerID, hash, password)
			if err != nil {
				log.Fatalf("Error sending request: %v", err)
			}

			if data == nil {
				log.Println("No data received for the request")
			} else {
				log.Printf("Received data: %d bytes", len(data))
			}

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
			data := []byte(key)
			hash := sha256.Sum256(data)
			mh, err := multihash.EncodeName(hash[:], "sha2-256")
			if err != nil {
				fmt.Printf("Error encoding multihash: %v\n", err)
				continue
			}
			c := cid.NewCidV1(cid.Raw, mh)
			providers := dht.FindProvidersAsync(ctx, c, 20)

			fmt.Println("Searching for providers...")
			for p := range providers {
				if p.ID == peer.ID("") {
					break
				}
				fmt.Printf("Found provider: %s\n", p.ID.String())
				for _, addr := range p.Addrs {
					fmt.Printf(" - Address: %s\n", addr.String())
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
			provideKey(ctx, dht, key)

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
