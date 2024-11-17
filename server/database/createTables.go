package database

import (
	"database/sql"
	"fmt"
)

func CreateNewTables(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS hosting (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            hash TEXT NOT NULL UNIQUE,
            FileName TEXT NOT NULL,
            FileSize TEXT NOT NULL,
            sizeInGB REAL NOT NULL,
            DateListed TEXT NOT NULL,
            type TEXT NOT NULL,
            downloads INTEGER NOT NULL,
            price REAL NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS sharing (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            hash TEXT NOT NULL UNIQUE,
            FileName TEXT NOT NULL,
            FileSize TEXT NOT NULL,
            sizeInGB REAL NOT NULL,
            DateListed TEXT NOT NULL,
            type TEXT NOT NULL,
            downloads INTEGER NOT NULL,
            price REAL NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS purchased (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            hash TEXT NOT NULL UNIQUE,
            FileName TEXT NOT NULL,
            FileSize TEXT NOT NULL,
            sizeInGB REAL NOT NULL,
            DateListed TEXT NOT NULL,
            type TEXT NOT NULL,
            downloads INTEGER NOT NULL,
            price REAL NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS explore (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            hash TEXT NOT NULL UNIQUE,
            FileName TEXT NOT NULL,
            FileSize TEXT NOT NULL,
            sizeInGB REAL NOT NULL,
            DateListed TEXT NOT NULL,
            type TEXT NOT NULL,
            downloads INTEGER NOT NULL,
            price REAL NOT NULL
        );`,
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
	}
	return nil
}

// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"
// 	"server/database"
// )

// func main() {

// 	// btcdCmd, btcwalletCmd, btcd, btcwallet, err := btc.Start("simnet", false)
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// 	return
// 	// }

// 	// btc.GetBlockCount(btcd)
// 	// btc.GetBalance(btcwallet)

// 	// _, err = btcwallet.Generate(500)
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// 	return
// 	// }
// 	// time.Sleep(time.Second * 10)

// 	// btc.GetBlockCount(btcd)
// 	// btc.GetBalance(btcwallet)

// 	// btc.ShutdownClient(btcd)
// 	// btc.ShutdownClient(btcwallet)

// 	// btc.InterruptCmd(btcwalletCmd)
// 	// btc.InterruptCmd(btcdCmd)

// 	// for restarting the database comment it if you dont want to but its going to create clash because of unique hash
// 	err := os.Remove("./data.db")
// 	if err != nil && !os.IsNotExist(err) {
// 		log.Println("Error deleting existing database file:", err)
// 		return
// 	}

// 	db, err := database.SetupDatabase("./data.db")
// 	if err != nil {
// 		log.Println("Error setting up database:", err)
// 		return
// 	}
// 	defer db.Close()

// 	// Createing the new tables
// 	err = CreateNewTables(db)
// 	if err != nil {
// 		log.Println("Error creating new tables:", err)
// 		return
// 	}

// 	// Populating the database
// 	err = database.PopulateDatabase(db)
// 	if err != nil {
// 		log.Println("Error populating database:", err)
// 		return
// 	}
// 	time.Sleep(time.Second * 3)

// 	btc.GetBlockCount(btcd)
// 	btc.GetBalance(btcwallet)

// 	btc.ShutdownClient(btcd)
// 	btc.ShutdownClient(btcwallet)

// 	btc.InterruptCmd(btcwalletCmd)
// 	btc.InterruptCmd(btcdCmd)

// 	// db, err := database.SetupDatabase("./data.db")
// 	// if err != nil {
// 	// 	log.Println("Error setting up database:", err)
// 	// 	return
// 	// }
// 	// defer db.Close()

// 	// p2p.P2P(db)
// 	// gateway()
// 	// server()
// }

// func CreateNewTables(db *sql.DB) error {
// 	tables := []string{
// 		`CREATE TABLE IF NOT EXISTS hosting (
//             id INTEGER PRIMARY KEY AUTOINCREMENT,
//             hash TEXT NOT NULL UNIQUE,
//             FileName TEXT NOT NULL,
//             FileSize TEXT NOT NULL,
//             sizeInGB REAL NOT NULL,
//             DateListed TEXT NOT NULL,
//             type TEXT NOT NULL,
//             downloads INTEGER NOT NULL,
//             price REAL NOT NULL
//         );`,
// 		`CREATE TABLE IF NOT EXISTS sharing (
//             id INTEGER PRIMARY KEY AUTOINCREMENT,
//             hash TEXT NOT NULL UNIQUE,
//             FileName TEXT NOT NULL,
//             FileSize TEXT NOT NULL,
//             sizeInGB REAL NOT NULL,
//             DateListed TEXT NOT NULL,
//             type TEXT NOT NULL,
//             downloads INTEGER NOT NULL,
//             price REAL NOT NULL
//         );`,
// 		`CREATE TABLE IF NOT EXISTS purchased (
//             id INTEGER PRIMARY KEY AUTOINCREMENT,
//             hash TEXT NOT NULL UNIQUE,
//             FileName TEXT NOT NULL,
//             FileSize TEXT NOT NULL,
//             sizeInGB REAL NOT NULL,
//             DateListed TEXT NOT NULL,
//             type TEXT NOT NULL,
//             downloads INTEGER NOT NULL,
//             price REAL NOT NULL
//         );`,
// 		`CREATE TABLE IF NOT EXISTS explore (
//             id INTEGER PRIMARY KEY AUTOINCREMENT,
//             hash TEXT NOT NULL UNIQUE,
//             FileName TEXT NOT NULL,
//             FileSize TEXT NOT NULL,
//             sizeInGB REAL NOT NULL,
//             DateListed TEXT NOT NULL,
//             type TEXT NOT NULL,
//             downloads INTEGER NOT NULL,
//             price REAL NOT NULL
//         );`,
// 	}

// 	for _, table := range tables {
// 		_, err := db.Exec(table)
// 		if err != nil {
// 			return fmt.Errorf("error creating table: %v", err)
// 		}
// 	}
// 	return nil
// }
