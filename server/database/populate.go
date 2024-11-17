// populate.go
package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func PopulateDatabase(db *sql.DB) error {

	jsonFiles := map[string]string{
		"Database/tableData1.json": "hosting",
		"Database/tableData2.json": "sharing",
		"Database/tableData3.json": "purchased",
		"Database/tableData4.json": "explore",
	}

	for fileName, tableName := range jsonFiles {
		data, err := loadJSONFile(fileName)
		if err != nil {
			return fmt.Errorf("error loading data from %s: %v", fileName, err)
		}

		for _, file := range data {
			// Inserting the data into the appropriate tables
			err := AddFileDataToTable(db, tableName, file)
			if err != nil {
				fmt.Printf("Error adding the file data to table %s: %v\n", tableName, err)
				return err
			}
		}
	}

	fmt.Println("Database populated successfully :) ")
	return nil
}

func loadJSONFile(filePath string) ([]FileData, error) {
	var data []FileData

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON data from %s: %v", filePath, err) // decodes json
	}

	return data, nil
}
