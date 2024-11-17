// db_tables.go
package database

import (
	"database/sql"
	"fmt"
)

// Make the function public by capitalizing the first letter
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
