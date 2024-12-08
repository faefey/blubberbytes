package btc

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Start the btcd process.
func startBtcd(net string, miningaddr string, debug bool) (*exec.Cmd, error) {
	netCmd := ""
	if net != "mainnet" {
		netCmd = "--" + net
	}

	publicNode := "130.245.173.221:8333"
	if net == "testnet" {
		publicNode = "130.245.173.221:18333"
	}

	miningaddrCmd := ""
	if miningaddr != "" {
		miningaddrCmd = "--miningaddr=" + miningaddr
	}

	cmd := exec.Command("./btcd/btcd", "-C", "./btc/conf/btcd.conf", netCmd, "--connect="+publicNode, miningaddrCmd)

	cmd.SysProcAttr = sysProcAttr

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(cmdStdout)

	defer func() {
		go func() {
			for scanner.Scan() {
				if debug {
					fmt.Println(scanner.Text())
				}
			}
		}()
	}()

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		if debug {
			fmt.Println(scanner.Text())
		}

		if net == "mainnet" || net == "testnet" {
			if strings.Contains(scanner.Text(), "Syncing to block height") {
				return cmd, nil
			}
		} else {
			if strings.Contains(scanner.Text(), "RPC server listening") {
				return cmd, nil
			}
		}
	}

	return nil, errors.New("failed to start btcd")
}

// Start the btcwallet process.
func startBtcwallet(net string, debug bool) (*exec.Cmd, error) {
	netCmd := ""
	if net != "mainnet" {
		netCmd = "--" + net
	}

	cmd := exec.Command("./btcwallet/btcwallet", "-C", "./btc/conf/btcwallet.conf", netCmd)

	cmd.SysProcAttr = sysProcAttr

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(cmdStdout)

	defer func() {
		go func() {
			for scanner.Scan() {
				if debug {
					fmt.Println(scanner.Text())
				}
			}
		}()
	}()

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
