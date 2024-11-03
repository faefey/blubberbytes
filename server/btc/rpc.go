package btc

import (
	"fmt"
	"log"
	"os"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
)

// Store the wallet address.
func storeAddress(btcwallet *rpcclient.Client) (btcutil.Address, error) {
	// Query the RPC server for the wallet address.
	address, err := btcwallet.GetAccountAddress("default")
	if err != nil {
		return nil, err
	}

	// Store wallet address for transactions.
	err = os.WriteFile("walletaddress.txt", []byte(address.String()), 0755)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func GetBlockCount(btcd *rpcclient.Client) {
	// Query the RPC server for the current block count and display it.
	blocks, err := btcd.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Block count:", blocks)
}

func GetBalance(btcwallet *rpcclient.Client) {
	// Query the RPC server for the current wallet balance and display it.
	balance, err := btcwallet.GetBalance("default")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Balance:", balance)
}
