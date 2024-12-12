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

	// Add the necessary packages from libp2p, for example:
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/libp2p/go-libp2p/core/host"    // for host.Host
	"github.com/libp2p/go-libp2p/core/network" // for network.Stream
	"github.com/libp2p/go-libp2p/core/peer"
)

var (
	storing               []models.Storing // Global variable to hold Storing objects
	storingMutex          sync.Mutex       // Mutex to ensure thread-safe access to the global variable
	receivedFileData      []byte
	receivedFileExt       string
	receivedFileName      string
	receivedInfo          models.JoinedHosting
	receivedWalletAddress string
	infoSignal            = make(chan struct{})
	passwordSignalChan    = make(chan struct{})
	hashSignalChan        = make(chan struct{})
	hostingUpdateSignal   = make(chan struct{})
	successSignal         = make(chan struct{})
	failureSignal         = make(chan struct{})
	proxyList             []models.Proxy        // Global list to store received proxies
	proxySignal           = make(chan struct{}) // Channel to signal when a response is received
	hostingList           []models.JoinedHosting
	dataMutex             sync.Mutex
)

// Channel for signaling when data is ready
var signalChan = make(chan struct{}, 1)

func receiveDataFromPeer(node host.Host, db *sql.DB, folderPath string, btcwallet *rpcclient.Client, netParams *chaincfg.Params) {
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

		} else if header == "confirmation" {
			// Handle confirmation message
			log.Printf("Processing confirmation from peer: %s", s.Conn().RemotePeer())

			// Read the confirmation message
			message, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error reading confirmation message from peer %s: %v", s.Conn().RemotePeer(), err)
				return
			}
			message = strings.TrimSpace(message)

			// Log the confirmation message
			log.Printf("Confirmation received: %s", message)

			// Send a signal based on the confirmation message
			if message == "Processing successful" {
				log.Println("Processing was successful. Sending success signal.")
				successSignal <- struct{}{}
			} else if message == "Processing failed" {
				log.Println("Processing failed. Sending failure signal.")
				failureSignal <- struct{}{}
			} else {
				log.Printf("Unknown confirmation message received: %s", message)
			}
		} else if header == "ProxyBill" {
			log.Printf("Processing 'ProxyBill' from peer: %s", s.Conn().RemotePeer())

			// Read JSON data for ProxyBill
			data, err := io.ReadAll(reader)
			if err != nil {
				log.Printf("Error reading ProxyBill data from peer %s: %v", s.Conn().RemotePeer(), err)
				return
			}

			// Unmarshal JSON into ProxyBill struct
			var proxyBill models.ProxyBill
			err = json.Unmarshal(data, &proxyBill)
			if err != nil {
				log.Printf("Error unmarshaling ProxyBill data from peer %s: %v", s.Conn().RemotePeer(), err)
				log.Printf("Received data was: %s", string(data))
				return
			}

			// Call helper function to handle the ProxyBill and send confirmation back
			err = handleProxyBill(node, proxyBill, s.Conn().RemotePeer().String(), btcwallet, netParams, db)
			if err != nil {
				log.Printf("Error processing ProxyBill: %v", err)
			}
		} else if header == "requested_proxy" {
			log.Printf("Processing 'requested_proxy' response from peer: %s", s.Conn().RemotePeer())

			// Log the start of reading the response
			log.Println("Attempting to read the 'requested_proxy' response.")

			// Read the next line to check for the response type
			response, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error reading 'requested_proxy' response from peer %s: %v", s.Conn().RemotePeer(), err)
				log.Println("Possible reasons: Peer closed the stream prematurely or did not send data.")
				log.Println("Ensure the sending peer is properly sending the 'requested_proxy' response.")
				// Send a signal even if there's an error
				proxySignal <- struct{}{}
				return
			}
			response = strings.TrimSpace(response)

			// Log the raw response received
			log.Printf("Raw response received: %s", response)

			// Check the response
			if response == "no proxy anymore" {
				log.Println("No proxy available. Peer has no proxies to provide.")
				// Send a signal and exit
				proxySignal <- struct{}{}
				return
			}

			// Log that a proxy response was received
			log.Println("Received proxy data, attempting to unmarshal JSON.")

			// If a proxy is available, unmarshal the JSON data
			var proxy models.Proxy
			err = json.Unmarshal([]byte(response), &proxy)
			if err != nil {
				log.Printf("Error unmarshaling proxy data from peer %s: %v", s.Conn().RemotePeer(), err)
				log.Printf("Received data was: %s", response)
				log.Println("Ensure the response is valid JSON in the format: {\"IP\":\"<ip>\",\"Rate\":<rate>},\"Wallet\":<wallet>}")
				// Send a signal even if there's an unmarshaling error
				proxySignal <- struct{}{}
				return
			}

			// Log the received proxy data
			log.Printf("Received proxy from peer: %+v", proxy)

			// Add the received proxy to the global list
			dataMutex.Lock()
			proxyList = append(proxyList, proxy)
			dataMutex.Unlock()

			log.Printf("Proxy added to global list. Current list size: %d", len(proxyList))

			// Send a signal after adding the proxy
			proxySignal <- struct{}{}
		} else if header == "proxy_request" {
			log.Printf("Processing 'proxy_request' request from peer: %s", s.Conn().RemotePeer())

			// Use the helper function to handle the proxy response
			err := sendProxyResponseToPeer(node, s.Conn().RemotePeer().String(), db)
			if err != nil {
				log.Printf("Error processing 'proxy_request': %v", err)
			}
		} else if header == "download_request" {
			handleDownloadRequest(s, db, node, s.Conn().RemotePeer().String())
		} else if header == "request_info" {
			handleInfoRequest(s, db, node)
		} else if header == "requested_info" {
			// Handle received info
			log.Printf("Receiving requested info from peer: %s", s.Conn().RemotePeer())

			// Read JSON data
			data, err := io.ReadAll(reader)
			if err != nil {
				log.Printf("Error reading requested info from peer %s: %v", s.Conn().RemotePeer(), err)
				return
			}

			// Parse JSON into JoinedHosting struct
			var info models.JoinedHosting
			err = json.Unmarshal(data, &info)
			if err != nil {
				log.Printf("Error unmarshaling requested info: %v", err)
				return
			}

			// Store the received info globally
			dataMutex.Lock()
			receivedInfo = info
			dataMutex.Unlock()
			infoSignal <- struct{}{}
		} else if header == "message" {
			// Handle message transfer
			message, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error reading message from stream: %v", err)
				return
			}

			// Trim the message to avoid unnecessary whitespaces or newlines
			message = strings.TrimSpace(message)

			// Handle specific messages
			switch message {
			case "Invalid password":
				log.Println("Received 'Invalid password' message from peer.")
				signalChan <- struct{}{}
				passwordSignalChan <- struct{}{} // Notify the password signal channel
				return

			case "File not found":
				log.Println("Received 'File not found' message from peer.")
				signalChan <- struct{}{}
				hashSignalChan <- struct{}{} // Notify the file signal channel
				return

			default:
				log.Printf("Received unknown message from peer: %s", message)
				return
			}
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
			signalChan <- struct{}{}
		} else if header == "request_all" {
			log.Printf("Received 'send_all' request from peer: %s", s.Conn().RemotePeer())
			handleSendAllRequest(s, db, node, s.Conn().RemotePeer().String())
		} else if header == "requested_hostings" {
			log.Printf("Handling 'requested_hostings' from peer: %s", s.Conn().RemotePeer())

			// Read JSON data
			data, err := io.ReadAll(reader)
			if err != nil {
				log.Printf("Error reading 'requested_hostings' data: %v", err)
				return
			}

			// Parse JSON into a slice of JoinedHosting objects
			var receivedHostings []models.JoinedHosting
			err = json.Unmarshal(data, &receivedHostings)
			if err != nil {
				log.Printf("Error unmarshalling 'requested_hostings' data: %v", err)
				return
			}

			// Safely add the received hostings to the global hosting list
			dataMutex.Lock()
			hostingList = append(hostingList, receivedHostings...)
			dataMutex.Unlock()

			log.Printf("Successfully added %d hostings to the global hosting list", len(receivedHostings))

			hostingUpdateSignal <- struct{}{}
			log.Println("Signal sent for hosting updates")

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

		} else if header == "requested_file_name" {
			// Handle file name transfer
			log.Printf("Handling file name transfer from peer: %s", s.Conn().RemotePeer())

			name, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error reading file name from stream: %v", err)
				return
			}
			name = strings.TrimSpace(name)

			// Safely store the received name
			dataMutex.Lock()
			receivedFileName = name
			dataMutex.Unlock()

			log.Printf("File name received and stored: %s", receivedFileName)
		} else if header == "requested_wallet_address" {
			log.Printf("Handling wallet address transfer from peer: %s", s.Conn().RemotePeer())

			// Read the wallet address string
			data, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error reading wallet address from stream: %v", err)
				return
			}
			data = strings.TrimSpace(data)

			// Check for "No wallet address available" response
			if data == "No wallet address available" {
				log.Println("No wallet address available from peer.")
				return
			}

			// Safely store the received wallet address
			dataMutex.Lock()
			receivedWalletAddress = data
			dataMutex.Unlock()

			log.Printf("Wallet address received and stored: %s", receivedWalletAddress)
		} else {
			log.Printf("Unknown header type received: %s", header)
		}
	})
}

