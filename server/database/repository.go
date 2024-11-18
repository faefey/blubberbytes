package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// AddFileMetadata inserts a new record into the FileMetadata table
func AddFileMetadata(db *sql.DB, fileSize int64, extension, fileName string, filePrice float64, fileHash, path string) (int64, error) {
	// Initially, passwords will be an empty list.
	passwordsJSON, err := json.Marshal([]string{})
	if err != nil {
		return 0, fmt.Errorf("error marshaling empty passwords: %v", err)
	}

	query := `INSERT INTO FileMetadata (file_size, extension, file_name, file_price, file_hash, passwords, path) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, fileSize, extension, fileName, filePrice, fileHash, string(passwordsJSON), path)
	if err != nil {
		return 0, fmt.Errorf("error adding file metadata: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error fetching last insert ID: %v", err)
	}

	fmt.Printf("FileMetadata added with ID: %d\n", id)
	return id, nil
}

// DeleteFileMetadata removes a record from the FileMetadata table by its ID
func DeleteFileMetadata(db *sql.DB, id int64) error {
	query := `DELETE FROM FileMetadata WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting file metadata with ID %d: %v", id, err)
	}

	fmt.Printf("FileMetadata with ID %d deleted successfully.\n", id)
	return nil
}

// FindFileMetadataByHash retrieves a file's metadata from the FileMetadata table by its file hash
func FindFileMetadataByHash(db *sql.DB, fileHash string) (*FileMetadata, error) {
	var file FileMetadata
	var passwordsJSON string

	// Query to find file metadata by file hash
	query := `SELECT id, file_size, extension, file_name, file_price, file_hash, passwords, path FROM FileMetadata WHERE file_hash = ?`
	err := db.QueryRow(query, fileHash).Scan(
		&file.ID,
		&file.FileSize,
		&file.Extension,
		&file.FileName,
		&file.FilePrice,
		&file.FileHash,
		&passwordsJSON,
		&file.Path,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No result found
		}
		return nil, fmt.Errorf("error querying file metadata by hash: %v", err)
	}

	err = json.Unmarshal([]byte(passwordsJSON), &file.Passwords)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling passwords: %v", err)
	}

	return &file, nil
}

// DeleteFileMetadataByHash removes a record from the FileMetadata table by its file hash
func DeleteFileMetadataByHash(db *sql.DB, fileHash string) error {
	query := `DELETE FROM FileMetadata WHERE file_hash = ?`
	result, err := db.Exec(query, fileHash)
	if err != nil {
		return fmt.Errorf("error deleting file metadata with hash %s: %v", fileHash, err)
	}

	// Check if any row was affected (i.e., if the file was found and deleted)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error fetching affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no file found with the given hash: %s", fileHash)
	}

	fmt.Printf("File metadata with hash %s deleted successfully.\n", fileHash)
	return nil
}

// AddPasswordToFile adds a new password to the list of passwords for the given file hash.
func AddPasswordToFile(db *sql.DB, fileHash, newPassword string) error {
	// Step 1: Retrieve the current passwords for the file hash
	var passwordsJSON string
	query := `SELECT passwords FROM FileMetadata WHERE file_hash = ?`
	err := db.QueryRow(query, fileHash).Scan(&passwordsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("file hash not found: %s", fileHash)
		}
		return fmt.Errorf("error retrieving passwords for file hash %s: %v", fileHash, err)
	}

	// Step 2: Unmarshal the current passwords into a slice
	var passwords []string
	err = json.Unmarshal([]byte(passwordsJSON), &passwords)
	if err != nil {
		return fmt.Errorf("error unmarshaling passwords for file hash %s: %v", fileHash, err)
	}

	// Step 3: Add the new password to the list (prevent duplicates)
	for _, p := range passwords {
		if p == newPassword {
			return fmt.Errorf("password already exists for file hash %s", fileHash)
		}
	}
	passwords = append(passwords, newPassword)

	// Step 4: Marshal the updated passwords back to JSON
	updatedPasswordsJSON, err := json.Marshal(passwords)
	if err != nil {
		return fmt.Errorf("error marshaling updated passwords for file hash %s: %v", fileHash, err)
	}

	// Step 5: Update the passwords in the database
	updateQuery := `UPDATE FileMetadata SET passwords = ? WHERE file_hash = ?`
	_, err = db.Exec(updateQuery, updatedPasswordsJSON, fileHash)
	if err != nil {
		return fmt.Errorf("error updating passwords for file hash %s: %v", fileHash, err)
	}

	fmt.Printf("Password added successfully for file hash: %s\n", fileHash)
	return nil
}
