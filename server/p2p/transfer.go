package p2p

import (
	"bufio"
	"context" // for context usage
	"io"
	"log"           // for logging
	"os"            // for file operations
	"path/filepath" // for file path manipulations
	"strings"

	// Add the necessary packages from libp2p, for example:
	"github.com/libp2p/go-libp2p/core/host"    // for host.Host
	"github.com/libp2p/go-libp2p/core/network" // for network.Stream
	"github.com/libp2p/go-libp2p/core/peer"
)

func receiveDataFromPeer(node host.Host, folderPath string) {
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
		} else {
			log.Printf("Unknown header type received: %s", header)
		}
	})
}

func sendDataToPeer(node host.Host, targetPeerID, filePath, message string) {
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

	// Choose to send a message or a file based on the inputs
	if message != "" {
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
		log.Println("No file or message provided to send.")
	}
}
