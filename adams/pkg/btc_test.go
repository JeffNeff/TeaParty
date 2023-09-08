package be

import (
	"fmt"
	"testing"

	btcRPC "github.com/btcsuite/btcd/rpcclient"
)

func TestGenerateBTCAccount(t *testing.T) {
	walletName := "for7ae0758b-5cca-441a-932d-bebc82ad2b26"
	var connCfg = &btcRPC.ConnConfig{
		Host:         "192.168.50.8:7332",
		User:         "dockeruser",
		Pass:         "dockerpass",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	btcClient, err := btcRPC.New(connCfg, nil)
	if err != nil {
		fmt.Println("err creating client: ", err)
	}

	_, err = btcClient.CreateWallet(
		walletName,
	)
	if err != nil {
		panic(err)
	}

	// update the connection config to use the new wallet
	connCfg.Host = connCfg.Host + "/wallet/" + walletName
	address, err := btcClient.GetNewAddress(walletName)
	if err != nil {
		panic(err)
	}

	t.Log(address)

	// get wallet balance
	balance, err := btcClient.GetBalance("*")
	if err != nil {
		panic(err)
	}

	t.Log(balance)

	// get wallet private key
	privateKey, err := btcClient.DumpPrivKey(address)
	if err != nil {
		panic(err)
	}

	t.Log(privateKey)
}
