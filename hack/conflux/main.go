package main

import (
	"fmt"

	sdk "github.com/Conflux-Chain/go-conflux-sdk"
)

func main() {
	client, err := sdk.NewClient("http://0.0.0.0:12539", sdk.ClientOption{
		KeystorePath: "../context/keystore",
	})

	if err != nil {
		panic(err)
	}

	// get epoc
	epoch, err := client.GetEpochNumber()
	if err != nil {
		panic(err)
	}

	fmt.Println(epoch)

	chainID, err := client.GetNetworkID()
	if err != nil {
		panic(err)
	}
	fmt.Println(chainID)

}
