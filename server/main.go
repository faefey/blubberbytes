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
	defer db.Close()

	// Creates the tables in the database
	err = database.CreateNewTables(db)
	if err != nil {
		log.Println("Error creating new tables:", err)
		return
	}

	// // Populates the database
	err = database.PopulateDatabase(db)
	if err != nil {
		log.Println("Error populating database:", err)
		return
	}

	net := "testnet"

	// err = os.Remove("./btc/walletaddress.txt")
	// if err != nil && !os.IsNotExist(err) {
	// 	log.Println("Error deleting existing walletaddress file:", err)
	// 	return
	// }

	// walletDir := btcutil.AppDataDir("btcwallet", false)
	// err = os.Remove(filepath.Join(walletDir, net+"/wallet.db"))
	// if err != nil && !os.IsNotExist(err) {
	// 	log.Println("Error deleting existing wallet file:", err)
	// 	return
	// }

	// Starts btc-related processes and saves wallet address
	btcdCmd, btcwalletCmd, btcd, btcwallet, miningaddr, err := btc.Start(net, false)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		btc.ShutdownClient(btcd)
		btc.ShutdownClient(btcwallet)

		btc.InterruptCmd(btcwalletCmd)
		btc.InterruptCmd(btcdCmd)
	}()

	node, dht, err := p2p.P2PSync()
	if err != nil {
		log.Println(err)
		return
	}

	go p2p.P2PAsync(node, dht, db)
	go gateway.Gateway(node, db)
	go server.Server(node, btcwallet, miningaddr, db)

	// Blocks until a signal is received
	<-sigs
}
