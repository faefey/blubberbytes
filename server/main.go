package main

import (
	"log"
	"os"
	"os/signal"
	"server/btc"
	"server/database"
	"server/gateway"
	"server/p2p"
	"server/server"
	"syscall"
)

func main() {
	// Creates a channel to receive signals
	sigs := make(chan os.Signal, 1)

	// Notifies the channel on signals
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Resets the database
	err := os.Remove("./database/data.db")
	if err != nil && !os.IsNotExist(err) {
		log.Println("Error deleting existing database file:", err)
		return
	}

	// Initializes the database
	db, err := database.SetupDatabase("./database/data.db")
	if err != nil {
		log.Println("Error setting up database:", err)
		return
	}

	// Creates the tables in the database
	err = database.CreateNewTables(db)
	if err != nil {
		log.Println("Error creating new tables:", err)
		return
	}

	// Populates the database
	err = database.PopulateDatabase(db)
	if err != nil {
		log.Println("Error populating database:", err)
		return
	}

	// Starts btc-related processes and saves wallet address
	btcdCmd, btcwalletCmd, btcd, btcwallet, err := btc.Start("simnet", false)
	if err != nil {
		log.Println(err)
		return
	}

	node, dht := p2p.P2PSync()
	go p2p.P2PAsync(node, dht, db)
	go gateway.Gateway(node, db)
	go server.Server(node, btcwallet, db)

	// Blocks until a signal is received
	sig := <-sigs
	log.Println("Received signal:", sig)

	// Performs cleanup or graceful shutdown here
	log.Println("Performing cleanup...")

	btc.ShutdownClient(btcd)
	btc.ShutdownClient(btcwallet)

	btc.InterruptCmd(btcwalletCmd)
	btc.InterruptCmd(btcdCmd)

	db.Close()
}
