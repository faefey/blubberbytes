package p2p

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"server/database/models"
	"server/database/operations"
	"time"

	"math/rand"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/host"
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

func GetProviderIDs(node host.Host, key string) ([]string, error) {
	// Assign dhtRouting to a local variable for clarity
	dht := dhtRouting

	// Check if the DHT is initialized
	if dht == nil {
		return []string{}, fmt.Errorf("dhtRouting is not initialized")
	}

	// Use global context
	ctx := globalCtx

	// Convert the key to a multihash
	data := []byte(key)
	hash := sha256.Sum256(data)
	mh, err := multihash.EncodeName(hash[:], "sha2-256")
	if err != nil {
		return []string{}, fmt.Errorf("error encoding multihash: %v", err)
	}

	// Create a CID from the multihash
	c := cid.NewCidV1(cid.Raw, mh)

	// Find providers asynchronously
	providers := dht.FindProvidersAsync(ctx, c, 20)

	// Collect provider IDs
	var ids []string
	for p := range providers {

		// Skip the node's own PeerID
		if p.ID == node.ID() {
			continue
		}

		ids = append(ids, p.ID.String())
	}

	return ids, nil
}

func SimplyDownload(node host.Host, targetPeerID, hash string) (string, []byte, string, string, error) {
	// Log the start of the function
	log.Printf("Starting SendDownloadRequest to peer %s for hash %s", targetPeerID, hash)

	// Call sendDataToPeer to send the download request
	log.Println("Calling sendDataToPeer to send the download request...")
	err := sendDataToPeer(node, targetPeerID, "", "", "download_request", hash, "")
	if err != nil {
		log.Printf("Failed to send download request to peer %s: %v", targetPeerID, err)
		return "", nil, "", "", err
	}
	log.Println("Download request sent successfully. Waiting for signal...")

	// Wait for the first signal
	<-signalChan
	log.Println("First signal received. Proceeding...")

	time.Sleep(500 * time.Millisecond)

	// Check for hash signal
	select {
	case <-hashSignalChan:
		log.Println("Received hash signal indicating the hash is invalid.")
		return "", nil, "", "", fmt.Errorf("hash is invalid")
	case <-time.After(100 * time.Millisecond):
		log.Println("No hash signal received within 100ms. Continuing...")
	}

	// Lock the data mutex
	log.Println("Acquiring lock to access global variables...")
	dataMutex.Lock()
	defer dataMutex.Unlock()
	log.Println("Global variables locked. Checking received data...")

	if receivedFileData == nil || receivedFileExt == "" || receivedFileName == "" || receivedWalletAddress == "" {
		log.Println("File data, name, extension, or wallet address is missing in the received data.")
		return "", nil, "", "", fmt.Errorf("file data, name, extension, or wallet address is missing")
	}

	// Retrieve the file name, data, and extension
	log.Printf("Received file details:\n - Name: %s\n - Extension: %s\n - Data Size: %d bytes", receivedFileName, receivedFileExt, len(receivedFileData))
	name := receivedFileName
	data := receivedFileData
	ext := receivedFileExt

	// Retrieve wallet address directly
	walletAddress := receivedWalletAddress
	log.Printf("Retrieved wallet address: %s", walletAddress)

	// Clear the global variables for the next request
	log.Println("Clearing global variables for the next request...")
	receivedFileData = nil
	receivedFileExt = ""
	receivedFileName = ""
	receivedWalletAddress = "" // Clear wallet address

	return name, data, ext, walletAddress, nil
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

func RandomProxiesInfo(node host.Host) ([]models.Proxy, error) {
	// Get a list of provider IDs for the "PROXY" key from the DHT
	providerIDs, err := GetProviderIDs(node, "PROXY")
	if err != nil {
		log.Printf("Failed to get provider IDs for PROXY key: %v", err)
		return []models.Proxy{}, err
	}

	// Log the original list of provider IDs
	log.Println("Original provider IDs:", providerIDs)

	// Shuffle the list of provider IDs using a random generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(providerIDs), func(i, j int) {
		providerIDs[i], providerIDs[j] = providerIDs[j], providerIDs[i]
	})

	// Log the shuffled list
	log.Println("Shuffled provider IDs:", providerIDs)

	// Select up to 5 random providers
	selectedProviders := providerIDs
	if len(providerIDs) > 5 {
		selectedProviders = providerIDs[:5]
	}

	// Log the selected providers
	log.Println("Selected providers:", selectedProviders)

	// Send proxy requests and wait for signals
	for _, targetPeerID := range selectedProviders {
		// Send the proxy request
		err := sendDataToPeer(node, targetPeerID, "", "", "proxy_request", "", "")
		if err != nil {
			log.Printf("Failed to send proxy request to peer %s: %v", targetPeerID, err)
			// Continue to the next peer even if one fails
			continue
		}

		// Wait for a signal indicating the proxy response was processed
		select {
		case <-proxySignal:
			log.Printf("Signal received after sending proxy request to peer %s. Proceeding to next peer.", targetPeerID)
		case <-time.After(5 * time.Second): // Timeout after 5 seconds
			log.Printf("Timeout waiting for response signal from peer %s. Skipping to next peer.", targetPeerID)
		}
	}

	// Read, log, and clear the global list of proxies
	dataMutex.Lock()
	defer dataMutex.Unlock()

	log.Printf("Returning global proxy list: %+v", proxyList)
	result := make([]models.Proxy, len(proxyList))
	copy(result, proxyList) // Create a copy of the proxy list to return
	proxyList = nil         // Clear the global list
	log.Println("Global proxy list cleared.")

	return result, nil
}

