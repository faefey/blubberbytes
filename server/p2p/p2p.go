package p2p

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
)

var (
	relay_node_addr     = "/ip4/130.245.173.221/tcp/4001/p2p/12D3KooWDpJ7As7BWAwRMfu1VU2WCqNjvq387JEYKDBj4kx6nXTN"
	bootstrap_node_addr = "/ip4/127.0.1.0/tcp/61000/p2p/12D3KooWKLP2W9BDZhSjNkUyChEF6jVhoVbkfztu7o5mbHvQ4XcM"
)

func P2PSync() (host.Host, *dht.IpfsDHT) {
	fmt.Print("Enter your student ID: ")
	_, err := fmt.Scanln(&node_id)
	if err != nil {
		log.Fatalf("Failed to read student ID: %s", err)
	}

	node, dht, err := createNode()
	if err != nil {
		log.Fatalf("Failed to create node: %s", err)
	}

	return node, dht
}

func P2PAsync(node host.Host, dht *dht.IpfsDHT, db *sql.DB) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	globalCtx = ctx

	fmt.Println("Node Peer ID:", node.ID())

	connectToPeer(node, relay_node_addr) // connect to relay node
	makeReservation(node)                // make reservation on relay node
	go refreshReservation(node, 10*time.Minute)
	connectToPeer(node, bootstrap_node_addr) // connect to bootstrap node
	go handlePeerExchange(node)
	go receiveDataFromPeer(node, db, "D:/blubberbytes/cse416-dht-go-main/") // Ensures a folder path is used
	go handleInput(ctx, dht, node, db)                                      // Pass db connection to handleInput

	defer node.Close()

	select {}
}
