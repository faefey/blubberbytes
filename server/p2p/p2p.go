package p2p

import (
	"context"
	"fmt"
	"log"
	"time"
)

var (
	node_id             string
	relay_node_addr     = "/ip4/130.245.173.221/tcp/4001/p2p/12D3KooWDpJ7As7BWAwRMfu1VU2WCqNjvq387JEYKDBj4kx6nXTN"
	bootstrap_node_addr = "/ip4/130.245.173.222/tcp/61000/p2p/12D3KooWQd1K1k8XA9xVEzSAu7HUCodC7LJB6uW5Kw4VwkRdstPE"
	globalCtx           context.Context
)

func P2P() {
	fmt.Print("Enter your student ID: ")
	_, err := fmt.Scanln(&node_id)
	if err != nil {
		log.Fatalf("Failed to read student ID: %s", err)
	}

	node, dht, err := createNode()
	if err != nil {
		log.Fatalf("Failed to create node: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	globalCtx = ctx

	fmt.Println("Node Peer ID:", node.ID())

	connectToPeer(node, relay_node_addr) // connect to relay node
	makeReservation(node)                // make reservation on relay node
	go refreshReservation(node, 10*time.Minute)
	connectToPeer(node, bootstrap_node_addr) // connect to bootstrap node
	go handlePeerExchange(node)
	go receiveDataFromPeer(node, "D:/blubberbytes/cse416-dht-go-main/") // Ensures a folder path is used
	go handleInput(ctx, dht, node)

	defer node.Close()

	select {}
}
