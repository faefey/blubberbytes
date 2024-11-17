package main

import (
	"log"
	"os"
	"server/database"
)

func main() {
	// btcdCmd, btcwalletCmd, btcd, btcwallet, err := btc.Start("simnet", false)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// btc.GetBlockCount(btcd)
	// btc.GetBalance(btcwallet)

	// _, err = btcwallet.Generate(500)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// time.Sleep(time.Second * 10)

	// btc.GetBlockCount(btcd)
	// btc.GetBalance(btcwallet)

	// btc.ShutdownClient(btcd)
	// btc.ShutdownClient(btcwallet)

	// btc.InterruptCmd(btcwalletCmd)
	// btc.InterruptCmd(btcdCmd)

	// for restarting the database comment it if you dont want to but its going to create clash because of unique hash
	err := os.Remove("./data.db")
	if err != nil && !os.IsNotExist(err) {
		log.Println("Error deleting existing database file:", err)
		return
	}

	db, err := database.SetupDatabase("./data.db")
	if err != nil {
		log.Println("Error setting up database:", err)
		return
	}
	defer db.Close()

	// Creating the new tables
	err = database.CreateNewTables(db)
	if err != nil {
		log.Println("Error creating new tables:", err)
		return
	}

	// Populating the database
	err = database.PopulateDatabase(db)
	if err != nil {
		log.Println("Error populating database:", err)
		return
	}

	// p2p.P2P(db)
	// gateway()
	server(db)
}
