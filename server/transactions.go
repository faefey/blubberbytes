package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"server/btc"
	"syscall"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
)

func transactions(w http.ResponseWriter, r *http.Request, btcwallet *rpcclient.Client) {
	file, err := os.Open("./btc/walletaddress.txt")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	miningaddr, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}

	addr, err := btcutil.DecodeAddress(string(miningaddr), &chaincfg.SimNetParams)
	if err != nil {
		log.Println(err)
	}

	_, err = btcwallet.Generate(500)
	if err != nil {
		log.Println(err)
	}
	time.Sleep(time.Second * 3)
	bal, err := btcwallet.GetBalance("default")
	if err != nil {
		log.Println(err)
	}
	log.Println(bal)
	err = btcwallet.WalletPassphrase("hi", 10000)
	if err != nil {
		log.Println(err)
	}
	_, err = btcwallet.SendFrom("default", addr, 100)
	if err != nil {
		log.Println(err)
	}
	_, err = btcwallet.SendFrom("default", addr, 150)
	if err != nil {
		log.Println(err)
	}
	_, err = btcwallet.SendFrom("default", addr, 200)
	if err != nil {
		log.Println(err)
	}
	_, err = btcwallet.SendFrom("default", addr, 250)
	if err != nil {
		log.Println(err)
	}
	bal, err = btcwallet.GetBalance("default")
	if err != nil {
		log.Println(err)
	}
	log.Println(bal)

	// a, err := btcwallet.GetRawMempool()
	// if err != nil {
	// 	log.Println("1 failed: " + err.Error())
	// }
	// fmt.Println(1, a)

	// b, err := btcwallet.GetRawMempoolVerbose()
	// if err != nil {
	// 	log.Println("2 failed: " + err.Error())
	// }
	// fmt.Println(2, b)

	f, err := btcwallet.ListAddressTransactions([]btcutil.Address{addr}, "*")
	if err != nil {
		log.Println("7 failed: " + err.Error())
	}
	fmt.Println(7, f)

	g, err := btcwallet.ListSinceBlock(nil)
	if err != nil {
		log.Println("8 failed: " + err.Error())
	}
	fmt.Println(8, g)

	h, err := btcwallet.ListTransactions("*")
	if err != nil {
		log.Println("9 failed: " + err.Error())
	}
	fmt.Println(9, h)

	i, err := btcwallet.ListTransactionsCount("*", 10)
	if err != nil {
		log.Println("10 failed: " + err.Error())
	}
	fmt.Println(10, i)

	// j, err := btcwallet.ListTransactionsCountFrom("*", 10, 0)
	// if err != nil {
	// 	log.Println("11 failed: " + err.Error())
	// }
	// fmt.Println(11, j)

	// k, err := btcwallet.SearchRawTransactions(addr, 0, 10, false, nil)
	// if err != nil {
	// 	log.Println("12 failed: " + err.Error())
	// }
	// fmt.Println(12, k)

	l, err := btcwallet.SearchRawTransactionsVerbose(addr, 0, 10, true, false, nil)
	if err != nil {
		log.Println("13 failed: " + err.Error())
	}
	fmt.Println(13, l)

	// c, err := btcwallet.GetRawTransaction(nil)
	// if err != nil {
	// 	log.Println("3 failed: " + err.Error())
	// }
	// fmt.Println(3, c)

	// d, err := btcwallet.GetRawTransactionVerbose(nil)
	// if err != nil {
	// 	log.Println("4 failed: " + err.Error())
	// }
	// fmt.Println(4, d)

	// e, err := btcwallet.GetTransaction(nil)
	// if err != nil {
	// 	log.Println("5 failed: " + err.Error())
	// }
	// fmt.Println(5, e)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(l)
}

func main_() {
	// Creates a channel to receive signals
	sigs := make(chan os.Signal, 1)

	// Notifies the channel on signals
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	net := "simnet"

	err := os.Remove("./btc/walletaddress.txt")
	if err != nil && !os.IsNotExist(err) {
		log.Println("Error deleting existing database file:", err)
		return
	}

	walletDir := btcutil.AppDataDir("btcwallet", false)
	err = os.Remove(filepath.Join(walletDir, net+"/wallet.db"))
	if err != nil && !os.IsNotExist(err) {
		log.Println("Error deleting existing database file:", err)
		return
	}

	// Starts btc-related processes and saves wallet address
	btcdCmd, btcwalletCmd, btcd, btcwallet, err := btc.Start(net, false)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		btc.ShutdownClient(btcd)
		btc.ShutdownClient(btcwallet)

		btc.InterruptCmd(btcwalletCmd)
		btc.InterruptCmd(btcdCmd)
	}()

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		transactions(w, r, btcwallet)
	})

	// Run the server
	fmt.Println("Server is running on port 3005...")
	go http.ListenAndServe(":3005", nil)

	<-sigs
}
