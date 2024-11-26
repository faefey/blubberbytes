package operations

import (
	"database/sql"
	"fmt"
	"server/database/models"
)

func AddUploads(db *sql.DB, id int64, date, hash, name, extension string, size int64) error {
	query := `INSERT INTO Uploads (id, date, hash, name, extension, size) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, id, date, hash, name, extension, size)
	if err != nil {
		return fmt.Errorf("error adding record to Upload History: %v", err)
	}

	fmt.Printf("Record added to Upload History with id: %d\n", id)
	return nil
}

func GetAllUploads(db *sql.DB) ([]models.Uploads, error) {
	query := `SELECT id, date, hash, name, extension, size FROM Uploads`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying Upload History table: %v", err)
	}
	defer rows.Close()

	var uploadsRecords []models.Uploads
	for rows.Next() {
		var record models.Uploads
		err := rows.Scan(&record.Id, &record.Date, &record.Hash, &record.Name, &record.Extension, &record.Size)
		if err != nil {
			return nil, fmt.Errorf("error scanning Upload History record: %v", err)
		}
		uploadsRecords = append(uploadsRecords, record)
	}

	return uploadsRecords, nil
}