func sendProxyResponseToPeer(node host.Host, targetPeerID string, db *sql.DB) error {
	log.Printf("Preparing to send proxy response to peer %s", targetPeerID)

	// Decode the target peer ID
	targetPeerIDParsed, err := peer.Decode(targetPeerID)
	if err != nil {
		log.Printf("Failed to decode target peer ID: %v", err)
		return err
	}

	// Open a new stream to the target peer
	ctx := context.Background()
	s, err := node.NewStream(network.WithAllowLimitedConn(ctx, "/senddata/p2p"), targetPeerIDParsed, "/senddata/p2p")
	if err != nil {
		log.Printf("Failed to open stream to peer %s: %v", targetPeerIDParsed, err)
		return err
	}
	defer func() {
		log.Printf("Closing stream to peer %s", targetPeerIDParsed)
		s.Close()
	}()
	log.Printf("Stream opened successfully to peer %s", targetPeerIDParsed)

	// Send the "requested_proxy" header
	_, err = s.Write([]byte("requested_proxy\n"))
	if err != nil {
		log.Printf("Failed to send 'requested_proxy' header to peer %s: %v", targetPeerIDParsed, err)
		return err
	}

	// Retrieve the proxy from the database
	proxy, err := operations.GetProxy(db)
	if err != nil {
		// Send "no proxy anymore" if there's a database error
		_, _ = s.Write([]byte("no proxy anymore\n"))
		log.Printf("Error retrieving proxy from database: %v", err)
		return err
	}

	if proxy == nil {
		// No proxy found, send "no proxy anymore"
		_, err = s.Write([]byte("no proxy anymore\n"))
		if err != nil {
			log.Printf("Error sending 'no proxy anymore' message to peer %s: %v", targetPeerIDParsed, err)
			return err
		}
		log.Println("Sent 'no proxy anymore' message.")
		return nil
	}

	// Proxy found, send it back as JSON
	proxyData, err := json.Marshal(proxy)
	if err != nil {
		// Send "no proxy anymore" if JSON marshaling fails
		_, _ = s.Write([]byte("no proxy anymore\n"))
		log.Printf("Error marshaling proxy data: %v", err)
		return err
	}

	_, err = s.Write(append(proxyData, '\n'))
	if err != nil {
		log.Printf("Error sending proxy data to peer %s: %v", targetPeerIDParsed, err)
		return err
	}

	log.Printf("Successfully sent proxy data to peer %s: %+v", targetPeerIDParsed, proxy)
	return nil
}

