package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"server/database/models"
	"server/database/operations"
	"slices"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
)

func UploadsHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	uploadsRecords, err := operations.GetAllUploads(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploadsRecords)
}

func DownloadsHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	downloadsRecords, err := operations.GetAllDownloads(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(downloadsRecords)
}

func TransactionsHandler(w http.ResponseWriter, _ *http.Request, btcwallet *rpcclient.Client, db *sql.DB) {
	walletInfo, err := operations.GetWalletInfo(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	address := walletInfo.Address

	listSinceBlockResult, err := btcwallet.ListSinceBlock(nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transactions := listSinceBlockResult.Transactions

	var temp []btcjson.ListTransactionsResult
	for _, transaction := range transactions {
		if !(transaction.Category == "send" && transaction.Address == address) {
			temp = append(temp, transaction)
		}
	}
	transactions = temp

	slices.SortStableFunc(transactions, func(a, b btcjson.ListTransactionsResult) int {
		return int(b.Time - a.Time)
	})

	var transactionsRecords []models.Transactions
	for _, transaction := range transactions {
		transactionsRecords = append(transactionsRecords, models.Transactions{
			Id:            transaction.TxID,
			Date:          time.Unix(transaction.Time, 0).Local().Format("01/02/2006"),
			Wallet:        transaction.Address,
			Amount:        transaction.Amount,
			Category:      transaction.Category,
			Confirmations: transaction.Confirmations,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactionsRecords)
}
