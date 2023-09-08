package be

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

func (e *ExchangeServer) CreateSolanaAccount() AccountResponse {
	// Create a new account:
	account := solana.NewWallet()
	fmt.Println("account private key:", account.PrivateKey.String())
	fmt.Println("account public key:", account.PublicKey())

	ar := &AccountResponse{
		PrivateKey: account.PrivateKey.String(),
		PublicKey:  account.PublicKey().String(),
	}

	return *ar
}

func (a *ExchangeServer) SendCoreSOLAsset(fromWalletPrivateKey, toAddress, txid string, amount *big.Int) error {
	privateKey, err := solana.PrivateKeyFromBase58(fromWalletPrivateKey)
	if err != nil {
		a.logger.Error("error creating solana public key from string: " + err.Error())
		return err
	}

	toAddressPublicKey, err := solana.PublicKeyFromBase58(toAddress)
	if err != nil {
		a.logger.Error("error creating solana public key from string: " + err.Error())
		return err
	}

	recent, err := a.solNode.rpcClient.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}

	// big.int to lamports
	amountLamparts := amount.Mul(amount, big.NewInt(1000000000))

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				amountLamparts.Uint64(),
				privateKey.PublicKey(),
				toAddressPublicKey,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(privateKey.PublicKey()),
	)
	if err != nil {
		a.logger.Error("error creating solana transaction: " + err.Error())
		return err
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if privateKey.PublicKey().Equals(key) {
				return &privateKey
			}
			return nil
		},
	)
	if err != nil {
		panic(fmt.Errorf("unable to sign transaction: %w", err))
	}

	// TODO: Migrate to ws client so we can use the sendandconfirmtransaction method
	// Send transaction, and wait for confirmation:
	opts := rpc.TransactionOpts{}
	sig, err := a.solNode.rpcClient.SendTransactionWithOpts(
		context.Background(),
		tx,
		opts,
	)

	if err != nil {
		a.logger.Error("error sending solana transaction: " + err.Error())
		return err
	}

	a.logger.Info("sent transaction " + sig.String() + " to " + toAddress + " for " + amount.String() + " SOL")

	return nil

}

func (a *ExchangeServer) waitAndVerifySOLChain(request AccountWatchRequest) {
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

	// the request.Amount is currently in ETH big.Int format convert to uint64
	amount, err := strconv.ParseUint(request.Amount.String(), 10, 64)
	if err != nil {
		a.logger.Errorw("error converting amount to uint64", "error", err)
		return
	}

	// convert from wei to lamports
	amount = amount / 1000000000

	// create a ticker that ticks every 60 seconds
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()
	if a.dev {
		ticker = time.NewTicker(time.Second * 10)
	}

	// create a timer that times out after the specified timeout
	timer := time.NewTimer(time.Second * time.Duration(request.TimeOut))
	defer timer.Stop()

	a.logger.Info("Watching for " + request.Account + " to have a payment of " + fmt.Sprint(amount) + " on chain " + request.Chain)
	// start a for loop that checks the balance of the address
	canILive := true
	for canILive {
		select {
		case <-ticker.C:
			// create new solana public key from string
			pk, err := solana.PublicKeyFromBase58(request.Account)
			if err != nil {
				a.logger.Error("error creating solana public key from string: " + err.Error())
				break
			}

			balance, err := a.solNode.rpcClient.GetBalance(context.Background(), pk, rpc.CommitmentFinalized)
			if err != nil {
				a.logger.Error("occured getting balance of " + request.Account + ": " + err.Error())
				break
			}

			a.logger.Infof("balance of %v is %v on chain %v", request.Account, balance.Value, request.Chain)
			// if the balance is equal to the amount, verify with the
			// second RPC server.
			if balance.Value >= amount {
				verifiedBalance, err := a.solNode.rpcClientTwo.GetBalance(context.Background(), pk, rpc.CommitmentFinalized)
				if err != nil {
					a.logger.Error("occured getting balance of " + request.Account + ": " + err.Error() + " from the secondary SOL RPC server")
					break
				}

				if verifiedBalance.Value >= amount {
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
					break
				}
			}
		case <-timer.C:
			// if the timer times out, return an error
			e := fmt.Sprintf("timeout occured waiting for " + request.Account + " to have a payment of " + request.Amount.String())
			a.logger.Info(e)
			awrr := &AccountWatchRequestResult{
				AccountWatchRequest: request,
				Result:              e,
			}

			if err := a.Dispatch(awrr); err != nil {
				a.logger.Error("error dispatching account watch request result: " + err.Error())
			}

			canILive = false

			return
		}
	}
}
