package p2p

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
	"github.com/multiformats/go-multiaddr"
)

var (
	peerIDList []string
	listMutex  sync.Mutex // Mutex to ensure thread-safe access
)

func makeReservation(node host.Host) {
	ctx := globalCtx
	relayInfo, err := peer.AddrInfoFromString(relay_node_addr)
	if err != nil {
		panic(fmt.Sprintf("Failed to create addrInfo from string representation of relay multiaddr: %v", err))
	}
	_, err = client.Reserve(ctx, node, *relayInfo)
	if err != nil {
		panic(fmt.Sprintf("Failed to make reservation on relay: %v", err))
	}
	fmt.Printf("Reservation successfull \n")
}

func refreshReservation(node host.Host, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			makeReservation(node)
		case <-globalCtx.Done():
			fmt.Println("Context done, stopping reservation refresh.")
			return
		}
	}
}

func connectToPeer(node host.Host, peerAddr string) {
	// Log every time the function is called
	log.Printf("Attempting to connect to peer with address: %s", peerAddr)

	addr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		log.Printf("Failed to parse peer address: %s", err)
		return
	}

	info, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		log.Printf("Failed to get AddrInfo from address: %s", err)
		return
	}

	node.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	err = node.Connect(context.Background(), *info)
	if err != nil {
		log.Printf("Failed to connect to peer: %s", err)
		return
	}

}

func connectToPeerUsingRelay(node host.Host, targetPeerID string) {
	ctx := globalCtx
	targetPeerID = strings.TrimSpace(targetPeerID)
	relayAddr, err := multiaddr.NewMultiaddr(relay_node_addr)
	if err != nil {
		log.Printf("Failed to create relay multiaddr: %v", err)
	}
	peerMultiaddr := relayAddr.Encapsulate(multiaddr.StringCast("/p2p-circuit/p2p/" + targetPeerID))

	relayedAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMultiaddr)
	if err != nil {
		log.Println("Failed to get relayed AddrInfo: %w", err)
		return
	}
	// Connect to the peer through the relay
	err = node.Connect(ctx, *relayedAddrInfo)
	if err != nil {
		log.Println("Failed to connect to peer through relay: %w", err)
		return
	}
	addPeerID(targetPeerID)

}

func handlePeerExchange(node host.Host) {
	relayInfo, _ := peer.AddrInfoFromString(relay_node_addr)
	node.SetStreamHandler("/orcanet/p2p", func(s network.Stream) {
		defer s.Close()

		buf := bufio.NewReader(s)
		peerAddr, err := buf.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("error reading from stream: %v", err)
			}
		}
		peerAddr = strings.TrimSpace(peerAddr)
		var data map[string]interface{}
		err = json.Unmarshal([]byte(peerAddr), &data)

		if err != nil {
			fmt.Printf("error unmarshaling JSON: %v", err)
		}
		if knownPeers, ok := data["known_peers"].([]interface{}); ok {
			for _, peer := range knownPeers {
				fmt.Println("Peer:")
				if peerMap, ok := peer.(map[string]interface{}); ok {
					if peerID, ok := peerMap["peer_id"].(string); ok {
						if string(peerID) != string(relayInfo.ID) {
							connectToPeerUsingRelay(node, peerID)
						}
					}
				}
			}
		}
	})
}

func addPeerID(peerID string) {
	listMutex.Lock()
	defer listMutex.Unlock()

	// Check if the peerID already exists in the list
	for _, id := range peerIDList {
		if id == peerID {
			return // Peer ID already exists, do nothing
		}
	}

	// Add the peerID to the global list
	peerIDList = append(peerIDList, peerID)
}

func printPeerList() {
	listMutex.Lock()
	defer listMutex.Unlock()

	fmt.Println("Current Peer IDs:")
	for _, id := range peerIDList {
		fmt.Println(id)
	}
}
