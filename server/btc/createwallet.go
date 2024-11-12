package btc

import (
	"log"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcwallet/wallet"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb"
)

func createWallet(walletDir string, net string) {
	//Choose which network parameters to use based on net
	netParams := &chaincfg.MainNetParams
	if net == "simnet" {
		netParams = &chaincfg.SimNetParams
	} else if net == "testnet" {
		netParams = &chaincfg.TestNet3Params
	}

	loader := wallet.NewLoader(netParams, walletDir+"\\"+net, true, 10*time.Second, 250)

	pubPassphraseString := "public_passphrase"
	pubPassphrase := []byte(pubPassphraseString)
	privPassphraseString := generatePrivatePassphrase(32)

	privPassphrase := []byte(privPassphraseString)
	seed := generateSeed(32)

	_, err := loader.CreateNewWallet(pubPassphrase, privPassphrase, seed, time.Now())
	if err != nil {
		panic(err)
	}

	log.Println("New wallet successfully created with public passphrase " + pubPassphraseString + " and private " + privPassphraseString)
}
