package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // go get modernc.org/sqlite
)

// SetupDatabase initializes the SQLite database with configuration but no tables.
func SetupDatabase(dbPath string) (*sql.DB, error) {
	// Open a connection to the SQLite database.
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Enable WAL mode
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %v", err)
	}

	// Set a busy timeout
	_, err = db.Exec("PRAGMA busy_timeout = 5000;")
	if err != nil {
		return nil, fmt.Errorf("failed to set busy timeout: %v", err)
	}

	fmt.Println("Database setup complete (without tables).")
	return db, nil
}

// CreateNewTables creates all the required tables by calling setup functions.
func CreateNewTables(db *sql.DB) error {
	// Create histories tables
	err := SetupHistoriesTables(db)
	if err != nil {
		return fmt.Errorf("failed to set up histories tables: %v", err)
	}

	// Create files tables
	err = SetupFilesTables(db)
	if err != nil {
		return fmt.Errorf("failed to set up files tables: %v", err)
	}

	// Create WalletInfo table
	err = SetupWalletInfoTable(db)
	if err != nil {
		return fmt.Errorf("failed to set up WalletInfo table: %v", err)
	}

	// Create Proxy table
	err = SetupProxyTable(db)
	if err != nil {
		return fmt.Errorf("failed to set up Proxy table: %v", err)
	}

	// Create ProxyLogs table
	err = SetupProxyLogsTable(db)
	if err != nil {
		return fmt.Errorf("failed to set up ProxyLogs table: %v", err)
	}

	// Create IPtoNode table
	err = SetupIPtoNodeTable(db)
	if err != nil {
		return fmt.Errorf("failed to set up IPtoNode table: %v", err)
	}

	fmt.Println("All new tables created successfully.")
	return nil
}
