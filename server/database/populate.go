package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"server/database/models"
	"server/database/operations"
)

// PopulateDatabase populates the database with data from JSON files.
func PopulateDatabase(db *sql.DB) error {
	// Populate Storing table
	err := populateStoring(db, "database/test_data/storing.json")
	if err != nil {
		return fmt.Errorf("error populating Storing table: %v", err)
	}

	// Populate Hosting table
	err = populateHosting(db, "database/test_data/hosting.json")
	if err != nil {
		return fmt.Errorf("error populating Hosting table: %v", err)
	}

	// Populate Sharing table
	err = populateSharing(db, "database/test_data/sharing.json")
	if err != nil {
		return fmt.Errorf("error populating Sharing table: %v", err)
	}

	// Populate Saved table
	err = populateSaved(db, "database/test_data/saved.json")
	if err != nil {
		return fmt.Errorf("error populating Saved table: %v", err)
	}

	// Populate downloads table
	err = populateDownloads(db, "database/test_data/downloads.json")
	if err != nil {
		return fmt.Errorf("error populating Download History table: %v", err)
	}

	// Populate uploads table
	err = populateUploads(db, "database/test_data/uploads.json")
	if err != nil {
		return fmt.Errorf("error populating Upload History table: %v", err)
	}

	fmt.Println("Database populated successfully.")
	return nil
}

func populateStoring(db *sql.DB, filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading %s: %v", filePath, err)
	}

	var storingRecords []models.Storing
	err = json.Unmarshal(data, &storingRecords)
	if err != nil {
		return fmt.Errorf("error parsing %s: %v", filePath, err)
	}

	for _, record := range storingRecords {
		err = operations.AddStoring(db, record.Hash, record.Name, record.Extension, record.Path, record.Date, record.Size)
		if err != nil {
			return fmt.Errorf("error inserting into Storing: %v", err)
		}
	}

	return nil
}

func populateHosting(db *sql.DB, filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading %s: %v", filePath, err)
	}

	var hostingRecords []models.Hosting
	err = json.Unmarshal(data, &hostingRecords)
	if err != nil {
		return fmt.Errorf("error parsing %s: %v", filePath, err)
	}

	for _, record := range hostingRecords {
		err = operations.AddHosting(db, record.Hash, record.Price)
		if err != nil {
			return fmt.Errorf("error inserting into Hosting: %v", err)
		}
	}

	return nil
}

func populateSharing(db *sql.DB, filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading %s: %v", filePath, err)
	}

	var sharingRecords []models.Sharing
	err = json.Unmarshal(data, &sharingRecords)
	if err != nil {
		return fmt.Errorf("error parsing %s: %v", filePath, err)
	}

	for _, record := range sharingRecords {
		err = operations.AddSharing(db, record.Hash, record.Password)
		if err != nil {
			return fmt.Errorf("error inserting into Sharing: %v", err)
		}
	}

	return nil
}

func populateSaved(db *sql.DB, filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading %s: %v", filePath, err)
	}

	var savedRecords []models.Saved
	err = json.Unmarshal(data, &savedRecords)
	if err != nil {
		return fmt.Errorf("error parsing %s: %v", filePath, err)
	}

	for _, record := range savedRecords {
		err = operations.AddSaved(db, record.Hash, record.Name, record.Extension, record.Size)
		if err != nil {
			return fmt.Errorf("error inserting into Saved: %v", err)
		}
	}

	return nil
}

func populateUploads(db *sql.DB, filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading %s: %v", filePath, err)
	}

	var uploadsRecords []models.Uploads
	err = json.Unmarshal(data, &uploadsRecords)
	if err != nil {
		return fmt.Errorf("error parsing %s: %v", filePath, err)
	}

	for _, record := range uploadsRecords {
		err = operations.AddUploads(db, record.Date, record.Hash, record.Name, record.Extension, record.Size)
		if err != nil {
			return fmt.Errorf("error inserting into Upload History: %v", err)
		}
	}

	return nil
}

func populateDownloads(db *sql.DB, filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading %s: %v", filePath, err)
	}

	var downloadsRecords []models.Downloads
	err = json.Unmarshal(data, &downloadsRecords)
	if err != nil {
		return fmt.Errorf("error parsing %s: %v", filePath, err)
	}

	for _, record := range downloadsRecords {
		err = operations.AddDownloads(db, record.Date, record.Hash, record.Name, record.Extension, record.Size, record.Price)
		if err != nil {
			return fmt.Errorf("error inserting into Download History: %v", err)
		}
	}

	return nil
}
