package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"server/database/models"
	"server/database/operations"

	"github.com/btcsuite/btcd/rpcclient"
)

func WalletHandler(w http.ResponseWriter, _ *http.Request, btcwallet *rpcclient.Client, db *sql.DB) {
	walletInfo, err := operations.GetWalletInfo(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	address := walletInfo.Address

	currentBalance, err := btcwallet.GetBalance("*")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pendingBalance, err := btcwallet.GetBalanceMinConf("*", 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wallet := models.Wallet{
		Address:        address,
		CurrentBalance: currentBalance.ToBTC(),
		PendingBalance: pendingBalance.ToBTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallet)
}

func GenerateHandler(w http.ResponseWriter, _ *http.Request, btcwallet *rpcclient.Client, db *sql.DB) {
	block, err := btcwallet.Generate(1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(block)
}
