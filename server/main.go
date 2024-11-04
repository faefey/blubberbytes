package main

import (
	"log"
	"server/btc"
	"server/database"
	"time"
)

func main() {
	btcdCmd, btcwalletCmd, btcd, btcwallet, err := btc.Start("simnet")
	if err != nil {
		log.Println(err)
		return
	}

	btc.GetBlockCount(btcd)
	btc.GetBalance(btcwallet)

	_, err = btcwallet.Generate(500)
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Second * 3)

	btc.GetBlockCount(btcd)
	btc.GetBalance(btcwallet)

	btc.ShutdownClients(btcd, btcwallet)

	btc.InterruptProcesses(btcdCmd, btcwalletCmd)

	//Database testing
	db, err := database.SetupDatabase("./data.db")
	if err != nil {
		log.Println("Error setting up database:", err)
		return
	}
	defer db.Close()

	// p2p.P2P(db)
	// gateway()
	// server()
}
