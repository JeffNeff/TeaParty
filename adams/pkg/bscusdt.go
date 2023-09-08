package be

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type BalanceResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

func fetchBNBUSDTBalance(address string) (*big.Int, error) {
	// remove "0x" from the address
	address = strings.Replace(address, "0x", "", -1)
	url := "https://magical-still-snowflake.bsc.discover.quiknode.pro/865ae0a9a366e32ef1f4a8ad7cccf67bfa59a661/"
	method := "POST"

	payload := strings.NewReader(`{
	  "id":67,
	  "jsonrpc":"2.0",
	  "method":"eth_call",
	  "params":[{"data":"0x70a08231000000000000000000000000` + address + `","to":"0x55d398326f99059fF775485246999027B3197955"}, "latest"]
	}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-qn-api-version", "1")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var balanceResponse BalanceResponse
	err = json.Unmarshal(body, &balanceResponse)
	if err != nil {
		return nil, err
	}

	// convert the result from a hex string to a decimal string
	dec, err := strconv.ParseInt(balanceResponse.Result, 0, 64)
	if err != nil {
		return nil, err
	}

	return big.NewInt(dec), nil
}

func (a *ExchangeServer) waitAndVerifyBSCUSDT(request AccountWatchRequest) {
	ctx := context.Background()
	if !a.watch {
		a.logger.Info("dev mode is on, not watching for payment. Returning success")
		awrr := &AccountWatchRequestResult{
			AccountWatchRequest: request,
			Result:              "suceess",
		}

		if err := a.Dispatch(awrr); err != nil {
			a.logger.Error("error dispatching account watch request result: " + err.Error())
		}
		return
	}
	a.logger.Info("Watching for " + request.Account + " to have a payment of " + request.Amount.String() + " with BSCUSDT ")

	// create a ticker that ticks every 30 seconds
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
	if a.dev {
		ticker = time.NewTicker(time.Second * 10)
	}

	// create a timer that times out after the specified timeout
	timer := time.NewTimer(time.Second * time.Duration(request.TimeOut))
	defer timer.Stop()
	account := common.HexToAddress(request.Account)
	// start a for loop that checks the balance of the address
	canILive := true
	for canILive {
		select {
		case <-ctx.Done():
			a.logger.Info("context done, exiting")
			awrr := &AccountWatchRequestResult{
				AccountWatchRequest: request,
				Result:              "error",
			}
			if err := a.Dispatch(awrr); err != nil {
				a.logger.Error("error dispatching account watch request result: " + err.Error())
			}
			return
		case <-ticker.C:
			// get the balance of the address
			balance, err := fetchBNBUSDTBalance(request.Account)
			if err != nil {
				a.logger.Error("error getting balance of address: " + err.Error())
				continue
			}
			a.logger.Infof("balance of %v is %v on chain %v", account, balance, request.Chain)
			// if the balance is greater than the amount, return success
			if balance.Cmp(request.Amount) >= 0 {
				a.logger.Info("payment received, returning success")
				awrr := &AccountWatchRequestResult{
					AccountWatchRequest: request,
					Result:              "suceess",
				}

				if err := a.Dispatch(awrr); err != nil {
					a.logger.Error("error dispatching account watch request result: " + err.Error())
				}
				canILive = false
				return
			}
		case <-timer.C:
			a.logger.Info("timeout reached, returning failure")
			awrr := &AccountWatchRequestResult{
				AccountWatchRequest: request,
				Result:              "error",
			}
			if err := a.Dispatch(awrr); err != nil {
				a.logger.Error("error dispatching account watch request result: " + err.Error())
				MetricsFailedAccountWatchRequestIncrement()
			}
			canILive = false
			return
		}
	}
}
