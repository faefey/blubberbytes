package btc

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcwallet/wallet"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb"
)

// Function for generating the seed for the wallet
// func generateSeed(length int) []byte {
// 	seed := make([]byte, length)
// 	rand.Read(seed)
// 	return seed
// }

// Function for generating the private passphrase for the wallet
func generatePrivatePassphrase(length int) string {
	passphrase := make([]byte, length)
	rand.Read(passphrase)
	return base64.StdEncoding.EncodeToString(passphrase)
}

func createWallet(walletDir string, net string) {
	//Choose which network parameters to use based on net
	netParams := &chaincfg.MainNetParams
	if net == "simnet" {
		netParams = &chaincfg.SimNetParams
	} else if net == "testnet" {
		netParams = &chaincfg.TestNet3Params
	}

	loader := wallet.NewLoader(netParams, walletDir+"\\"+net, true, 10*time.Second, 250)

	pubPassphraseString := "public"
	pubPassphrase := []byte(pubPassphraseString)
	privPassphraseString := generatePrivatePassphrase(32)

	privPassphrase := []byte(privPassphraseString)
	//seed := generateSeed(32)

	_, err := loader.CreateNewWallet(pubPassphrase, privPassphrase, nil, time.Now())
	if err != nil {
		panic(err)
	}

	if err := loader.UnloadWallet(); err != nil {
		fmt.Printf("error unloading wallet: %v", err)
	}

	log.Println("New wallet successfully created with public passphrase " + pubPassphraseString + " and private " + privPassphraseString)
}
