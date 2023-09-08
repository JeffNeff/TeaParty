package be

import (
	"context"
	"fmt"
	"time"

	kasrpc "github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
)

func (e *ExchangeServer) waitAndVerifyKASChain(ctx context.Context, client, client2 *kasrpc.RPCClient, request AccountWatchRequest) {
	e.logger.Info("Watching for " + request.Account + " to have a payment of " + request.Amount.String() + " on chain " + request.Chain)
	if !e.watch {
		e.logger.Info("dev mode is on, not watching for payment. Returning success")
		awrr := &AccountWatchRequestResult{
			AccountWatchRequest: request,
			Result:              "suceess",
		}

		if err := e.Dispatch(awrr); err != nil {
			e.logger.Error("error dispatching account watch request result: " + err.Error())
		}
		return
	}

	// create a ticker that ticks every 30 seconds
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	if e.dev {
		ticker = time.NewTicker(time.Second * 10)
	}

	// create a timer that times out after the specified timeout
	timer := time.NewTimer(time.Second * time.Duration(request.TimeOut))
	defer timer.Stop()

	// start a for loop that checks the balance of the address
	canILive := true
	for canILive {
		select {
		case <-ticker.C:
			response, err := client.GetBalanceByAddress(request.Account)
			if err != nil {
				e.logger.Error("occured getting balance of " + request.Account + ": " + err.Error())
				break
			}

			rBal := response.Balance

			e.logger.Infof("balance of %v is %v on chain %v", request.Account, rBal, request.Chain)

			if request.Amount.Uint64() <= rBal {
				verifiedBalance, err := client2.GetBalanceByAddress(request.Account)
				if err != nil {
					e.logger.Error("occured getting balance of " + request.Account + ": " + err.Error() + " from the secondary Kas RPC server")
					break
				}

				if request.Amount.Uint64() <= verifiedBalance.Balance {
					e.logger.Info("attempting to complete order " + request.TransactionID)
					// send a complete order event
					awr := &AccountWatchRequestResult{
						AccountWatchRequest: request,
						Result:              "suceess",
					}
					if err := e.Dispatch(awr); err != nil {
						e.logger.Error("error dispatching account watch request result: " + err.Error())
					}
					e.logger.Info("completed order " + request.TransactionID)
					canILive = false
					break
				} else {
					e.logger.Error("balance of " + request.Account + " is not equal to " + request.Amount.String())
					break
				}

			}

		case <-timer.C:
			// if the timer times out, return an error
			response := fmt.Sprintf("timeout occured waiting for " + request.Account + " to have a payment of " + request.Amount.String())
			e.logger.Info(response)
			awrr := &AccountWatchRequestResult{
				AccountWatchRequest: request,
				Result:              response,
			}
			if err := e.Dispatch(awrr); err != nil {
				e.logger.Error("error dispatching account watch request result: " + err.Error())
			}

			canILive = false

			break
		}
	}
	return
}
