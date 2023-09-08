package be

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	bridge "github.com/teapartycrypto/TeaParty/adams/pkg/contract/bridge"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (e *ExchangeServer) generateEVMAccount(chain string) *ecdsa.PrivateKey {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	pk := hexutil.Encode(privateKeyBytes)[2:]
	e.logger.Debug("Generated " + chain + " Private Key: " + pk)
	return privateKey
}

const (
	MININGGAMECONTRACTADDRESS = "0x970A8b10147E3459D3CBF56329B76aC18D329728"
)

func (a *ExchangeServer) waitAndVerifyThatNFTIsAvalibleOnEVMChain(request AccountWatchRequest) {
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

	a.logger.Info("Watching for the NFT to be avalible on " + request.Chain)

	// create a ticker that ticks every 30 seconds
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()
	if a.dev {
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
			// check the balance of the address
			a.logger.Info("Checking if the NFT is avalible at address " + request.Account)
			contractAddress := common.HexToAddress(MININGGAMECONTRACTADDRESS)

			awrr, err := a.defineTheOwnerOfNFT(ctx, request.Chain, contractAddress.String(), request.Account, request.NFTID, &request)
			if err != nil {
				a.logger.Error("error checking balance: " + err.Error())
				continue
			}

			if awrr != nil {
				if err := a.Dispatch(awrr); err != nil {
					a.logger.Error("error dispatching account watch request result: " + err.Error())
				}
				canILive = false
			}

		case <-timer.C:
			// the timer has timed out, return an error
			a.logger.Info("Timeout reached, returning error")
			awrr := &AccountWatchRequestResult{
				AccountWatchRequest: request,
				Result:              "error",
			}

			if err := a.Dispatch(awrr); err != nil {
				a.logger.Error("error dispatching account watch request result: " + err.Error())
			}
			canILive = false

		case <-ctx.Done():
			// the context has been canceled, return an error
			a.logger.Info("Context canceled, returning error")
			awrr := &AccountWatchRequestResult{
				AccountWatchRequest: request,
				Result:              "error",
			}

			if err := a.Dispatch(awrr); err != nil {
				a.logger.Error("error dispatching account watch request result: " + err.Error())
			}

			canILive = false
		}
	}

}

func (a *ExchangeServer) defineTheOwnerOfNFT(ctx context.Context, chain, contract, account string, id int64, awr *AccountWatchRequest) (*AccountWatchRequestResult, error) {
	// account = "0x68bd627A508441011511eC2980a113C06eDBf965"
	// url := "https://nft.api.infura.io/networks/137/nfts/0x970A8b10147E3459D3CBF56329B76aC18D329728/4/owners"
	// convert the id of type int64 to a string
	idstring := strconv.FormatInt(id, 10)
	// we need to know what network the nft is on
	var chainid string
	switch chain {
	case MiningGame:
		chainid = "137"
	default:
		return nil, fmt.Errorf("chain not supported: %s", chain)
	}

	url := "https://nft.api.infura.io/networks/" + chainid + "/nfts/" + contract + "/" + idstring + "/owners"
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic NjE5Nzk3OTdhOGJiNGJmZTlkZGRkNGZmOTY3NWRiN2U6ODBjM2ZjZjk2ZmFhNDlmN2JkNmIzMmNjYzQxOTAyZDg=")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var t Token
	err = json.Unmarshal(body, &t)
	if err != nil {
		return nil, err
	}

	fmt.Printf("comparig %s to %s", strings.ToLower(strings.TrimSpace(t.Owners[0].OwnerOf)), strings.ToLower(strings.TrimSpace(account)))
	if strings.Contains(strings.ToLower(strings.TrimSpace(t.Owners[0].OwnerOf)), strings.ToLower(strings.TrimSpace(account))) {
		// compare the address to the address we are looking for strings

		a.logger.Info("NFT is avalible")
		awrr := &AccountWatchRequestResult{
			AccountWatchRequest: *awr,
			Result:              "suceess",
		}

		a.logger.Infof("Dispatching account watch request result: %+v", awrr)

		return awrr, nil
	}

	a.logger.Info("NFT is not avalible yet")
	return nil, nil
}

