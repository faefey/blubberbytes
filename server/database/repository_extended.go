// repository_extended.go
package database

import (
	"database/sql"
	"fmt"
)

func AddFileDataToTable(db *sql.DB, tableName string, fileData FileData) error {
	query := fmt.Sprintf(`INSERT INTO %s (hash, FileName, FileSize, sizeInGB, DateListed, type, downloads, price)
                          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, tableName)
	_, err := db.Exec(query, fileData.Hash, fileData.FileName, fileData.FileSize, fileData.SizeInGB, fileData.DateListed, fileData.Type, fileData.Downloads, fileData.Price)
	if err != nil {
		return fmt.Errorf("error adding file data to table %s: %v", tableName, err)
	}
	return nil
}

func GetFileDataFromTable(db *sql.DB, tableName string) ([]FileData, error) {
	query := fmt.Sprintf(`SELECT id, hash, FileName, FileSize, sizeInGB, DateListed, type, downloads, price FROM %s`, tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying table %s: %v", tableName, err)
	}
	defer rows.Close()

	var files []FileData
	for rows.Next() {
		var file FileData
		err := rows.Scan(&file.ID, &file.Hash, &file.FileName, &file.FileSize, &file.SizeInGB, &file.DateListed, &file.Type, &file.Downloads, &file.Price)
		if err != nil {
			return nil, fmt.Errorf("error scanning file data from table %s: %v", tableName, err)
		}
		files = append(files, file)
	}
	return files, nil
}

// DeleteHosting deletes a file from the hosting table by its hash
func DeleteHosting(db *sql.DB, fileHash string) error {
	query := `DELETE FROM hosting WHERE hash = ?`
	result, err := db.Exec(query, fileHash)
	if err != nil {
		return fmt.Errorf("error deleting file from hosting with hash %s: %v", fileHash, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows in hosting: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no file found in hosting with hash %s", fileHash)
	}
	return nil
}

// DeleteSharing deletes a file from the sharing table by its hash
func DeleteSharing(db *sql.DB, fileHash string) error {
	query := `DELETE FROM sharing WHERE hash = ?`
	result, err := db.Exec(query, fileHash)
	if err != nil {
		return fmt.Errorf("error deleting file from sharing with hash %s: %v", fileHash, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows in sharing: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no file found in sharing with hash %s", fileHash)
	}
	return nil
}

// FileExistsInTable checks if a file with the given hash exists in the specified table
func FileExistsInTable(db *sql.DB, tableName string, fileHash string) (bool, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE hash = ?`, tableName)
	var count int
	err := db.QueryRow(query, fileHash).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking file existence in table %s: %v", tableName, err)
	}
	return count > 0, nil
}