func handleDownloadRequest(s network.Stream, db *sql.DB, node host.Host, targetPeerID string) {
	log.Printf("Handling download request from peer %s", targetPeerID)

	reader := bufio.NewReader(s)

	// Read the file hash
	fileHash, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading file hash from stream from peer %s: %v", targetPeerID, err)
		sendDataToPeer(node, targetPeerID, "", "File not found", "message", "", "")
		return
	}
	fileHash = strings.TrimSpace(fileHash)
	log.Printf("Received file hash: %s", fileHash)

	// Retrieve file metadata from the database
	log.Printf("Searching for file metadata in the database for hash: %s", fileHash)
	storing, err := operations.FindStoring(db, fileHash)
	if err != nil || storing == nil {
		log.Printf("File not found or error occurred while fetching file metadata for hash %s: %v", fileHash, err)
		sendDataToPeer(node, targetPeerID, "", "File not found", "message", "", "")
		return
	}

	log.Printf("Found file metadata for file hash: %s", fileHash)

	// Send the file name
	fileName := storing.Name
	err = sendRequestedFileNameToPeer(node, targetPeerID, fileName)
	if err != nil {
		log.Printf("Error sending file name to peer %s: %v", targetPeerID, err)
		return
	}
	log.Printf("File name sent successfully to peer %s: %s", targetPeerID, fileName)

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

	// Send wallet address
	err = sendWalletAddressToPeer(node, targetPeerID, db)
	if err != nil {
		log.Printf("Error sending wallet address to peer %s: %v", targetPeerID, err)
		return
	}
	log.Printf("Wallet address sent successfully to peer %s", targetPeerID)

	// Use sendDataToPeer to send the requested file back
	log.Printf("Sending requested file back to peer %s from path: %s", targetPeerID, storing.Path)
	err = sendRequestedFileToPeer(node, targetPeerID, storing.Path)
	if err != nil {
		log.Printf("Error sending requested file to peer %s: %v", targetPeerID, err)
		return
	}

	log.Printf("File sent successfully to peer %s: %s", targetPeerID, storing.Path)
}

