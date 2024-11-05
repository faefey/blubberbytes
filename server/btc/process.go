package btc

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/btcsuite/btcd/btcutil"
	"golang.org/x/sys/windows"
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

	cmd := exec.Command("./btcd/btcd", "-C", "./conf/btcd.conf", netCmd, "-a", publicNode, miningaddrCmd)

	cmd.SysProcAttr = &windows.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_PROCESS_GROUP,
	}

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

		if net == "mainnet" {
			if strings.Contains(scanner.Text(), "Syncing to block height") {
				return cmd, nil
			}
		} else {
			if strings.Contains(scanner.Text(), "Server listening") {
				return cmd, nil
			}
		}
	}

	return nil, errors.New("failed to start btcd")
}

// Start the btcwallet process.
func startBtcwallet(net string, debug bool) (*exec.Cmd, error) {
	walletDir := btcutil.AppDataDir("btcwallet", false)
	if _, err := os.Stat(filepath.Join(walletDir, net+"/wallet.db")); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("the wallet does not exist, run ./btcwallet/btcwallet --create to initialize and create it")
	}

	netCmd := ""
	if net != "mainnet" {
		netCmd = "--" + net
	}

	cmd := exec.Command("./btcwallet/btcwallet", "-C", "./conf/btcwallet.conf", netCmd)

	cmd.SysProcAttr = &windows.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_PROCESS_GROUP,
	}

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

// Interrupt exec.Cmd processes.
func InterruptProcesses(cmds ...*exec.Cmd) {
	for _, cmd := range cmds {
		if runtime.GOOS == "windows" {
			err := sendCtrlBreak(cmd.Process.Pid)
			if err != nil {
				log.Println(err)
			}
		} else {
			err := cmd.Process.Signal(os.Interrupt)
			if err != nil {
				log.Println(err)
			}
		}

		err := cmd.Wait()
		if err != nil {
			log.Println(err)
		}
	}
}

func sendCtrlBreak(pid int) error {
	d, err := windows.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}
	defer d.Release()

	p, err := d.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return err
	}

	r, _, err := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
	if r == 0 {
		return err
	}

	return nil
}
