package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"server/database/models"

	"github.com/btcsuite/btcd/rpcclient"
)

func WalletHandler(w http.ResponseWriter, _ *http.Request, btcwallet *rpcclient.Client, db *sql.DB) {
	balance, err := btcwallet.GetBalance("default")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, err := os.Open("walletaddress.txt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	miningaddr, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wallet := models.Wallet{
		Address: string(miningaddr),
		Balance: balance.ToBTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallet)
}
