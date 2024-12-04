package btc

import (
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
