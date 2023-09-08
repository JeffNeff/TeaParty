package main

import (
	"fmt"

	"github.com/kaspanet/kaspad/domain/consensus/utils/constants"
	kasrpc "github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
)

func main() {
	// daemonClient, tearDown, err := client.Connect("localhost:16110")
	// if err != nil {
	// 	panic(err)
	// }
	// defer tearDown()

	// ctx, cancel := context.WithTimeout(context.Background(), 10)
	// defer cancel()
	// response, err := daemonClient.GetBalance(ctx, &pb.GetBalanceRequest{})
	// if err != nil {
	// 	panic(err)
	// }
	// println(response)

	// use kasrpc
	kasClient, err := kasrpc.NewRPCClient("localhost:16110")
	if err != nil {
		panic(err)
	}
	defer kasClient.Disconnect()

	// var adresses []string
	// adresses = append(adresses, "kaspa:qq6zldlce2zm67v7engqksseay3fsrynpnwsmvz6l7w5eun7a8akslf4u0qjv")

	response, err := kasClient.GetBalanceByAddress("kaspa:qq6zldlce2zm67v7engqksseay3fsrynpnwsmvz6l7w5eun7a8akslf4u0qjv")
	if err != nil {
		panic(err)
	}

	fmt.Println("Balance: ", formatKas(response.Balance))
	fmt.Println(response.Balance)
	// 152693645093
	// 12
	// fmt.Printf("response: %+v", response)

	// for _, balance := range response.Entries {
	// 	fmt.Printf("balance: %+v", balance.Balance)
	// }

	// // use kasClient
	// bc, err := kasClient.GetBlockCount()
	// if err != nil {
	// 	panic(err)
	// }

	// // fmt.Printf("Block count: %+v", bc)
	// adr := "kaspa:qre0sj4xhk9lsuzqzdy4vqak87szn74g2ss726qm3nql600xxfrkywtdtnrs3"
	// bal, err := kasClient.GetBalanceByAddress(adr)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Balance: %+v", bal)

}

func formatKas(amount uint64) string {
	res := "                   "
	if amount > 0 {
		res = fmt.Sprintf("%19.8f", float64(amount)/constants.SompiPerKaspa)
	}
	return res
}

// func (a *warrenadapter) waitAndVerifyKASChain(ctx context.Context, client, client2 *kasrpc.RPCClient, request AccountWatchRequest) {
// 	// create a ticker that ticks every 30 seconds
// 	ticker := time.NewTicker(time.Second * 30)
// 	defer ticker.Stop()

// 	// create a timer that times out after the specified timeout
// 	timer := time.NewTimer(time.Second * time.Duration(request.TimeOut))
// 	defer timer.Stop()

// 	a.logger.Info("Watching for " + request.Account + " to have a payment of " + request.Amount.String() + " on chain " + request.Chain)
// 	// start a for loop that checks the balance of the address
// 	canILive := true
// 	for canILive {
// 		select {
// 		case <-ticker.C:
// 			balance, err := client.GetBalanceByAddress(request.Account)
// 			if err != nil {
// 				a.logger.Error("occured getting balance of " + request.Account + ": " + err.Error())
// 				return
// 			}

// 			a.logger.Infof("balance of %v is %v on chain %v", request.Account, balance.Balance, request.Chain)
// 			// if the balance is equal to the amount, verify with the second RPC server.
// 			bal := new(big.Int).SetUint64(balance.Balance)
// 			if request.Amount.Cmp(bal) == 0 || request.Amount.Cmp(bal) == 1 {
// 				verifiedBalance, err := client2.GetBalanceByAddress(request.Account)
// 				if err != nil {
// 					a.logger.Error("occured getting balance of " + request.Account + ": " + err.Error() + " from the secondary ETH RPC server")
// 					return
// 				}

// 				bal = new(big.Int).SetUint64(verifiedBalance.Balance)
// 				if request.Amount.Cmp(bal) == 0 || request.Amount.Cmp(bal) == 1 {
// 					a.logger.Info("attempting to complete order " + request.TransactionID)
// 					// send a complete order event
// 					awr := &AccountWatchRequestResult{
// 						AccountWatchRequest: request,
// 						Result:              "",
// 					}
// 					a.emitCE(awr)
// 					canILive = false
// 					return
// 				} else {
// 					a.logger.Error("balance of " + request.Account + " is not equal to " + request.Amount.String())
// 					return
// 				}

// 			}

// 		case <-timer.C:
// 			// if the timer times out, return an error
// 			e := fmt.Sprintf("timeout occured waiting for " + request.Account + " to have a payment of " + request.Amount.String())
// 			a.logger.Info(e)
// 			awrr := &AccountWatchRequestResult{
// 				AccountWatchRequest: request,
// 				Result:              e,
// 			}
// 			a.emitCE(awrr)
// 			canILive = false

// 			return
// 		}
// 	}
// }
