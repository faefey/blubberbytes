package database

import (
	"database/sql"
	"fmt"
)

// SetupFilesTables initializes tables related to file management (storing, hosting, sharing, saved).
func SetupFilesTables(db *sql.DB) error {
	tables := map[string]string{
		"Storing": `
			CREATE TABLE IF NOT EXISTS Storing (
				hash TEXT PRIMARY KEY NOT NULL,
				name TEXT NOT NULL,
				extension TEXT NOT NULL,
				size INTEGER NOT NULL,
				path TEXT NOT NULL,
				date TEXT NOT NULL
			);`,
		"Hosting": `
			CREATE TABLE IF NOT EXISTS Hosting (
				hash TEXT PRIMARY KEY NOT NULL,
				price REAL NOT NULL,
				FOREIGN KEY(hash) REFERENCES Storing(hash)
			);`,
		"Sharing": `
			CREATE TABLE IF NOT EXISTS Sharing (
				hash TEXT PRIMARY KEY NOT NULL,
				password TEXT NOT NULL,
				FOREIGN KEY(hash) REFERENCES Storing(hash)
			);`,
		"Saved": `
			CREATE TABLE IF NOT EXISTS Saved (
				hash TEXT PRIMARY KEY NOT NULL,
				name TEXT NOT NULL,
				extension TEXT NOT NULL,
				size INTEGER NOT NULL
			);`,
	}

	// Execute each table creation statement
	for tableName, createStmt := range tables {
		_, err := db.Exec(createStmt)
		if err != nil {
			return fmt.Errorf("error creating %s table: %v", tableName, err)
		}
		fmt.Printf("%s table created successfully.\n", tableName)
	}

	return nil
}

// SetupHistoriesTables initializes tables related to histories (uploads, downloads, transactions, proxies).
func SetupHistoriesTables(db *sql.DB) error {
	tables := map[string]string{
		"Uploads": `
			CREATE TABLE IF NOT EXISTS Uploads (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date TEXT NOT NULL,
				hash TEXT NOT NULL,
				name TEXT NOT NULL,
				extension TEXT NOT NULL,
				size INTEGER NOT NULL
			);`,
		"Downloads": `
			CREATE TABLE IF NOT EXISTS Downloads (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date TEXT NOT NULL,
				hash TEXT NOT NULL,
				name TEXT NOT NULL,
				extension TEXT NOT NULL,
				size INTEGER NOT NULL,
				price REAL NOT NULL
			);`,
		"Transactions": `
			CREATE TABLE IF NOT EXISTS Transactions (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date TEXT NOT NULL,
				wallet TEXT NOT NULL,
				amount REAL NOT NULL,
				balance REAL NOT NULL
			);`,
		"Proxies": `
			CREATE TABLE IF NOT EXISTS Proxies (
				-- Placeholder table for proxies
				id INTEGER PRIMARY KEY AUTOINCREMENT
			);`,
	}

	// Execute each table creation statement
	for tableName, createStmt := range tables {
		_, err := db.Exec(createStmt)
		if err != nil {
			return fmt.Errorf("error creating %s table: %v", tableName, err)
		}
		fmt.Printf("%s table created successfully.\n", tableName)
	}

	return nil
}

// SetupWalletInfoTable initializes the WalletInfo table with a placeholder row.
func SetupWalletInfoTable(db *sql.DB) error {
	createTable :=
		`CREATE TABLE IF NOT EXISTS WalletInfo (
			address TEXT PRIMARY KEY NOT NULL,
			pubPassphrase TEXT NOT NULL,
			privPassphrase TEXT NOT NULL
		);`

	// Execute the table creation statement
	_, err := db.Exec(createTable)
	if err != nil {
		return fmt.Errorf("error creating WalletInfo table: %v", err)
	}
	fmt.Printf("WalletInfo table created successfully.\n")

	query := `INSERT INTO WalletInfo (address, pubPassphrase, privPassphrase) VALUES (?, ?, ?)`
	_, err = db.Exec(query, "", "", "")
	if err != nil {
		return fmt.Errorf("error initializing WalletInfo table: %v", err)
	}
	fmt.Printf("WalletInfo table initialized successfully.\n")

	return nil
}

// SetupProxyTable initializes the Proxy table with a placeholder row.
func SetupProxyTable(db *sql.DB) error {
	createTable :=
		`CREATE TABLE IF NOT EXISTS Proxy (
			ip TEXT PRIMARY KEY NOT NULL,
			rate REAL NOT NULL,
			node TEXT NOT NULL,
			wallet TEXT NOT NULL
		);`

	// Execute the table creation statement
	_, err := db.Exec(createTable)
	if err != nil {
		return fmt.Errorf("error creating Proxy table: %v", err)
	}
	fmt.Printf("Proxy table created successfully.\n")

	query := `INSERT INTO Proxy (ip, rate, node, wallet) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, "", 0, "", "")
	if err != nil {
		return fmt.Errorf("error initializing Proxy table: %v", err)
	}
	fmt.Printf("Proxy table initialized successfully.\n")

	return nil
}

func SetupProxyLogsTable(db *sql.DB) error {
	createTable :=
		`CREATE TABLE IF NOT EXISTS ProxyLogs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ip TEXT NOT NULL,
			bytes INTEGER NOT NULL,
			time INTEGER NOT NULL
		);`

	// Execute the table creation statement
	_, err := db.Exec(createTable)
	if err != nil {
		return fmt.Errorf("error creating ProxyLogs table: %v", err)
	}
	fmt.Printf("ProxyLogs table created successfully.\n")

	return nil
}

func SetupIPtoNodeTable(db *sql.DB) error {
	createTable :=
		`CREATE TABLE IF NOT EXISTS IPtoNode (
			ip TEXT PRIMARY KEY NOT NULL,
			node TEXT NOT NULL
		);`

	// Execute the table creation statement
	_, err := db.Exec(createTable)
	if err != nil {
		return fmt.Errorf("error creating IPtoNode table: %v", err)
	}
	fmt.Printf("IPtoNode table created successfully.\n")

	return nil
}