func sendWalletAddressToPeer(node host.Host, targetPeerID string, db *sql.DB) error {
	log.Printf("Preparing to send wallet address to peer %s", targetPeerID)

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

	// Write the "requested_wallet_address" header
	_, err = s.Write([]byte("requested_wallet_address\n"))
	if err != nil {
		log.Printf("Failed to send requested_wallet_address header to peer %s: %v", targetPeerIDParsed, err)
		return err
	}
	log.Printf("Sent 'requested_wallet_address' header to peer %s", targetPeerID)

	// Retrieve wallet info from the database
	walletInfo, err := operations.GetWalletInfo(db)
	if err != nil {
		log.Printf("Error retrieving wallet info from database: %v", err)
		return err
	}
	if walletInfo == nil || walletInfo.Address == "" {
		log.Printf("No wallet address found in the database.")
		_, _ = s.Write([]byte("No wallet address available\n"))
		return nil
	}

	// Send the wallet address as plain text
	_, err = s.Write([]byte(walletInfo.Address + "\n"))
	if err != nil {
		log.Printf("Error sending wallet address to peer %s: %v", targetPeerIDParsed, err)
		return err
	}
	log.Printf("Wallet address sent to peer %s: %s", targetPeerID, walletInfo.Address)

	return nil
}

func handleSendAllRequest(s network.Stream, db *sql.DB, node host.Host, targetPeerID string) {
	log.Printf("Handling 'send_all' request for peer: %s", targetPeerID)

	// Retrieve all hosting records from the database
	hostingRecords, err := operations.GetAllHosting(db)
	if err != nil {
		log.Printf("Error retrieving hosting records: %v", err)
		return // Exit the function if retrieval fails
	}

	// Decode the target peer ID
	targetPeerIDParsed, err := peer.Decode(targetPeerID)
	if err != nil {
		log.Printf("Failed to decode target peer ID: %v", err)
		return
	}

	// Open a stream to the target peer
	ctx := context.Background()
	stream, err := node.NewStream(network.WithAllowLimitedConn(ctx, "/senddata/p2p"), targetPeerIDParsed, "/senddata/p2p")
	if err != nil {
		log.Printf("Failed to open stream to peer %s: %v", targetPeerIDParsed, err)
		return
	}
	defer stream.Close()
	log.Printf("Stream opened successfully to peer %s", targetPeerIDParsed)

	// Send the header to indicate the type of data being sent
	header := "requested_hostings\n"
	_, err = stream.Write([]byte(header))
	if err != nil {
		log.Printf("Error sending header to peer %s: %v", targetPeerIDParsed, err)
		return
	}

	// Serialize the hosting records to JSON
	jsonData, err := json.Marshal(hostingRecords)
	if err != nil {
		log.Printf("Error serializing hosting records to JSON: %v", err)
		return // Exit the function if serialization fails
	}

	// Send the JSON data back to the requesting peer
	_, err = stream.Write(jsonData)
	if err != nil {
		log.Printf("Error sending hosting records to peer %s: %v", targetPeerIDParsed, err)
		return
	}

	log.Printf("All hosting records sent successfully to peer: %s", targetPeerIDParsed)
}

