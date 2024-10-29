package main

import (
	"fmt"
	"log"
	"syscall"
)

func main() {
	btcdCmd, err := startBtcd(false)
	if err != nil {
		log.Println(err)
		return
	}

	btcwalletCmd, err := startBtcwallet(false)
	if err != nil {
		log.Println(err)
		return
	}

	btcd, err := createClient("8334")
	if err != nil {
		log.Println(err)
		return
	}

	btcwallet, err := createClient("8332")
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		btcd.Shutdown()
		btcwallet.Shutdown()
		btcd.WaitForShutdown()
		btcwallet.WaitForShutdown()

		err = btcdCmd.Process.Signal(syscall.SIGKILL)
		if err != nil {
			log.Println(err)
		}
		err = btcwalletCmd.Process.Signal(syscall.SIGKILL)
		if err != nil {
			log.Println(err)
		}
	}()

	// Query the RPC server for the current block count and display it.
	blocks, err := btcd.GetBlockCount()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Block count:", blocks)

	// Query the RPC server for the current wallet balance and display it.
	balance, err := btcwallet.GetBalance("default")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Balance:", balance)

	// Query the RPC server for the wallet address and display it.
	address, err := btcwallet.GetAccountAddress("default")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Address:", address)

	// Mine coins for the wallet and display the mined block.
	// blockHashes, err := btcwallet.SendFrom("default", address, 1)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// fmt.Println("Balance:", blockHashes)

	// Run the server.
	// server()
}
