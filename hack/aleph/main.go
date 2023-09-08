package main

import (
	"fmt"
	"net/http"
	"time"

	alephium "github.com/alephium/go-sdk"
)

func main() {
	// create new Configuration
	// curl http://0.0.0.0:12973/addresses/1333nA4p3e1ShxrmVB5rHBzqqqhpz2dsepTzcV8VtWimG/balance
	config := &alephium.Configuration{
		Host: "http://0.0.0.0:12973",
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	// Create a new client
	client := alephium.NewAPIClient(config)

	fmt.Printf("%+v", client)

	// fetch block height
	height, _, err := client.BlocksApi.GetBlockHeight(nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(height)

}
