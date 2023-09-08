package main

import (
	"fmt"

	btcRPC "github.com/btcsuite/btcd/rpcclient"
	"github.com/google/uuid"
)

// //create wallet
// curl --user dockeruser:dockerpass   --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "createwallet", "params": ["ts"]}' -H 'content-type: text/plain;' 192.168.50.8:7332
// // get new address
// curl --user dockeruser:dockerpass  --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "getnewaddress", "params": []}' -H 'content-type: text/plain;' 192.168.50.8:7332/wallet/ts
// // get private key
// curl --user dockeruser:dockerpass --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "dumpprivkey", "params": ["13WSwSibLwgXXiwa6jyEEmazzu7Ygx56W3"]}' -H 'content-type: text/plain;' http://192.168.50.8:7332/wallet/ts
// // get balance
// curl  --user dockeruser:dockerpass --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "getbalance", "params": ["*", 6]}' -H 'content-type: text/plain;' http://192.168.50.8:7332/wallet/ts

// curl --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "getnewaddress", "params": []}' -H 'content-type: text/plain;' http://192.168.50.8:7332/wallet/testwallet

// gpt get balance
 curl --user "dockeruser":"dockerpass" \
 --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "getbalance", "params": ["1EKJZUuyF6GzK6xT51LATSFDpTXzzeJiVm"]}' \
     -H 'content-type: text/plain;' \
     http://192.168.50.8:7332/wallet


func main() {

	walletName := uuid.New().String()
	var connCfg = &btcRPC.ConnConfig{
		Host: "localhost:7332/wallet/" + walletName,
		// Host: "https://newest-purple-smoke.btc-testnet.discover.quiknode.pro/ea6fd836108a00c7cb4a62b5e491f53c11b4476b",
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

	address, err := btcClient.GetNewAddress(walletName)
	if err != nil {
		panic(err)
	}
	fmt.Println(address)

	// // get wallet balance
	// balance, err := btcClient.GetBalance("*")
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(balance)

	// walletName = uuid.New().String()
	// connCfg = &btcRPC.ConnConfig{
	// 	Host: "localhost:7332",
	// 	// Host: "localhost:7332/wallet/" + walletName,
	// 	// Host: "https://newest-purple-smoke.btc-testnet.discover.quiknode.pro/ea6fd836108a00c7cb4a62b5e491f53c11b4476b",
	// 	User:         "dockeruser",
	// 	Pass:         "dockerpass",
	// 	HTTPPostMode: true,
	// 	DisableTLS:   true,
	// }

	// btcClient, err = btcRPC.New(connCfg, nil)
	// if err != nil {
	// 	fmt.Println("err creating client: ", err)
	// }

	// wal, err := btcClient.CreateWallet(
	// 	walletName,
	// )
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("%+v", wal)

	// // get wallet balance
	// balance, err = btcClient.GetBalance("*")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(balance)

	// address, err = btcClient.GetNewAddress(walletName)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(address)

}

func test() {
	url := "http://192.168.50.8:7332/wallet/ts"
	payload := []byte(`{"jsonrpc": "1.0", "id": "curltest", "method": "getbalance", "params": ["*", 6]}`)
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	
	req.Header.Set("Content-Type", "text/plain")
	
	auth := []byte("dockeruser:dockerpass")
	b64Auth := base64.StdEncoding.EncodeToString(auth)
	req.Header.Set("Authorization", "Basic "+b64Auth)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	
	defer resp.Body.Close()
	
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}
	
	fmt.Println(result)
}