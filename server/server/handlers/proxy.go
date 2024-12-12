package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"server/database/models"
	"server/database/operations"
	"server/p2p"

	"github.com/libp2p/go-libp2p/core/host"
)

func UpdateProxyHandler(w http.ResponseWriter, r *http.Request, node host.Host, db *sql.DB) {
	decoder := json.NewDecoder(r.Body)
	var m models.Proxy
	err := decoder.Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	walletInfo, err := operations.GetWalletInfo(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	address := walletInfo.Address

	err = p2p.ProvideKey("PROXY")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = operations.UpdateProxy(db, m.IP, m.Rate, node.ID().String(), address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RefreshProxiesHandler(w http.ResponseWriter, _ *http.Request, node host.Host, db *sql.DB) {
	proxies, err := p2p.RandomProxiesInfo(node)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proxies)
}

func ConnectToProxyHandler(w http.ResponseWriter, r *http.Request, node host.Host, db *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	connectRequest := models.ProxyBill{
		IP:     "connect",
		Rate:   0,
		Bytes:  0,
		Amount: 0,
		Wallet: "",
	}

	err = p2p.SendProxyBillWithConfirmation(node, string(body), connectRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ProxyLogsHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	proxyLogsRecords, err := operations.GetProxyLogs(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proxyLogsRecords)
}
