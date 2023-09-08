package be

import (
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	btcRPC "github.com/btcsuite/btcd/rpcclient"
)

type Transaction struct {
	TxId               string `json:"txid"`
	SourceAddress      string `json:"source_address"`
	DestinationAddress string `json:"destination_address"`
	Amount             int64  `json:"amount"`
	UnsignedTx         string `json:"unsignedtx"`
	SignedTx           string `json:"signedtx"`
}

func (e *ExchangeServer) generateBTCAccount(name string, rpcConfig btcRPC.ConnConfig, rc btcRPC.Client) (*BTCResponse, error) {
	e.logger.Infof("generating BTCish account %v", name)
	_, err := rc.CreateWallet(name)
	if err != nil {
		e.logger.Errorw("create wallet failed", "error", err)
		return nil, err
	}
	// update the rpcConfig to use the new wallet
	rpcConfig.Host = rpcConfig.Host + "/wallet/" + name
	// btcClient, err := btcRPC.New(&rpcConfig, nil)
	if err != nil {
		fmt.Println("err creating client: ", err)
	}

	address, err := rc.GetNewAddress(name)
	if err != nil {
		e.logger.Errorw("generating a new address for" + name + " failed. with error: " + err.Error())
		return nil, err
	}
	// get the private key for the address
	privKey, err := rc.DumpPrivKey(address)
	if err != nil {
		e.logger.Errorw("getting the private key for " + name + " failed. with error: " + err.Error())
		return nil, err
	}
	return &BTCResponse{
		PublicKey:  address.String(),
		PrivateKey: privKey.String(),
	}, nil
}

