package p2p

import (
    "context"        // for context usage
    "log"            // for logging
    "os"             // for file operations
    "path/filepath"  // for file path manipulations
	"io"

    // Add the necessary packages from libp2p, for example:
    "github.com/libp2p/go-libp2p/core/host"   // for host.Host
    "github.com/libp2p/go-libp2p/core/network" // for network.Stream
	"github.com/libp2p/go-libp2p/core/peer"

)



func receiveDataFromPeer(node host.Host, folderPath string) {
	log.Println("Setting up stream handler to listen for incoming streams on '/senddata/p2p' protocol.")

	// Set a stream handler to listen for incoming streams on the "/senddata/p2p" protocol
	node.SetStreamHandler("/senddata/p2p", func(s network.Stream) {
		log.Printf("New stream opened from peer: %s", s.Conn().RemotePeer())
		defer func() {
			log.Printf("Stream closed by peer: %s", s.Conn().RemotePeer())
			s.Close()
		}()

		// Generate a unique file name for each incoming file based on timestamp (no extension for now)
		fileName := "node_file.pdf"
		filePath := filepath.Join(folderPath, fileName)

		// Log the attempt to create the file
		log.Printf("Attempting to create file at path: %s", filePath)

		// Open the file to store the received data
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Failed to create file in folder %s: %v", folderPath, err)
			return
		}
		defer func() {
			log.Printf("File closed: %s", filePath)
			file.Close()
		}()

		log.Printf("Writing received data to file: %s", filePath)

		// Read the entire stream content and write it to the file
		data, err := io.ReadAll(s)
		if err != nil {
			log.Printf("Error reading from stream from peer %s: %v", s.Conn().RemotePeer(), err)
			return
		}

		// Write the data to the file
		n, err := file.Write(data)
		if err != nil {
			log.Printf("Error writing to file %s: %v", filePath, err)
			return
		}

		log.Printf("File writing complete. Total bytes written: %d to file: %s", n, filePath)
	})

	log.Println("Listening for incoming streams on '/senddata/p2p' protocol.")
}








func sendDataToPeer(node host.Host, targetPeerID string, filePath string) {
	ctx := context.Background()

	// Decode the target peer ID
	targetPeerIDParsed, err := peer.Decode(targetPeerID)
	if err != nil {
		log.Printf("Failed to decode target peer ID: %v", err)
		return
	}

	log.Printf("Target peer ID successfully decoded: %s", targetPeerIDParsed)

	// Log file path and check if the file exists
	log.Printf("Attempting to send file at path: %s to peer %s", filePath, targetPeerIDParsed)

	// Open the file to read its content
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return
	}
	defer func() {
		log.Printf("File closed after sending: %s", filePath)
		file.Close()
	}()

	// Log file size and name
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Could not retrieve file info for: %s", filePath)
		return
	}
	log.Printf("File Info - Name: %s, Size: %d bytes", fileInfo.Name(), fileInfo.Size())

	// Attempt to open a stream to the target peer
	s, err := node.NewStream(network.WithAllowLimitedConn(ctx, "/senddata/p2p"), targetPeerIDParsed, "/senddata/p2p")
	if err != nil {
		log.Printf("Failed to open stream to %s: %v", targetPeerIDParsed, err)
		return
	}
	defer func() {
		log.Printf("Closing stream to peer %s", targetPeerIDParsed)
		s.Close()
	}()

	log.Printf("Stream opened successfully to peer %s", targetPeerIDParsed)

	// Read the entire file content and send it
	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file content: %v", err)
		return
	}

	// Send the entire file content over the stream
	n, err := s.Write(fileContent)
	if err != nil {
		log.Printf("Failed to send file content to peer %s: %v", targetPeerIDParsed, err)
		return
	}

	log.Printf("File sent successfully. Total bytes sent: %d to peer %s", n, targetPeerIDParsed)
}