package p2p

import (
	"crypto/sha256"
	"fmt"
	"log"
	"server/database/models"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multihash"
)

func RequestFileInfo(node host.Host, targetPeerID, hash string) (models.JoinedHosting, error) {
	log.Printf("Preparing to request file info from peer %s for hash: %s", targetPeerID, hash)

	// Send the "request_info" command
	err := sendDataToPeer(node, targetPeerID, "", "", "request_info", hash, "")
	if err != nil {
		log.Printf("Failed to request file info from peer %s: %v", targetPeerID, err)
		return models.JoinedHosting{}, err
	}

	log.Printf("File info request sent successfully to peer %s for hash: %s", targetPeerID, hash)

	// Wait for the infoSignal or timeout
	select {
	case <-infoSignal:
		log.Println("Signal received: File info is ready.")

		// Safely retrieve and clear the global `receivedInfo`
		dataMutex.Lock()
		info := receivedInfo
		receivedInfo = models.JoinedHosting{} // Clear the global variable
		dataMutex.Unlock()

		return info, nil
	case <-time.After(10 * time.Second): // Timeout
		log.Println("Timeout: No response received.")
		return models.JoinedHosting{}, fmt.Errorf("timed out waiting for response from peer %s", targetPeerID)
	}
}

func ProvideKey(key string) error {
	// Log the start of the provideKey process
	log.Printf("Starting to provide key: %s\n", key)
	dht := dhtRouting
	// Generate context
	ctx := globalCtx

	// Convert the key to bytes and compute the hash
	log.Printf("Converting key to bytes and hashing...")
	data := []byte(key)
	hash := sha256.Sum256(data)

	// Encode the hash as a multihash
	log.Printf("Encoding hash into multihash...")
	mh, err := multihash.EncodeName(hash[:], "sha2-256")
	if err != nil {
		log.Printf("Error encoding multihash: %v\n", err)
		return fmt.Errorf("error encoding multihash: %v", err)
	}

	// Create a CID from the multihash
	log.Printf("Creating CID from multihash...")
	c := cid.NewCidV1(cid.Raw, mh)
	log.Printf("Generated CID: %s\n", c.String())

	// Start providing the key
	log.Printf("Announcing key to the DHT...")
	err = dht.Provide(ctx, c, true)
	if err != nil {
		log.Printf("Failed to start providing key: %v\n", err)
		return fmt.Errorf("failed to start providing key: %v", err)
	}

	// Log success
	log.Printf("Successfully started providing key: %s\n", key)
	return nil
}

func GetProviderIDs(key string) ([]string, error) {
	// Assign dhtRouting to a local variable for clarity
	dht := dhtRouting

	// Check if the DHT is initialized
	if dht == nil {
		return nil, fmt.Errorf("dhtRouting is not initialized")
	}

	// Use global context
	ctx := globalCtx

	// Convert the key to a multihash
	data := []byte(key)
	hash := sha256.Sum256(data)
	mh, err := multihash.EncodeName(hash[:], "sha2-256")
	if err != nil {
		return nil, fmt.Errorf("error encoding multihash: %v", err)
	}

	// Create a CID from the multihash
	c := cid.NewCidV1(cid.Raw, mh)

	// Find providers asynchronously
	providers := dht.FindProvidersAsync(ctx, c, 20)

	// Collect provider IDs
	var ids []string
	for p := range providers {
		if p.ID == peer.ID("") {
			break
		}
		ids = append(ids, p.ID.String())
	}

	return ids, nil
}

func simply_download(node host.Host, targetPeerID, hash string) (string, []byte, string, error) {
	// Log the start of the function
	log.Printf("Starting SendDownloadRequest to peer %s for hash %s", targetPeerID, hash)

	// Call sendDataToPeer to send the request
	log.Println("Calling sendDataToPeer to send the download request...")
	err := sendDataToPeer(node, targetPeerID, "", "", "download_request", hash, "")
	if err != nil {
		log.Printf("Failed to send download request to peer %s: %v", targetPeerID, err)
		return "", nil, "", err
	}
	log.Println("Download request sent successfully. Waiting for signal...")

	// Wait for the first signal
	<-signalChan
	log.Println("First signal received. Proceeding...")

	time.Sleep(500 * time.Millisecond)

	// Check for hash signal
	select {
	case <-hashSignalChan: // Replace with your actual hash signal channel
		log.Println("Received hash signal indicating the hash is invalid.")
		return "", nil, "", fmt.Errorf("hash is invalid")
	case <-time.After(100 * time.Millisecond):
		log.Println("No hash signal received within 100ms. Continuing...")
	}

	// Lock the data mutex
	log.Println("Acquiring lock to access global variables...")
	dataMutex.Lock()
	defer dataMutex.Unlock()
	log.Println("Global variables locked. Checking received data...")

	if receivedFileData == nil || receivedFileExt == "" || receivedFileName == "" {
		log.Println("File data, name, or extension is missing in the received data.")
		return "", nil, "", fmt.Errorf("file data, name, or extension is missing")
	}

	// Retrieve the file name, data, and extension
	log.Printf("Received file details:\n - Name: %s\n - Extension: %s\n - Data Size: %d bytes", receivedFileName, receivedFileExt, len(receivedFileData))
	name := receivedFileName
	data := receivedFileData
	ext := receivedFileExt

	// Clear the global variables
	log.Println("Clearing global variables for the next request...")
	receivedFileData = nil
	receivedFileExt = ""
	receivedFileName = ""

	// Log success and return the results
	log.Println("SendDownloadRequest completed successfully.")
	return name, data, ext, nil
}

func SendRequest(node host.Host, targetPeerID, hash, password string) (string, []byte, string, error) {
	// Call sendDataToPeer to send the request
	err := sendDataToPeer(node, targetPeerID, "", "", "request", hash, password)
	if err != nil {
		return "", nil, "", err
	}

	<-signalChan // Wait for the first signal

	time.Sleep(500 * time.Millisecond)

	// Check for hash signal
	select {
	case <-hashSignalChan: // Replace with your actual hash signal channel
		return "", nil, "", fmt.Errorf("hash is invalid")
	case <-time.After(100 * time.Millisecond):
		// No hash signal received, continue
	}

	// Check for password signal
	select {
	case <-passwordSignalChan: // Replace with your actual password signal channel
		return "", nil, "", fmt.Errorf("password is invalid")
	case <-time.After(100 * time.Millisecond):
		// No password signal received, continue
	}

	dataMutex.Lock() // Lock the mutex to safely access the global variables
	defer dataMutex.Unlock()

	if receivedFileData == nil || receivedFileExt == "" || receivedFileName == "" {
		return "", nil, "", fmt.Errorf("file data, name, or extension is missing")
	}

	// Retrieve the file name, data, and extension
	name := receivedFileName // Copy the file name
	data := receivedFileData // Copy the data
	ext := receivedFileExt   // Copy the file extension

	// Clear the global variables
	receivedFileData = nil
	receivedFileExt = ""
	receivedFileName = ""

	return name, data, ext, nil
}
