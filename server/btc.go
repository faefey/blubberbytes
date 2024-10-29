package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
)

// Create a new RPC client using websockets.
func createClient(port string) (*rpcclient.Client, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:       "localhost:" + port,
		Endpoint:   "ws",
		User:       "user",
		Pass:       "password",
		DisableTLS: true,
	}

	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Start the btcd process.
func startBtcd(debug bool) (*exec.Cmd, error) {
	cmd := exec.Command("./btcd/btcd", "-C", "./conf/btcd.conf", "-a", "130.245.173.221:8333")

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer cmdStdout.Close()

	scanner := bufio.NewScanner(cmdStdout)

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		if debug {
			fmt.Println(scanner.Text())
		}

		if strings.Contains(scanner.Text(), "Syncing to block height") {
			return cmd, nil
		}
	}

	return nil, errors.New("failed to start btcd")
}

// Start the btcwallet process.
func startBtcwallet(debug bool) (*exec.Cmd, error) {
	walletDir := btcutil.AppDataDir("btcwallet", false)
	if _, err := os.Stat(filepath.Join(walletDir, "mainnet/wallet.db")); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("the wallet does not exist, run ./btcwallet/btcwallet --create to initialize and create it")
	}

	cmd := exec.Command("./btcwallet/btcwallet", "-C", "./conf/btcwallet.conf")

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer cmdStdout.Close()

	scanner := bufio.NewScanner(cmdStdout)

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	rpc := false
	wallet := false
	for scanner.Scan() {
		if debug {
			fmt.Println(scanner.Text())
		}

		if strings.Contains(scanner.Text(), "Established connection to RPC server") {
			rpc = true
		} else if strings.Contains(scanner.Text(), "Opened wallet") {
			wallet = true
		}

		if rpc && wallet {
			return cmd, nil
		}
	}

	return nil, errors.New("failed to start btcwallet")
}
