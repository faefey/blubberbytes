package p2p

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	record "github.com/libp2p/go-libp2p-record"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	"github.com/multiformats/go-multiaddr"
)

var (
	dhtRouting *dht.IpfsDHT
	globalCtx  context.Context
	node_id    string
)

type CustomValidator struct{}

func (v *CustomValidator) Validate(key string, value []byte) error {
	return nil
}

func (v *CustomValidator) Select(key string, values [][]byte) (int, error) {
	return 0, nil
}

func generatePrivateKeyFromSeed(seed []byte) (crypto.PrivKey, error) {
	hash := sha256.Sum256(seed) // Generate deterministic key material
	// Create an Ed25519 private key from the hash
	privKey, _, err := crypto.GenerateEd25519Key(
		bytes.NewReader(hash[:]),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	return privKey, nil
}

func createNode() (host.Host, *dht.IpfsDHT, error) {
	ctx := context.Background()
	globalCtx = ctx

	seed := []byte(node_id)
	customAddr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/0")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse multiaddr: %w", err)
	}
	privKey, err := generatePrivateKeyFromSeed(seed)
	if err != nil {
		panic(err)
	}
	relayAddr, err := multiaddr.NewMultiaddr(relay_node_addr)
	if err != nil {
		panic(fmt.Sprintf("Failed to create relay multiaddr: %v", err))
	}

	// Convert the relay multiaddress to AddrInfo
	relayInfo, err := peer.AddrInfoFromP2pAddr(relayAddr)
	if err != nil {
		panic(fmt.Sprintf("Failed to create AddrInfo from relay multiaddr: %v", err))
	}

	node, err := libp2p.New(
		libp2p.ListenAddrs(customAddr),
		libp2p.Identity(privKey),
		libp2p.NATPortMap(),
		libp2p.EnableNATService(),
		libp2p.EnableAutoRelayWithStaticRelays([]peer.AddrInfo{*relayInfo}),
		libp2p.EnableRelayService(),
		libp2p.EnableHolePunching(),
	)

	if err != nil {
		return nil, nil, err
	}
	_, err = relay.New(node)
	if err != nil {
		log.Printf("Failed to instantiate the relay: %v", err)
	}

	dhtRouting, err := dht.New(ctx, node, dht.Mode(dht.ModeClient))
	if err != nil {
		return nil, nil, err
	}
	namespacedValidator := record.NamespacedValidator{
		"orcanet": &CustomValidator{}, // Add a custom validator for the "orcanet" namespace
	}

	dhtRouting.Validator = namespacedValidator // Configure the DHT to use the custom validator

	err = dhtRouting.Bootstrap(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Set up notifications for new connections
	node.Network().Notify(&network.NotifyBundle{
		ConnectedF: func(n network.Network, conn network.Conn) {
			peerID := conn.RemotePeer().String()

			// Show a specific message based on the peer type after a successful connection
			switch peerID {
			case "12D3KooWDpJ7As7BWAwRMfu1VU2WCqNjvq387JEYKDBj4kx6nXTN":
				fmt.Println("Connected to Relay Node")
			case "12D3KooWM8uovScE5NPihSCKhXe8sbgdJAi88i2aXT2MmwjGWoSX":
				fmt.Println("Connected to Bootstrap Node1")
			case "12D3KooWE1xpVccUXZJWZLVWPxXzUJQ7kMqN8UQ2WLn9uQVytmdA":
				fmt.Println("Connected to Bootstrap Node2")

				// Log additional details for debugging
				fmt.Printf("Bootstrap Node2 Multiaddr: %s\n", conn.RemoteMultiaddr().String())
				fmt.Printf("Local Multiaddr: %s\n", conn.LocalMultiaddr().String())

				// Check if the connection is inbound or outbound
				if conn.Stat().Direction == network.DirInbound {
					fmt.Println("Connection direction: Inbound")
				} else {
					fmt.Println("Connection direction: Outbound")
				}

				// Check the number of connected peers to see if it's repeatedly connecting
				connectedPeers := node.Network().Peers()
				fmt.Printf("Total connected peers: %d\n", len(connectedPeers))

			default:
				addPeerID(peerID)
				fmt.Println("Connected to peer:", peerID)
			}
		},
		DisconnectedF: func(n network.Network, conn network.Conn) {
			log.Printf("Disconnected from peer: %s", conn.RemotePeer().String())
		},
	})

	return node, dhtRouting, nil
}