func (a *ExchangeServer) waitAndVerifyEVMChain(ctx context.Context, client, client2 *ethclient.Client, request AccountWatchRequest) {
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
	a.logger.Info("Watching for " + request.Account + " to have a payment of " + request.Amount.String() + " on chain " + request.Chain)

	// create a ticker that ticks every 60 seconds
	ticker := time.NewTicker(time.Second * 60)
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
		case <-ticker.C:
			balance, err := client.BalanceAt(context.Background(), account, nil)
			if err != nil {
				a.logger.Error("occured getting balance of " + request.Account + ": " + err.Error())
				return
			}
			a.logger.Infof("balance of %v is %v on chain %v", account, balance, request.Chain)
			// if the balance is equal to the amount, verify with the
			// second RPC server.
			if balance.Cmp(request.Amount) == 0 || balance.Cmp(request.Amount) == 1 {
				verifiedBalance, err := client2.BalanceAt(context.Background(), account, nil)
				if err != nil {
					a.logger.Error("occured getting balance of " + request.Account + ": " + err.Error() + " from the secondary ETH RPC server")
					return
				}

				if verifiedBalance.Cmp(request.Amount) == 0 || verifiedBalance.Cmp(request.Amount) == 1 {
					a.logger.Info("attempting to complete order " + request.TransactionID)
					// send a complete order event
					awrr := &AccountWatchRequestResult{
						AccountWatchRequest: request,
						Result:              "suceess",
					}

					if err := a.Dispatch(awrr); err != nil {
						a.logger.Error("error dispatching account watch request result: " + err.Error())
					}
					canILive = false
					return
				} else {
					a.logger.Error("balance of " + request.Account + " is not equal to " + request.Amount.String())
					return
				}
			}
		case <-timer.C:
			// if the timer times out, return an error
			e := fmt.Sprintf("timeout occured waiting for " + request.Account + " to have a payment of " + request.Amount.String())
			a.logger.Info(e)
			awrr := &AccountWatchRequestResult{
				AccountWatchRequest: request,
				Result:              "error",
			}

			if err := a.Dispatch(awrr); err != nil {
				a.logger.Error("error dispatching account watch request result: " + err.Error())
			}

			canILive = false

			return
		}
	}

}

func (a *ExchangeServer) waitAndVerifyWGRAMSBridgeTokenOnOctaSpace(ctx context.Context, client, client2 *ethclient.Client, request AccountWatchRequest) {
	// request.Account = "0x5D22D5c8675d3e3a6a1f296d740d6381CbD18769"
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
	a.logger.Info("Watching for " + request.Account + " to have a payment of " + request.Amount.String() + " tokens on chain " + request.Chain)

	// create a ticker that ticks every 60 seconds
	ticker := time.NewTicker(time.Second * 60)
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
		case <-ticker.C:
			balance, err := a.queryWGRAMSBridgeContractOnOctaSpaceUserAccountBalance(request.Account, client)
			if err != nil {
				a.logger.Error("occured getting balance of " + request.Account + ": " + err.Error())
				return
			}
			a.logger.Infof("balance of %v is %v on chain %v", account, balance, request.Chain)
			// if the balance is equal to the amount, verify with the
			// second RPC server.
			if balance.Cmp(request.Amount) == 0 || balance.Cmp(request.Amount) == 1 {
				verifiedBalance, err := a.queryWGRAMSBridgeContractOnOctaSpaceUserAccountBalance(request.Account, client)
				if err != nil {
					a.logger.Error("occured getting balance of " + request.Account + ": " + err.Error() + " from the secondary ETH RPC server")
					return
				}

				if verifiedBalance.Cmp(request.Amount) == 0 || verifiedBalance.Cmp(request.Amount) == 1 {
					a.logger.Info("attempting to complete order " + request.TransactionID)
					// send a complete order event
					awrr := &AccountWatchRequestResult{
						AccountWatchRequest: request,
						Result:              "suceess",
					}

					if err := a.Dispatch(awrr); err != nil {
						a.logger.Error("error dispatching account watch request result: " + err.Error())
					}
					canILive = false
					return
				} else {
					a.logger.Error("balance of " + request.Account + " is not equal to " + request.Amount.String())
					return
				}
			}
		case <-timer.C:
			// if the timer times out, return an error
			e := fmt.Sprintf("timeout occured waiting for " + request.Account + " to have a payment of " + request.Amount.String())
			a.logger.Info(e)
			awrr := &AccountWatchRequestResult{
				AccountWatchRequest: request,
				Result:              "error",
			}

			if err := a.Dispatch(awrr); err != nil {
				a.logger.Error("error dispatching account watch request result: " + err.Error())
			}

			canILive = false

			return
		}
	}

}

