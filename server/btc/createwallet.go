package btc

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcwallet/wallet"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb"
)

func createWallet(walletDir string, net string, pubPassphraseString, privPassphraseString string, db *sql.DB) error {
	//Choose which network parameters to use based on net
	netParams := &chaincfg.MainNetParams
	if net == "simnet" {
		netParams = &chaincfg.SimNetParams
	} else if net == "testnet" {
		netParams = &chaincfg.TestNet3Params
	}

	loader := wallet.NewLoader(netParams, filepath.Join(walletDir, net), true, 10*time.Second, 250)

	pubPassphrase := []byte(pubPassphraseString)

	privPassphrase := []byte(privPassphraseString)

	_, err := loader.CreateNewWallet(pubPassphrase, privPassphrase, nil, time.Now())
	if err != nil {
		return fmt.Errorf("error creating wallet: %v", err)
	}

	if err := loader.UnloadWallet(); err != nil {
		return fmt.Errorf("error unloading wallet: %v", err)
	}

	log.Println("New wallet successfully created")

	return nil
}
