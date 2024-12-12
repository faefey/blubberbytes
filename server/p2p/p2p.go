package p2p

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
)

var (
	relay_node_addr     = "/ip4/130.245.173.221/tcp/4001/p2p/12D3KooWDpJ7As7BWAwRMfu1VU2WCqNjvq387JEYKDBj4kx6nXTN"
	bootstrap_node_addr = "/ip4/130.245.173.222/tcp/61020/p2p/12D3KooWM8uovScE5NPihSCKhXe8sbgdJAi88i2aXT2MmwjGWoSX"
)

func P2PSync() (host.Host, *dht.IpfsDHT, error) {
	fmt.Print("Enter your student ID: ")
	_, err := fmt.Scanln(&node_id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read student ID: %s", err)
	}

	node, dht, err := createNode()
	dhtRouting = dht
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create node: %s", err)
	}

	return node, dht, nil
}

func P2PAsync(node host.Host, dht *dht.IpfsDHT, db *sql.DB, btcwallet *rpcclient.Client, netParams *chaincfg.Params) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	globalCtx = ctx

	fmt.Println("Node Peer ID:", node.ID())

	connectToPeer(node, relay_node_addr) // connect to relay node
	makeReservation(node)                // make reservation on relay node
	go refreshReservation(node, 10*time.Minute)
	connectToPeer(node, bootstrap_node_addr) // connect to bootstrap node
	// go handlePeerExchange(node)
	go receiveDataFromPeer(node, db, "D:/blubberbytes/cse416-dht-go-main/", btcwallet, netParams) // Ensures a folder path is used
	go handleInput(ctx, dht, node, db)                                                            // Pass db connection to handleInput

	// Call the helper function to periodically provide keys
	go periodicTaskHelper(12*time.Hour, db)

	// Keep the program running
	<-ctx.Done()

	defer node.Close()
	fmt.Println("Node closed.")
}