func handleFileRequest(s network.Stream, db *sql.DB, node host.Host, targetPeerID string) {
	log.Printf("Handling file request from peer %s", targetPeerID)

	reader := bufio.NewReader(s)

	// Read the file hash
	fileHash, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading file hash from stream from peer %s: %v", targetPeerID, err)
		sendDataToPeer(node, targetPeerID, "", "File not found", "", "", "")
		return
	}
	fileHash = strings.TrimSpace(fileHash)
	log.Printf("Received file hash: %s", fileHash)

	// Read the password
	password, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading password from stream from peer %s: %v", targetPeerID, err)
		sendDataToPeer(node, targetPeerID, "", "Invalid password", "", "", "")
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

	// Send the file name
	fileName := storing.Name
	err = sendRequestedFileNameToPeer(node, targetPeerID, fileName)
	if err != nil {
		log.Printf("Error sending file name to peer %s: %v", targetPeerID, err)
		return
	}
	log.Printf("File name sent successfully to peer %s: %s", targetPeerID, fileName)

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

	// Use sendDataToPeer to send the requested file back
	log.Printf("Sending requested file back to peer %s from path: %s", targetPeerID, storing.Path)
	err = sendRequestedFileToPeer(node, targetPeerID, storing.Path)
	if err != nil {
		log.Printf("Error sending requested file to peer %s: %v", targetPeerID, err)
		return
	}

	log.Printf("File sent successfully to peer %s: %s", targetPeerID, storing.Path)

}

func sendRequestedFileNameToPeer(node host.Host, targetPeerID, fileName string) error {
	log.Printf("Preparing to send file name to peer %s, name: %s", targetPeerID, fileName)

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

	// Write the "requested_file_name" header
	_, err = s.Write([]byte("requested_file_name\n"))
	if err != nil {
		log.Printf("Failed to send requested_file_name header to peer %s: %v", targetPeerIDParsed, err)
		return err
	}
	log.Printf("Sent 'requested_file_name' header to peer %s", targetPeerID)

	// Write the file name
	_, err = s.Write([]byte(fileName + "\n"))
	if err != nil {
		log.Printf("Failed to send file name to peer %s: %v", targetPeerIDParsed, err)
		return err
	}
	log.Printf("Sent file name to peer %s: %s", targetPeerID, fileName)

	return nil
}

