package main

import (
	"log"
	"server/btc"
	"time"
)

func main() {
	btcdCmd, btcwalletCmd, btcd, btcwallet, err := btc.Start("simnet")
	if err != nil {
		log.Println(err)
		return
	}

	btc.GetBlockCount(btcd)
	btc.GetBalance(btcwallet)

	_, err = btcwallet.Generate(500)
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Second * 3)

	btc.GetBlockCount(btcd)
	btc.GetBalance(btcwallet)

	btc.ShutdownClients(btcd, btcwallet)

	btc.InterruptProcesses(btcdCmd, btcwalletCmd)

	// p2p.P2P()
	// gateway()
	// server()
}