func (a *ExchangeServer) waitAndVerifyWOCTABridgeTokenOnPartychain(ctx context.Context, client, client2 *ethclient.Client, request AccountWatchRequest) {
	// request.Account = "0x5D22D5c8675d3e3a6a1f296d740d6381CbD18769"
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
	a.logger.Info("Watching for " + request.Account + " to have a payment of " + request.Amount.String() + " tokens on chain " + request.Chain)

	// create a ticker that ticks every 60 seconds
	ticker := time.NewTicker(time.Second * 60)
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
		case <-ticker.C:
			balance, err := a.queryWOCTABridgeContractOnPartyChainUserAccountBalance(request.Account, client)
			if err != nil {
				a.logger.Error("occured getting balance of " + request.Account + ": " + err.Error())
				return
			}
			a.logger.Infof("balance of %v is %v on chain %v", account, balance, request.Chain)
			// if the balance is equal to the amount, verify with the
			// second RPC server.
			if balance.Cmp(request.Amount) == 0 || balance.Cmp(request.Amount) == 1 {
				verifiedBalance, err := a.queryWOCTABridgeContractOnPartyChainUserAccountBalance(request.Account, client)
				if err != nil {
					a.logger.Error("occured getting balance of " + request.Account + ": " + err.Error() + " from the secondary ETH RPC server")
					return
				}

				if verifiedBalance.Cmp(request.Amount) == 0 || verifiedBalance.Cmp(request.Amount) == 1 {
					a.logger.Info("attempting to complete order " + request.TransactionID)
					// send a complete order event
					awrr := &AccountWatchRequestResult{
						AccountWatchRequest: request,
						Result:              "suceess",
					}

					if err := a.Dispatch(awrr); err != nil {
						a.logger.Error("error dispatching account watch request result: " + err.Error())
					}
					canILive = false
					return
				} else {
					a.logger.Error("balance of " + request.Account + " is not equal to " + request.Amount.String())
					return
				}
			}
		case <-timer.C:
			// if the timer times out, return an error
			e := fmt.Sprintf("timeout occured waiting for " + request.Account + " to have a payment of " + request.Amount.String())
			a.logger.Info(e)
			awrr := &AccountWatchRequestResult{
				AccountWatchRequest: request,
				Result:              "error",
			}

			if err := a.Dispatch(awrr); err != nil {
				a.logger.Error("error dispatching account watch request result: " + err.Error())
			}

			canILive = false

			return
		}
	}

}

const (
	// WGRAMSOnOCTAAddress is the address of the WGRAMS bridge contract on the OCTA chain
	WGRAMSOnOCTAAddress      = "0x0eeAaF074B23942CD660175dEaE6e1A5849d6614"
	WOCTAOnPartyChainAddress = "0x01c8024597EF12c6aC4954eC415b9e01c9Fe5976"
)

func (e *ExchangeServer) queryWGRAMSBridgeContractOnOctaSpaceUserAccountBalance(account string, rpc *ethclient.Client) (*big.Int, error) {
	e.logger.Info("querying contract " + WGRAMSOnOCTAAddress + " for balance of " + account)

	contract, err := bridge.NewBe(common.HexToAddress(WGRAMSOnOCTAAddress), rpc)
	if err != nil {
		fmt.Println("error encounterd creating contract instance: " + err.Error())
		return nil, err
	}

	balance, err := contract.BalanceOf(nil, common.HexToAddress(account))
	if err != nil {
		fmt.Println("error encounterd creating contract instance: " + err.Error())
		return nil, err
	}

	fmt.Println("balance of " + account + " is " + balance.String())

	return balance, nil
}

func (e *ExchangeServer) queryWOCTABridgeContractOnPartyChainUserAccountBalance(account string, rpc *ethclient.Client) (*big.Int, error) {
	e.logger.Info("querying contract " + WOCTAOnPartyChainAddress + " for balance of " + account)

	contract, err := bridge.NewBe(common.HexToAddress(WOCTAOnPartyChainAddress), rpc)
	if err != nil {
		fmt.Println("error encounterd creating contract instance: " + err.Error())
		return nil, err
	}

	balance, err := contract.BalanceOf(nil, common.HexToAddress(account))
	if err != nil {
		fmt.Println("error encounterd creating contract instance: " + err.Error())
		return nil, err
	}

	fmt.Println("balance of " + account + " is " + balance.String())

	return balance, nil
}

// func (e *ExchangeServer) waitAndVerifyTokenOnEVMChain(ctx context.Context, client, client2 *ethclient.Client, request AccountWatchRequest) {
// 	// request.Account = "0x5D22D5c8675d3e3a6a1f296d740d6381CbD18769"
// 	if !e.watch {
// 		e.logger.Info("dev mode is on, not watching for payment. Returning success")
// 		awrr := &AccountWatchRequestResult{
// 			AccountWatchRequest: request,
// 			Result:              "suceess",
// 		}

// 		if err := e.Dispatch(awrr); err != nil {
// 			e.logger.Error("error dispatching account watch request result: " + err.Error())
// 		}
// 		return
// 	}
// 	e.logger.Info("Watching for " + request.Account + " to have a payment of " + request.Amount.String() + " on chain " + request.Chain)

// 	// create a ticker that ticks every 60 seconds
// 	ticker := time.NewTicker(time.Second * 60)
// 	defer ticker.Stop()
// 	if e.dev {
// 		ticker = time.NewTicker(time.Second * 10)
// 	}