func sendDataToPeer(node host.Host, targetPeerID, filePath, message, dataType string, hash string, password string) error {
	connectToPeerUsingRelay(node, targetPeerID)
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

	} else if dataType == "confirmation" {
		log.Printf("Sending confirmation to peer %s", targetPeerIDParsed)

		// Send "confirmation" header
		_, err = s.Write([]byte("confirmation\n"))
		if err != nil {
			log.Printf("Failed to send confirmation header to peer %s: %v", targetPeerIDParsed, err)
			return err
		}

		// Write the confirmation message
		_, err = s.Write([]byte(message + "\n"))
		if err != nil {
			log.Printf("Failed to send confirmation message to peer %s: %v", targetPeerIDParsed, err)
			return err
		}

		log.Printf("Confirmation sent successfully to peer %s with message: %s", targetPeerIDParsed, message)
	} else if dataType == "ProxyBill" {
		log.Printf("Sending ProxyBill to peer %s", targetPeerIDParsed)

		// Convert message to JSON
		var proxyBill models.ProxyBill
		err = json.Unmarshal([]byte(message), &proxyBill)
		if err != nil {
			log.Printf("Failed to parse ProxyBill JSON: %v", err)
			return err
		}

		// Serialize ProxyBill to JSON string
		proxyBillJSON, err := json.Marshal(proxyBill)
		if err != nil {
			log.Printf("Failed to serialize ProxyBill to JSON: %v", err)
			return err
		}

		// Send "ProxyBill" header
		_, err = s.Write([]byte("ProxyBill\n"))
		if err != nil {
			log.Printf("Failed to send ProxyBill header to peer %s: %v", targetPeerIDParsed, err)
			return err
		}

		// Send the serialized ProxyBill JSON
		_, err = s.Write(proxyBillJSON)
		if err != nil {
			log.Printf("Failed to send ProxyBill data to peer %s: %v", targetPeerIDParsed, err)
			return err
		}

		log.Printf("ProxyBill sent successfully to peer %s", targetPeerIDParsed)
	} else if dataType == "proxy_request" {
		_, err = s.Write([]byte("proxy_request\n"))

	} else if dataType == "download_request" {
		// Send a "download_request" header
		log.Printf("Sending download request to peer %s for hash: %s", targetPeerIDParsed, hash)
		_, err = s.Write([]byte("download_request\n"))
		if err != nil {
			log.Printf("Failed to send download_request header to peer %s: %v", targetPeerIDParsed, err)
			return err
		}

		// Write the hash of the file to download
		_, err = s.Write([]byte(hash + "\n"))
		if err != nil {
			log.Printf("Failed to send hash for download to peer %s: %v", targetPeerIDParsed, err)
			return err
		}

		log.Printf("Download request sent successfully to peer %s for hash: %s", targetPeerIDParsed, hash)
	} else if dataType == "request_all" {
		log.Printf("Sending 'request_all' signal to peer %s", targetPeerIDParsed)
		_, err = s.Write([]byte("request_all\n"))
		if err != nil {
			log.Printf("Failed to send 'request_all' signal to peer %s: %v", targetPeerIDParsed, err)
			return err
		}
		log.Printf("'Request all files' signal sent successfully to peer %s", targetPeerIDParsed)

		// Wait for half a second to let the file be received
	} else if dataType == "request_info" {
		log.Printf("Requesting file info from peer %s with hash: %s", targetPeerIDParsed, hash)

		// Send "request_info" header
		_, err = s.Write([]byte("request_info\n"))
		if err != nil {
			log.Printf("Failed to send request_info header to peer %s: %v", targetPeerIDParsed, err)
			return err
		}

		// Send the file hash
		_, err = s.Write([]byte(hash + "\n"))
		if err != nil {
			log.Printf("Failed to send file hash to peer %s: %v", targetPeerIDParsed, err)
			return err
		}

		log.Printf("File info request sent successfully to peer %s", targetPeerIDParsed)

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

func sendRequestedInfoToPeer(node host.Host, targetPeerID string, fileInfo *models.JoinedHosting) error {
	log.Printf("Preparing to send requested file info to peer %s", targetPeerID)

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

	// Write the "requested_info" header
	_, err = s.Write([]byte("requested_info\n"))
	if err != nil {
		log.Printf("Failed to send 'requested_info' header to peer %s: %v", targetPeerIDParsed, err)
		return err
	}
	log.Printf("Sent 'requested_info' header to peer %s", targetPeerID)

	// Serialize the file information into JSON
	responseData, err := json.Marshal(fileInfo)
	if err != nil {
		log.Printf("Error marshaling file information: %v", err)
		return err
	}

	// Send the JSON data
	_, err = s.Write(responseData)
	if err != nil {
		log.Printf("Failed to send file information to peer %s: %v", targetPeerIDParsed, err)
		return err
	}

	log.Printf("File information sent successfully to peer %s", targetPeerID)
	return nil
}

func handleInfoRequest(s network.Stream, db *sql.DB, node host.Host) {
	// Create a buffered reader for the stream
	reader := bufio.NewReader(s)

	// Read the hash from the stream
	hash, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading hash from peer %s: %v", s.Conn().RemotePeer(), err)
		_, _ = s.Write([]byte(fmt.Sprintf("error: %v\n", err)))
		return
	}
	hash = strings.TrimSpace(hash)
	log.Printf("Received file info request for hash: %s from peer: %s", hash, s.Conn().RemotePeer())

	// Query the database for the requested file info
	joinedHosting, err := operations.FindHosting(db, hash)
	if err != nil {
		log.Printf("Error retrieving file info for hash %s: %v", hash, err)
		_, _ = s.Write([]byte(fmt.Sprintf("error: %v\n", err)))
		return
	}

	// Send the file information back to the requesting peer
	err = sendRequestedInfoToPeer(node, s.Conn().RemotePeer().String(), joinedHosting)
	if err != nil {
		log.Printf("Failed to send requested file info for hash %s to peer %s: %v", hash, s.Conn().RemotePeer(), err)
		return
	}

	log.Printf("File info response sent successfully for hash %s to peer %s", hash, s.Conn().RemotePeer())
}
