package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"server/server/handlers"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/libp2p/go-libp2p/core/host"
)

// handler for HTTP proxy setup:
func setupHTTPProxy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Recieved request to setup HTTP proxy")
	w.Write([]byte("HTTP proxy setup request received"))
}

// handler for viewing a random neighbor's files:
func viewRandomNeighborFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request to view random neighbors files")
	w.Write([]byte("Random neighbors files displayed"))
}

func cors(w http.ResponseWriter, r *http.Request, handler func()) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	} else {
		handler()
	}
}

func Server(node host.Host, btcwallet *rpcclient.Client, netParams *chaincfg.Params, db *sql.DB) {
	http.HandleFunc("/setupHTTPProxy", setupHTTPProxy)
	http.HandleFunc("/viewRandomNeighborFiles", viewRandomNeighborFiles)

	// GET routes
	http.HandleFunc("/storing", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.StoringHandler(w, r, db) })
	})

	http.HandleFunc("/hosting", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.HostingHandler(w, r, db) })
	})

	http.HandleFunc("/sharing", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.SharingHandler(w, r, db) })
	})

	http.HandleFunc("/saved", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.SavedHandler(w, r, db) })
	})

	http.HandleFunc("/statistics", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.StatisticsHandler(w, r, db) })
	})

	http.HandleFunc("/uploads", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.UploadsHandler(w, r, db) })
	})

	http.HandleFunc("/downloads", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.DownloadsHandler(w, r, db) })
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.TransactionsHandler(w, r, btcwallet, db) })
	})

	http.HandleFunc("/proxies", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.ProxiesHandler(w, r, db) })
	})

	http.HandleFunc("/wallet", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.WalletHandler(w, r, btcwallet, db) })
	})

	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.GenerateHandler(w, r, btcwallet, db) })
	})

	http.HandleFunc("/refreshproxies", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.RefreshProxiesHandler(w, r, node, db) })
	})

	http.HandleFunc("/proxylogs", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.ProxyLogsHandler(w, r, db) })
	})

	// POST routes
	http.HandleFunc("/getproviders", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.GetProvidersHandler(w, r, node, db) })
	})

	http.HandleFunc("/requestmetadata", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.RequestMetadataHandler(w, r, node, db) })
	})

	http.HandleFunc("/downloadfile", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.DownloadFileHandler(w, r, node, btcwallet, netParams, db) })
	})

	http.HandleFunc("/explore", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.ExploreHandler(w, r, node, db) })
	})

	http.HandleFunc("/addstoring", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.AddStoringHandler(w, r, db) })
	})

	http.HandleFunc("/deletestoring", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.DeleteStoringHandler(w, r, db) })
	})

	http.HandleFunc("/addhosting", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.AddHostingHandler(w, r, db) })
	})

	http.HandleFunc("/deletehosting", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.DeleteHostingHandler(w, r, db) })
	})

	http.HandleFunc("/addsharing", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.AddSharingHandler(w, r, node, db) })
	})

	http.HandleFunc("/deletesharing", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.DeleteSharingHandler(w, r, db) })
	})

	http.HandleFunc("/sharinglink", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.SharingLinkHandler(w, r, node, db) })
	})

	http.HandleFunc("/addsaved", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.AddSavedHandler(w, r, db) })
	})

	http.HandleFunc("/deletesaved", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.DeleteSavedHandler(w, r, db) })
	})

	http.HandleFunc("/updateproxy", func(w http.ResponseWriter, r *http.Request) {
		cors(w, r, func() { handlers.UpdateProxyHandler(w, r, node, db) })
	})

	// Run the server
	fmt.Println("Server is running on port 3001...")
	if err := http.ListenAndServe(":3001", nil); err != nil {
		panic(fmt.Sprintf("Server failed: %s", err))
	}
}
