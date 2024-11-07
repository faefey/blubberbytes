package btc

import (
	"io"
	"os"
	"os/exec"

	"github.com/btcsuite/btcd/rpcclient"
)

func Start(net string, debug bool) (*exec.Cmd, *exec.Cmd, *rpcclient.Client, *rpcclient.Client, error) {
	init := true
	file, err := os.Open("walletaddress.txt")
	if err != nil {
		init = false
	}
	init = false

	var btcdCmd, btcwalletCmd *exec.Cmd
	var btcd, btcwallet *rpcclient.Client

	if init {
		miningaddr, err := io.ReadAll(file)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		btcdCmd, btcwalletCmd, btcd, btcwallet, err = startBtc(net, string(miningaddr), debug)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	} else {
		btcdCmd, btcwalletCmd, btcd, btcwallet, err = startBtc(net, "", debug)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		miningaddr, err := storeAddress(btcwallet)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		ShutdownClient(btcd)
		ShutdownClient(btcwallet)

		InterruptCmd(btcwalletCmd)
		InterruptCmd(btcdCmd)

		btcdCmd, btcwalletCmd, btcd, btcwallet, err = startBtc(net, miningaddr.String(), debug)
		if err != nil {
			return nil, nil, nil, nil, err
		}
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
