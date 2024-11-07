package btc

import (
	"github.com/btcsuite/btcd/rpcclient"
)

// Create a new RPC client using websockets.
func createClient(port string, net string) (*rpcclient.Client, error) {
	netParam := net
	if net == "testnet" {
		netParam = "testnet3"
	}

	connCfg := &rpcclient.ConnConfig{
		Host:       "localhost:" + port,
		Endpoint:   "ws",
		User:       "user",
		Pass:       "password",
		DisableTLS: true,
		Params:     netParam,
	}

	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Create a new RPC client for btcd using websockets.
func createBtcdClient(net string) (*rpcclient.Client, error) {
	return createClient("8334", net)
}

// Create a new RPC client for btcwallet using websockets.
func createBtcwalletClient(net string) (*rpcclient.Client, error) {
	return createClient("8332", net)
}

// Shutdown a client.
func ShutdownClient(client *rpcclient.Client) {
	client.Shutdown()
	client.WaitForShutdown()
}
