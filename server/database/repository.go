package database

import (
	"database/sql"
	"fmt"
)

// AddFileMetadata inserts a new record into the FileMetadata table
func AddFileMetadata(db *sql.DB, fileSize int64, extension, fileName string, filePrice float64, fileHash string) (int64, error) {
	query := `INSERT INTO FileMetadata (file_size, extension, file_name, file_price, file_hash) 
	          VALUES (?, ?, ?, ?, ?)`
	result, err := db.Exec(query, fileSize, extension, fileName, filePrice, fileHash)
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

	// Query to find file metadata by file hash
	query := `SELECT id, file_size, extension, file_name, file_price, file_hash FROM FileMetadata WHERE file_hash = ?`
	err := db.QueryRow(query, fileHash).Scan(
		&file.ID,
		&file.FileSize,
		&file.Extension,
		&file.FileName,
		&file.FilePrice,
		&file.FileHash,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No result found
		}
		return nil, fmt.Errorf("error querying file metadata by hash: %v", err)
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
