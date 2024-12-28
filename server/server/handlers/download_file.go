package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/database/operations"
	"server/p2p"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/libp2p/go-libp2p/core/host"
)

func GetProvidersHandler(w http.ResponseWriter, r *http.Request, node host.Host, db *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	providers, err := p2p.GetProviderIDs(node, string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providers)
}

func RequestMetadataHandler(w http.ResponseWriter, r *http.Request, node host.Host, db *sql.DB) {
	decoder := json.NewDecoder(r.Body)
	var request struct {
		Peer string `json:"peer"`
		Hash string `json:"hash"`
	}
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metadata, err := p2p.RequestFileInfo(node, request.Peer, request.Hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metadata)
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request, node host.Host, btcwallet *rpcclient.Client, netParams *chaincfg.Params, db *sql.DB) {
	decoder := json.NewDecoder(r.Body)
	var request struct {
		Peer  string  `json:"peer"`
		Hash  string  `json:"hash"`
		Price float64 `json:"price"`
	}
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	walletInfo, err := operations.GetWalletInfo(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = btcwallet.WalletPassphrase(walletInfo.PrivPassphrase, 300)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name, data, ext, address, err := p2p.SimplyDownload(node, request.Peer, request.Hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	btcutilAddress, err := btcutil.DecodeAddress(address, netParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = btcwallet.SendFrom("default", btcutilAddress, btcutil.Amount(request.Price*1e8))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	date := time.Now().Local().Format("01/02/2006")
	err = operations.AddDownloads(db, date, request.Hash, name, ext, int64(len(data)), request.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var contentType string
	if ext == "" {
		contentType = "application/octet-stream"
	} else {
		contentType = ext
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", name))
	w.Write(data)
}