func (a *ExchangeServer) waitAndVerifyBTCChain(config, config2 btcRPC.ConnConfig, request AccountWatchRequest) {
	a.logger.Infof("waiting for %v to have a payment of %v on chain %v", request.Account, request.Amount, request.Chain)
	if !a.watch {
		a.logger.Info("dev mode is on, not watching for payment. Returning success")
		awrr := &AccountWatchRequestResult{
			AccountWatchRequest: request,
			Result:              "suceess",
		}

		if err := a.Dispatch(awrr); err != nil {
			a.logger.Error("error dispatching account watch request result: " + err.Error())
		}
	}

	config.Host = config.Host + "/wallet/" + request.TransactionID
	client, err := btcRPC.New(&config, nil)
	if err != nil {
		a.logger.Errorw("holy fucknuts batman!! The RPC for BTC-like chain is down!!: " + request.Chain)
		return
	}

	config2.Host = config2.Host + "/wallet/" + request.TransactionID
	client2, err := btcRPC.New(&config2, nil)
	if err != nil {
		a.logger.Errorw("holy fucknuts batman!! The RPC for the BTC-like chain is down!!: " + request.Chain)
		return
	}

	// the request.Amount is currently in ETH big.Int format
	// convert it to BTC
	amount := btcutil.Amount(request.Amount.Int64())
	// create a ticker that ticks every 30 seconds
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()
	if a.dev {
		ticker = time.NewTicker(time.Second * 10)
	}

	// create a timer that times out after the specified timeout
	timer := time.NewTimer(time.Second * time.Duration(request.TimeOut))
	defer timer.Stop()

	a.logger.Info("Watching for " + request.Account + " to have a payment of " + amount.String() + " on chain " + request.Chain)
	// start a for loop that checks the balance of the address
	canILive := true
	for canILive {
		select {
		case <-ticker.C:
			balance, err := client.GetBalance("*")
			if err != nil {
				a.logger.Error("occured getting balance of " + request.Account + ": " + err.Error() + "on " + request.Chain)
			}

			// convert to shatoshi
			balance = balance * 10000000000
			a.logger.Infof("balance of %v is %v on chain %v, looking for %v", request.Account, balance, request.Chain, amount)
			// if the balance is equal to the amount, verify with the
			// second RPC server.
			if balance == amount || balance > amount {
				verifiedBalance, err := client2.GetBalance("*")
				if err != nil {
					a.logger.Error("occured getting balance of " + request.Account + ": " + err.Error() + " from the secondary " + request.Chain + " RPC server")
				}
				// convert to shatoshi
				verifiedBalance = verifiedBalance * 10000000000

				if verifiedBalance == amount || verifiedBalance > amount {
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

// func (e *ExchangeServer) sendCoreBTCAsset(fromWalletPrivateKey string, toAddress string, txid string, amount *big.Int, rpcClient *btcRPC.Client) error {
// 	e.logger.Infof("sending BTC asset %v from %v to %v", amount, fromWalletPrivateKey, toAddress)

// 	// serialize the transaction
// 	var buf bytes.Buffer
// 	if err := tx.Serialize(&buf); err != nil {
// 		e.logger.Errorw("error serializing transaction", "error", err)
// 		return err
// 	}

// 	// send the transaction
// 	txid, err = rpcClient.SendRawTransaction(buf.Bytes(), true)
// 	if err != nil {
// 		e.logger.Errorw("error sending transaction", "error", err)
// 		return err
// 	}

// 	e.logger.Infof("sent BTC asset %v from %v to %v with txid %v", amount, fromWalletPrivateKey, toAddress, txid)

// 	return nil

// }

// func CreateTransaction(secret string, destination string, amount int64, txHash string) (Transaction, error) {
// 	var transaction Transaction
// 	wif, err := btcutil.DecodeWIF(secret)
// 	if err != nil {
// 		return Transaction{}, err
// 	}
// 	addresspubkey, _ := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.MainNetParams)
// 	sourceTx := wire.NewMsgTx(wire.TxVersion)
// 	sourceUtxoHash, _ := chainhash.NewHashFromStr(txHash)
// 	sourceUtxo := wire.NewOutPoint(sourceUtxoHash, 0)
// 	sourceTxIn := wire.NewTxIn(sourceUtxo, nil, nil)
// 	destinationAddress, err := btcutil.DecodeAddress(destination, &chaincfg.MainNetParams)
// 	sourceAddress, err := btcutil.DecodeAddress(addresspubkey.EncodeAddress(), &chaincfg.MainNetParams)
// 	if err != nil {
// 		return Transaction{}, err
// 	}
// 	destinationPkScript, _ := txscript.PayToAddrScript(destinationAddress)
// 	sourcePkScript, _ := txscript.PayToAddrScript(sourceAddress)
// 	sourceTxOut := wire.NewTxOut(amount, sourcePkScript)
// 	sourceTx.AddTxIn(sourceTxIn)
// 	sourceTx.AddTxOut(sourceTxOut)
// 	sourceTxHash := sourceTx.TxHash()
// 	redeemTx := wire.NewMsgTx(wire.TxVersion)
// 	prevOut := wire.NewOutPoint(&sourceTxHash, 0)
// 	redeemTxIn := wire.NewTxIn(prevOut, nil, nil)
// 	redeemTx.AddTxIn(redeemTxIn)
// 	redeemTxOut := wire.NewTxOut(amount, destinationPkScript)
// 	redeemTx.AddTxOut(redeemTxOut)
// 	sigScript, err := txscript.SignatureScript(redeemTx, 0, sourceTx.TxOut[0].PkScript, txscript.SigHashAll, wif.PrivKey, false)
// 	if err != nil {
// 		return Transaction{}, err
// 	}
// 	redeemTx.TxIn[0].SignatureScript = sigScript
// 	flags := txscript.StandardVerifyFlags
// 	vm, err := txscript.NewEngine(sourceTx.TxOut[0].PkScript, redeemTx, 0, flags, nil, nil, amount)
// 	if err != nil {
// 		return Transaction{}, err
// 	}
// 	if err := vm.Execute(); err != nil {
// 		return Transaction{}, err
// 	}
// 	var unsignedTx bytes.Buffer
// 	var signedTx bytes.Buffer
// 	sourceTx.Serialize(&unsignedTx)
// 	redeemTx.Serialize(&signedTx)
// 	transaction.TxId = sourceTxHash.String()
// 	transaction.UnsignedTx = hex.EncodeToString(unsignedTx.Bytes())
// 	transaction.Amount = amount
// 	transaction.SignedTx = hex.EncodeToString(signedTx.Bytes())
// 	transaction.SourceAddress = sourceAddress.EncodeAddress()
// 	transaction.DestinationAddress = destinationAddress.EncodeAddress()
// 	return transaction, nil
// }
