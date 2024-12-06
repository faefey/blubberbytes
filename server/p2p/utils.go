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

func requestFileInfo(node host.Host, targetPeerID, hash string) (*models.JoinedHosting, error) {
	log.Printf("Preparing to request file info from peer %s for hash: %s", targetPeerID, hash)

	// Send the "request_info" command
	err := sendDataToPeer(node, targetPeerID, "", "", "request_info", hash, "")
	if err != nil {
		log.Printf("Failed to request file info from peer %s: %v", targetPeerID, err)
		return nil, err
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

		return &info, nil
	case <-time.After(10 * time.Second): // Timeout
		log.Println("Timeout: No response received.")
		return nil, fmt.Errorf("timed out waiting for response from peer %s", targetPeerID)
	}
}

func provideKey(key string) error {
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

func getProviderIDs(key string) ([]string, error) {
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
