package btc

import (
	"os"

	"github.com/btcsuite/btcd/rpcclient"
)

// Store the wallet address.
func storeAddress(btcwallet *rpcclient.Client, path string) error {
	// Query the RPC server for the wallet address.
	address, err := btcwallet.GetAccountAddress("default")
	if err != nil {
		return err
	}

	// Store wallet address for transactions.
	err = os.WriteFile(path, []byte(address.String()), 0755)
	if err != nil {
		return err
	}

	return nil
}
