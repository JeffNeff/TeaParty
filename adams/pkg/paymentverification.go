package be

// import (
// 	"context"
// 	"fmt"
// 	"math/big"

// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/common"
// )

// func (e *ExchangeServer) VerifyPaymentWithMoPartyContract(address common.Address) (bool, error) {
// 	e.logger.Infof("Verifying payment for address %s", address)
// 	for _, v := range e.noPartyFeeAddresses {
// 		// fmt.Printf("compare %s to %s ", v, address.String())
// 		if v == address.String() {
// 			e.logger.Infof("Address %s is in the no fee list", address)
// 			e.logger.Info("Returning true")
// 			return true, nil
// 		}
// 	}

// 	// call the smart contract function GetTeaTransactions from a authenticated account
// 	opts := &bind.CallOpts{
// 		Pending: false,
// 		From:    e.moContractTransactOpts.From,
// 		Context: context.Background(),
// 	}

// 	addresses, _, err := e.partyContract.GetTeaPartyTransactions(opts)
// 	if err != nil {
// 		e.logger.Info("Error getting transactions for address %s", address.String())
// 		return false, err
// 	}

// 	fmt.Printf("transactions: %+v", addresses)

// 	for _, v := range addresses {
// 		e.logger.Debugf("comparing %s to %s", v.String(), address.String())
// 		if v == address {
// 			e.logger.Infof("Address %s has a transaction", address.String())
// 			return true, nil
// 		}
// 	}

// 	e.logger.Info("Address %s has no transactions", address.String())
// 	return false, nil
// }

// func (e *ExchangeServer) RemovePaymentFromMoPartyContract(address common.Address) error {
// 	e.logger.Info("Removing payment for address %s", address)
// 	for _, v := range e.noPartyFeeAddresses {
// 		if v == address.String() {
// 			e.logger.Infof("Address %s is in the no fee list. no need to remove", address.String())
// 			e.logger.Info("Returning true")
// 			return nil
// 		}
// 	}
// 	// fetch the proper nonce
// 	nonce, err := e.partyChain.rpcClient.PendingNonceAt(context.Background(), e.moContractTransactOpts.From)
// 	if err != nil {
// 		e.logger.Errorf("Error getting nonce for address %s", e.moContractTransactOpts.From.String())
// 		e.logger.Error(err)
// 		return err
// 	}
// 	e.moContractTransactOpts.Nonce = big.NewInt(int64(nonce))

// 	_, err2 := e.partyContract.RemoveTransaction(e.moContractTransactOpts, address)
// 	if err2 != nil {
// 		e.logger.Errorf("Error removing payment for address %s", address)
// 		e.logger.Error(err2)
// 	}
// 	return nil
// }
