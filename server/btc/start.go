package btc

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"server/database/operations"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
)

func Start(net string, db *sql.DB, debug bool) (*exec.Cmd, *exec.Cmd, *rpcclient.Client, *rpcclient.Client, error) {
	pubPassphrase := "public"
	var privPassphrase string
	fmt.Print("Enter your private passphrase: ")
	_, err := fmt.Scanln(&privPassphrase)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	walletDir := btcutil.AppDataDir("btcwallet", false)

	// err = os.Remove(filepath.Join(walletDir, net+"/wallet.db"))
	// if err != nil && !os.IsNotExist(err) {
	// 	return nil, nil, nil, nil, err
	// }

	// fmt.Println("Import your wallet (not done yet)")

	if _, err := os.Stat(filepath.Join(walletDir, net+"/wallet.db")); errors.Is(err, os.ErrNotExist) {
		err := createWallet(walletDir, net, pubPassphrase, privPassphrase, db)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}

	err = operations.UpdateWalletPassphrases(db, pubPassphrase, privPassphrase)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	walletInfo, err := operations.GetWalletInfo(db)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	address := walletInfo.Address

	if address == "" {
		btcdCmd, btcwalletCmd, btcd, btcwallet, err := startBtc(net, "", debug)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		address, err = operations.StoreAddress(btcwallet, db)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		ShutdownClient(btcd)
		ShutdownClient(btcwallet)

		InterruptCmd(btcwalletCmd)
		InterruptCmd(btcdCmd)
	}

	btcdCmd, btcwalletCmd, btcd, btcwallet, err := startBtc(net, address, debug)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return btcdCmd, btcwalletCmd, btcd, btcwallet, nil
}

// Start all btc-related processes.
func startBtc(net string, miningaddr string, debug bool) (*exec.Cmd, *exec.Cmd, *rpcclient.Client, *rpcclient.Client, error) {
	btcdCmd, err := startBtcd(net, miningaddr, debug)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	btcwalletCmd, err := startBtcwallet(net, debug)
	if err != nil {
		InterruptCmd(btcdCmd)
		return nil, nil, nil, nil, err
	}

	btcd, err := createBtcdClient(net)
	if err != nil {
		InterruptCmd(btcwalletCmd)
		InterruptCmd(btcdCmd)
		return nil, nil, nil, nil, err
	}

	btcwallet, err := createBtcwalletClient(net)
	if err != nil {
		ShutdownClient(btcd)
		InterruptCmd(btcwalletCmd)
		InterruptCmd(btcdCmd)
		return nil, nil, nil, nil, err
	}

	return btcdCmd, btcwalletCmd, btcd, btcwallet, nil
}