// 	// create a timer that times out after the specified timeout
// 	timer := time.NewTimer(time.Second * time.Duration(request.TimeOut))
// 	defer timer.Stop()

// 	account := common.HexToAddress(request.Account)

// 	// start a for loop that checks the balance of the address
// 	canILive := true
// 	for canILive {
// 		select {
// 		case <-ticker.C:
// 			balance, err := client.BalanceOf(context.Background(), account, request.ContractAddress)
// 			if err != nil {
// 				e.logger.Error("occured getting balance of " + request.Account + ": " + err.Error())
// 				return
// 			}
// 			e.logger.Infof("balance of %v is %v on chain %v", account, balance, request.Chain)
// 			// if the balance is equal to the amount, verify with the
// 			// second RPC server.
// 			if balance.Cmp(request.Amount) == 0 || balance.Cmp(request.Amount) == 1 {
// 				contractAddress := common.HexToAddress("0xB1937094fd9e72f05248880E1F4206b241f2BF07")
// 				verifiedBalance, err :=
// 				if err != nil {
// 					e.logger.Error("occured getting balance of " + request.Account + ": " + err.Error() + " from the secondary ETH RPC server")
// 					return
// 				}

// 				if verifiedBalance.Cmp(request.Amount) == 0 || verifiedBalance.Cmp(request.Amount) == 1 {
// 					e.logger.Info("attempting to complete order " + request.TransactionID)
// 					// send a complete order event
// 					awrr := &AccountWatchRequestResult{
// 						AccountWatchRequest: request,
// 						Result:              "suceess",
// 					}

// 					if err := e.Dispatch(awrr); err != nil {
// 						e.logger.Error("error dispatching account watch request result: " + err.Error())
// 					}
// 					canILive = false
// 					return
// 				} else {
// 					e.logger.Error("balance of " + request.Account + " is not equal to " + request.Amount.String())
// 					return
// 				}
// 			}
// 		case <-timer.C:
// 			// if the timer times out, return an error
// 			e := fmt.Sprintf("timeout occured waiting for " + request.Account + " to have a payment of " + request.Amount.String())
// 			e.logger.Info(e)
// 			awrr := &AccountWatchRequestResult{
// 				AccountWatchRequest: request,
// 				Result:              "error",
// 			}

// 			if err := e.Dispatch(awrr); err != nil {
// 				e.logger.Error("error dispatching account watch request result: " + err.Error())
// 			}

// 			canILive = false

// 			return
// 		}
// 	}
// }

func (e *ExchangeServer) sendCoreEVMAsset(fromAddress, privateKey string, toAddress string, amount *big.Int, txid string, rpcClient *ethclient.Client) error {
	// verify there are no missing or
	if toAddress == "" {
		e.logger.Error("toAddress is empty")
		return fmt.Errorf("toAddress is empty")
	}
	if amount == nil {
		e.logger.Error("amount is nil")
		return fmt.Errorf("amount is nil")
	}
	if rpcClient == nil {
		e.logger.Error("rpcClient is nil")
		return fmt.Errorf("rpcClient is nil")
	}
	if txid == "" {
		e.logger.Error("txid is empty")
		return fmt.Errorf("txid is empty")
	}

	// convert the string address to an address
	qualifiedFromAddress := common.HexToAddress(fromAddress)
	// send the currency to the buyer
	// read nonce
	nonce, err := rpcClient.PendingNonceAt(context.Background(), qualifiedFromAddress)
	if err != nil {
		e.logger.Error("cannot get nonce for " + fromAddress + ": " + err.Error())
		return err
	}

	qualifiedToAddress := common.HexToAddress(toAddress)

	// create gas params
	gasLimit := uint64(30000) // in units
	gasPrice, err := rpcClient.SuggestGasPrice(context.Background())
	if err != nil {
		e.logger.Error("error getting gas price: " + err.Error())
		return err
	}

	// create a transaction
	tx := types.NewTransaction(nonce, qualifiedToAddress, amount, gasLimit, gasPrice, nil)

	// fetch chain id
	chainID, err := rpcClient.NetworkID(context.Background())
	if err != nil {
		e.logger.Error("occured getting chain id: " + err.Error())
		return err
	}

	// convert the private key to a private key
	ecdsa, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		e.logger.Error("error converting private key to private key: " + err.Error())
		return err
	}

	// sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), ecdsa)
	if err != nil {
		e.logger.Error("error signing transaction: " + err.Error())
		return err
	}

	// send the transaction
	err = rpcClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		e.logger.Error("error sending transaction: " + err.Error())
		return err
	}

	e.logger.Info("tx sent: " + signedTx.Hash().Hex() + "txid: " + txid)
	return nil
}
