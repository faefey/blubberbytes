package btc

import (
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
)

func Start(net string, debug bool) (*exec.Cmd, *exec.Cmd, *rpcclient.Client, *rpcclient.Client, error) {
	walletAddrPath := "./btc/walletaddress.txt"

	if _, err := os.Stat(walletAddrPath); errors.Is(err, os.ErrNotExist) {
		walletDir := btcutil.AppDataDir("btcwallet", false)
		createWallet(walletDir, net)

		btcdCmd, btcwalletCmd, btcd, btcwallet, err := startBtc(net, "", debug)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		err = storeAddress(btcwallet, walletAddrPath)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		ShutdownClient(btcd)
		ShutdownClient(btcwallet)

		InterruptCmd(btcwalletCmd)
		InterruptCmd(btcdCmd)
	}

	file, err := os.Open(walletAddrPath)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer file.Close()

	miningaddr, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	btcdCmd, btcwalletCmd, btcd, btcwallet, err := startBtc(net, string(miningaddr), debug)
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
