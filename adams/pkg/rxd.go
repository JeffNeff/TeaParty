package be

type RPCRequest struct {
	Jsonrpc string   `json:"jsonrpc"`
	ID      string   `json:"id"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

type CreateWalletResponse struct {
	Result struct {
		Name    string `json:"name"`
		Warning string `json:"warning"`
	} `json:"result"`
	Error map[string]interface{} `json:"error"`
	ID    string                 `json:"id"`
}

type CreateAddressRespone struct {
	Result string                 `json:"result"`
	Error  map[string]interface{} `json:"error"`
	ID     string                 `json:"id"`
}

type AccountResponse struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

type BTCResponse struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publickey"`
}

type AccountRequest struct {
	Name string `json:"name"`
}

// func (e *ExchangeServer) waitAndVerifyRXD(request AccountWatchRequest) {
// 	a.logger.Infof("waiting for %v to have a payment of %v on chain %v", request.Account, request.Amount, request.Chain)
// 	if !a.watch {
// 		a.logger.Info("dev mode is on, not watching for payment. Returning success")
// 		awrr := &AccountWatchRequestResult{
// 			AccountWatchRequest: request,
// 			Result:              "suceess",
// 		}

// 		if err := a.Dispatch(awrr); err != nil {
// 			a.logger.Error("error dispatching account watch request result: " + err.Error())
// 		}
// 	}

// 	config := e.radiantNode.config.Host  + "/wallet/" + request.TransactionID
// 	client, err := btcRPC.New(&config, nil)
// 	if err != nil {
// 		a.logger.Errorw("holy fucknuts batman!! The RPC for BTC-like chain is down!!: " + request.Chain)
// 		return
// 	}

// 		// the request.Amount is currently in ETH big.Int format
// 	// convert it to BTC
// 	amount := btcutil.Amount(request.Amount.Int64())
// 	// create a ticker that ticks every 30 seconds
// 	ticker := time.NewTicker(time.Second * 60)
// 	defer ticker.Stop()
// 	if a.dev {
// 		ticker = time.NewTicker(time.Second * 10)
// 	}

// 	// create a timer that times out after the specified timeout
// 	timer := time.NewTimer(time.Second * time.Duration(request.TimeOut))
// 	defer timer.Stop()

// 	a.logger.Info("Watching for " + request.Account + " to have a payment of " + amount.String() + " on chain " + request.Chain)
// 	// start a for loop that checks the balance of the address
// 	canILive := true
// 	for canILive {
// 		select {
// 		case <-ticker.C:
// 			balance, err := e.fetchRXDBalance(request.TransactionID,)

// }

// func (e *ExchangeServer) fetchRXDBalance(config string) error {
// 	payload := []byte(`{"jsonrpc": "1.0", "id": "curltest", "method": "getbalance", "params": ["*", 6]}`)
// 	req, err := http.NewRequest("POST", config, bytes.NewBuffer(payload))
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("Content-Type", "text/plain")
// 	auth := []byte("dockeruser:dockerpass")

// 	b64Auth := base64.StdEncoding.EncodeToString(auth)
// 	req.Header.Set("Authorization", "Basic "+b64Auth)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}

// 	defer resp.Body.Close()

// 	var result map[string]interface{}
// 	err = json.NewDecoder(resp.Body).Decode(&result)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("balance")
// }

// curl --user dockeruser:dockerpass   --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "createwallet", "params": ["ts"]}' -H 'content-type: text/plain;' 192.168.50.8:7332
// // get new address
// curl --user dockeruser:dockerpass  --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "getnewaddress", "params": []}' -H 'content-type: text/plain;' 192.168.50.8:7332/wallet/ts
// // get private key
// curl --user dockeruser:dockerpass --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "dumpprivkey", "params": ["13WSwSibLwgXXiwa6jyEEmazzu7Ygx56W3"]}' -H 'content-type: text/plain;' http://192.168.50.8:7332/wallet/ts
// // get balance
// curl  --user dockeruser:dockerpass --data-binary '{"jsonrpc": "1.0", "id": "curltest", "method": "getbalance", "params": ["*", 6]}' -H 'content-type: text/plain;' http://192.168.50.8:7332/wallet/ts