func Explore(node host.Host, peerIDs []string) ([]models.JoinedHosting, error) {
	// Iterate through the list of peer IDs
	for _, peerID := range peerIDs {
		log.Printf("Requesting all files from peer: %s", peerID)
		if peerID == node.ID().String() {
			log.Printf("Skipping request to self for peer ID: %s", peerID)
			continue
		}
		// Send a generic "request all files" signal to the peer
		err := sendDataToPeer(node, peerID, "", "", "request_all", "", "")
		if err != nil {
			log.Printf("Error requesting all files from peer %s: %v", peerID, err)
			continue
		}

		log.Printf("Request sent to peer %s for all files. Waiting for response signal...", peerID)

		// Wait for a signal (blocking until the signal is received)
		<-hostingUpdateSignal
		log.Printf("Response signal received from peer %s", peerID)
	}

	// After processing all peers, collect and clear the global hosting list
	collectedHostings := hostingList
	hostingList = []models.JoinedHosting{} // Clear the global hosting list

	// Log and return the collected hostings
	log.Printf("Total collected hostings: %d", len(collectedHostings))
	return collectedHostings, nil
}

func SendProxyBillWithConfirmation(node host.Host, peerID string, proxyBill models.ProxyBill) error {
	// Serialize ProxyBill to JSON
	proxyBillJSON, err := json.Marshal(proxyBill)
	if err != nil {
		log.Printf("Failed to serialize ProxyBill: %v", err)
		return fmt.Errorf("failed to serialize ProxyBill: %w", err)
	}

	// Send the ProxyBill to the specified peer
	log.Printf("Sending ProxyBill to peer %s", peerID)
	err = sendDataToPeer(node, peerID, "", string(proxyBillJSON), "ProxyBill", "", "")
	if err != nil {
		log.Printf("Failed to send ProxyBill to peer: %v", err)
		return fmt.Errorf("failed to send ProxyBill to peer: %w", err)
	}

	// Wait for confirmation signal
	select {
	case <-successSignal:
		log.Println("ProxyBill processing confirmed as successful")
		return nil // No error, successful transaction

	case <-failureSignal:
		log.Println("ProxyBill processing confirmed as unsuccessful")
		return fmt.Errorf("proxyBill processing failed")

	case <-time.After(10 * time.Second): // Timeout for waiting confirmation
		log.Println("Timed out waiting for confirmation signal")
		return fmt.Errorf("confirmation signal timeout")
	}
}

func handleProxyBill(node host.Host, proxyBill models.ProxyBill, peerID string, btcwallet *rpcclient.Client, netParams *chaincfg.Params, db *sql.DB) error {
	log.Println("Received ProxyBill:")
	log.Printf("IP: %s", proxyBill.IP)
	log.Printf("Rate: %.2f", proxyBill.Rate)
	log.Printf("Bytes: %d", proxyBill.Bytes)
	log.Printf("Amount: %.2f", proxyBill.Amount)
	log.Printf("Wallet: %s", proxyBill.Wallet)

	// Placeholder function for additional processing
	err := processProxyBill(proxyBill, btcwallet, netParams, db)
	if err != nil {
		log.Printf("Failed to process ProxyBill: %v", err)

		// Send failure confirmation back to peer
		err = sendDataToPeer(node, peerID, "", "Processing failed", "confirmation", "", "")
		if err != nil {
			log.Printf("Failed to send failure confirmation to peer: %v", err)
		}
		return err
	}

	// Send success confirmation back to peer
	err = sendDataToPeer(node, peerID, "", "Processing successful", "confirmation", "", "")
	if err != nil {
		log.Printf("Failed to send success confirmation to peer: %v", err)
		return err
	}

	log.Println("Successfully processed ProxyBill and sent confirmation.")
	return nil
}

func processProxyBill(proxyBill models.ProxyBill, btcwallet *rpcclient.Client, netParams *chaincfg.Params, db *sql.DB) error {
	log.Println("Processing ProxyBill...")

	if proxyBill.Rate == -1 {
		err := operations.AddIPtoNode(db, proxyBill.IP, proxyBill.Wallet)
		if err != nil {
			return err
		}
	} else {
		walletInfo, err := operations.GetWalletInfo(db)
		if err != nil {
			return err
		}

		err = btcwallet.WalletPassphrase(walletInfo.PrivPassphrase, 300)
		if err != nil {
			return err
		}

		btcutilAddress, err := btcutil.DecodeAddress(proxyBill.Wallet, netParams)
		if err != nil {
			return err
		}

		_, err = btcwallet.SendToAddress(btcutilAddress, btcutil.Amount(proxyBill.Amount*1e8))
		if err != nil {
			return err
		}
	}

	// Simulate a processing error
	err := fmt.Errorf("simulated processing error for ProxyBill with IP: %s", proxyBill.IP)
	log.Printf("Error: %v", err)
	return err
}
