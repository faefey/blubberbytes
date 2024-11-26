package operations

import (
	"database/sql"
	"fmt"
	"server/database/models"
)

func AddDownloads(db *sql.DB, id int64, date, hash, name, extension string, size int64, price float64) error {
	query := `INSERT INTO Downloads (id, date, hash, name, extension, size, price) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, id, date, hash, name, extension, size, price)
	if err != nil {
		return fmt.Errorf("error adding record to Download History: %v", err)
	}

	fmt.Printf("Record added to Downloads with hash: %s\n", hash)
	return nil
}

func GetAllDownloads(db *sql.DB) ([]models.Downloads, error) {
	query := `SELECT id, date, hash, name, extension, size, price FROM Downloads`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying Download History table: %v", err)
	}
	defer rows.Close()

	var downloadsRecords []models.Downloads
	for rows.Next() {
		var record models.Downloads
		err := rows.Scan(&record.Id, &record.Date, &record.Hash, &record.Name, &record.Extension, &record.Size, &record.Price)
		if err != nil {
			return nil, fmt.Errorf("error scanning Download History record: %v", err)
		}
		downloadsRecords = append(downloadsRecords, record)
	}

	return downloadsRecords, nil
}
