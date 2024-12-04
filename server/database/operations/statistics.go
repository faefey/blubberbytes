package operations

import (
	"database/sql"
	"fmt"
	"server/database/models"
)

// CalcStatistics calculates statistics from the Storing, Hosting, Sharing, and Saved tables.
func CalcStatistics(db *sql.DB) (models.Statistics, error) {
	tables := []string{"Storing", "Hosting", "Sharing", "Saved"}
	var stats [8]int64
	for i, table := range tables {
		var num int64
		err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&num)
		if err != nil {
			return models.Statistics{}, fmt.Errorf("error counting rows in %s table: %v", table, err)
		}

		var size int64
		if table == "Storing" || table == "Saved" {
			err = db.QueryRow("SELECT SUM(size) FROM " + table).Scan(&size)
		} else {
			err = db.QueryRow(fmt.Sprintf("SELECT SUM(size) FROM %s JOIN Storing ON %s.hash == Storing.hash", table, table)).Scan(&size)
		}
		if err != nil {
			return models.Statistics{}, fmt.Errorf("error calculating total size in %s table: %v", table, err)
		}

		stats[i*2] = num
		stats[i*2+1] = size
	}

	statistics := models.Statistics{
		StoringNum:  stats[0],
		StoringSize: stats[1],
		HostingNum:  stats[2],
		HostingSize: stats[3],
		SharingNum:  stats[4],
		SharingSize: stats[5],
		SavedNum:    stats[6],
		SavedSize:   stats[7],
	}

	return statistics, nil
}
