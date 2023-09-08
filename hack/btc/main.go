package main

import (
	"fmt"

	btcRPC "github.com/btcsuite/btcd/rpcclient"
)

func main() {
	var connCfg = &btcRPC.ConnConfig{
		Host: "https://newest-purple-smoke.btc-testnet.discover.quiknode.pro/ea6fd836108a00c7cb4a62b5e491f53c11b4476b",
		// User:         "",
		// Pass:         "",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	btcClient, err := btcRPC.New(connCfg, nil)
	if err != nil {
		fmt.Println("err creating client: ", err)
	}

	// Get the current block count.
	blockCount, err := btcClient.GetBlockCount()
	if err != nil {
		panic(err)
	}

	fmt.Println(blockCount)

}
