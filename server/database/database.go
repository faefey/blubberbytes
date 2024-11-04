// database.go
package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // go get modernc.org/sqlite
)

// SetupDatabase initializes the SQLite database and creates the FileMetadata table.
func SetupDatabase(dbPath string) (*sql.DB, error) {
	// Open a connection to the SQLite database.
	db, err := sql.Open("sqlite", dbPath) // Use "sqlite" instead of "sqlite3"
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Create the FileMetadata table.
	fileMetadataTable := `
    CREATE TABLE IF NOT EXISTS FileMetadata (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        file_size INTEGER NOT NULL,
        extension TEXT NOT NULL,
        file_name TEXT NOT NULL,
        file_price REAL NOT NULL,
        file_hash TEXT UNIQUE NOT NULL
    );`

	// Execute the table creation.
	_, err = db.Exec(fileMetadataTable)
	if err != nil {
		return nil, fmt.Errorf("error creating FileMetadata table: %v", err)
	}

	fmt.Println("FileMetadata table created successfully.")
	return db, nil
}
