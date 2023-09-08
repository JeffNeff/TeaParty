package be

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (e *ExchangeServer) Dispatch(awrr *AccountWatchRequestResult) error {
	fmt.Printf("Dispatching: %+v", awrr)
	// if the account watch request was for an assisted sell order
	if awrr.AccountWatchRequest.TransactionID == "assisted-sell-order" {
		// if the account watch request result is a success
		if awrr.Result == "suceess" {
			e.logger.Infof("Account watch result: %s", awrr.Result)
			e.logger.Infof("creating a new assisted sell trade order in the marketplace...")
			// create a new assisted sell trade order in the marketplace
			if err := e.createAssistedSellTradeOrderInMarketplace(&awrr.AccountWatchRequest); err != nil {
				e.logger.Errorw("failed to create assisted sell trade order in marketplace", err)
				return err
			}
			// remove
			return nil
		}
		// if the account watch request result is an error
		if awrr.Result == "error" {
			e.logger.Infof("Account watch result: %s for assisted sell order", awrr.Result)
			e.logger.Info("adding the assisted sell order to the failed assisted sell orders db...")
			// add the assisted sell order to the failed assisted sell orders db
			if err := e.addAssistedSellOrderToFailedAssistedSellOrdersDB(awrr); err != nil {
				e.logger.Errorw("failed to add assisted sell order to failed assisted sell orders db", err)
				return err
			}
			return nil
		}
		return nil
	}

	cco, err := e.fetchCompleteOrdersFromDB()
	if err != nil {
		e.logger.Errorw("failed to fetch complete orders from db", err)
		MetricsFailedTradesIncrement("could not fetch complete orders from db")
		return err
	}

	// if the account watch request result is an error
	if awrr.Result == "error" {
		e.logger.Infof("Account watch result: %s", awrr.Result)
		e.logger.Infof("closing the transaction: %s", awrr.AccountWatchRequest.TransactionID)
		// fetch current complete orders from db

		// find the order
		for _, order := range cco {
			// by matching the order id to the transaction id
			if order.OrderID == awrr.AccountWatchRequest.TransactionID {
				// if the buyer funded the escrow wallet
				if order.BuyerPaymentComplete {
					e.logger.Infof("buyer funded the escrow wallet, refunding the buyer %s", order.BuyerNKNAddress)
					// if the buyer elected to be paid off chain, send the funds to the buyer's NKN address first.
					if order.BuyerToFinalizeOnChain {
						// refund the buyer
						if err := e.refundBuyerViaEscrowWallet(order); err != nil {
							e.logger.Errorw("failed to refund buyer via escrow wallet...", err)
							// if the refund fails via NKN, try to refund on-chain.
							if refundErr := e.sendREFUNDToBuyer(order); err != nil {
								e.logger.Errorw("failed to send refund to buyer on chain", err)
								MetricsFailedTradesIncrement("could not refund buyer via escrow wallet on-chain")
								// if the refund fails on-chain, update the order in the db to be handled manually.
								err := e.updateFailedOrdersInDB(order, fmt.Sprintf("failed to refund buyer via escrow wallet: %s", refundErr.Error()))
								if err != nil {
									e.logger.Errorw("failed to update the failed orders in the db", err)
								}
								return refundErr
							}
						}
						// if the buyer elected to be paid on-chain, attempt to send the funds on-chain first.
						// if the on-chain payment fails, refund the buyer via NKN.
					} else {
						if err := e.sendREFUNDToBuyer(order); err != nil {
							e.logger.Errorw("failed to send refund to buyer on chain", err)
							if refundErr := e.refundBuyerViaEscrowWallet(order); err != nil {
								e.logger.Errorw("failed to refund buyer via escrow wallet...", err)
								err := e.updateFailedOrdersInDB(order, fmt.Sprintf("failed to refund buyer via escrow wallet: %s", err.Error()))
								if err != nil {
									e.logger.Errorw("failed to update the failed orders in the db", err)
								}
								MetricsFailedTradesIncrement("could not refund buyer via escrow wallet off-chain")
								return refundErr
							}
						}
					}

					if order.SellerPaymentComplete {
						e.logger.Infof("seller funded the escrow wallet, refunding the seller %s", order.BuyerNKNAddress)
						// if the seller elected to be paid on chain, send the funds on-chain first.
						if order.SellerToFinalizeOnChain {
							if err := e.sendREFUNDToSeller(order); err != nil {
								e.logger.Errorw("failed to send refund to seller on chain", err)
								if refundErr := e.refundSellerViaEscrowWallet(order); err != nil {
									e.logger.Errorw("failed to refund seller via escrow wallet...", err)
									if err := e.updateFailedOrdersInDB(order, fmt.Sprintf("failed to refund seller via escrow wallet: %s", err.Error())); err != nil {
										e.logger.Errorf("failed to update the failed orders in the db: %+v , with error: %+v", order, err)
									}
									MetricsFailedTradesIncrement("could not refund seller via escrow wallet off-chain")
									return refundErr
								}
							}
							// if the seller elected to be paid off chain, refund the seller via NKN first.
						}
					} else {
						if err := e.refundSellerViaEscrowWallet(order); err != nil {
							e.logger.Errorw("failed to refund seller", err)
							if refundErr := e.sendREFUNDToSeller(order); err != nil {
								e.logger.Errorw("failed to send refund to seller on chain", err)
								if err := e.updateFailedOrdersInDB(order, fmt.Sprintf("failed to refund seller via escrow wallet: %s", err.Error())); err != nil {
									e.logger.Errorf("failed to update the failed orders in the db: %+v , with error: %+v", order, err)
								}
								MetricsFailedTradesIncrement("could not refund seller via escrow wallet on-chain")
								return refundErr
							}
						}
					}
					MetricsRemoveTradeInProgress()
					e.closeOrder(&order)
					return nil
				}
			}

		}
	}

	if awrr.Result == "suceess" {
		e.logger.Infof("Successfull Account watch result: %s", awrr.Result)
		if awrr.AccountWatchRequest.Seller {
			for i, order := range cco {
				e.logger.Infof("comparing %s to %s", order.OrderID, awrr.AccountWatchRequest.TransactionID)
				if order.OrderID == awrr.AccountWatchRequest.TransactionID {
					cco[i].SellerPaymentComplete = true
					e.logger.Infof("seller payment complete for order %s", order.OrderID)
					if err := e.updateCompleteOrdersInDBWithIndex(cco, i); err != nil {
						e.logger.Errorw("failed to update the complete orders in the db", err)
						return err
					}
					break
				}
			}
		} else {
			for i, order := range cco {
				if order.OrderID == awrr.AccountWatchRequest.TransactionID {
					e.logger.Infof("buyer payment complete for order %s", order.OrderID)
					cco[i].BuyerPaymentComplete = true
					if err := e.updateCompleteOrdersInDBWithIndex(cco, i); err != nil {
						e.logger.Errorw("failed to update the complete orders in the db", err)
						return err
					}
					break
				}
			}
		}

		// check if the order is complete
		for _, order := range cco {
			if order.OrderID == awrr.AccountWatchRequest.TransactionID {
				if order.BuyerPaymentComplete && order.SellerPaymentComplete {
					e.logger.Infof("order %s is complete", order.OrderID)
					e.logger.Info("Attempting to find a matching order and finalize it")
					var found bool
					for _, order := range cco {
						if order.OrderID == awrr.AccountWatchRequest.TransactionID {
							e.logger.Info("Found a matching order.. calling completeOrder")
							e.completeOrder(order)
							found = true
							return nil
						}
					}

					if !found {
						e.logger.Infof("No matching order found for transaction id: %s", awrr.AccountWatchRequest.TransactionID)
						return nil
					}
				}
			}
		}
	}
	return nil
}

func (e *ExchangeServer) createAssistedSellTradeOrderInMarketplace(awr *AccountWatchRequest) error {
	e.logger.Infof("Creating a new assisted sell trade order in the marketplace for order: %+v", awr)
	// create a new ID for the order
	so := &SellOrder{}
	so.TXID = uuid.New().String()
	so.Currency = awr.Chain
	so.Amount = awr.Amount
	so.Price = awr.AssistedSellOrderInformation.Price
	so.TradeAsset = awr.AssistedSellOrderInformation.TradeAsset
	so.OnChain = true
	so.Assisted = true
	so.AssistedTradeOrderInformation.SellerRefundAddress = awr.AssistedSellOrderInformation.SellerRefundAddress
	so.AssistedTradeOrderInformation.SellerShippingAddress = awr.AssistedSellOrderInformation.SellerShippingAddress
	so.SellerShippingAddress = awr.AssistedSellOrderInformation.SellerShippingAddress
	so.AssistedTradeOrderInformation.SellersEscrowWallet.Chain = awr.AssistedSellOrderInformation.SellersEscrowWallet.Chain
	so.AssistedTradeOrderInformation.SellersEscrowWallet.PublicAddress = awr.AssistedSellOrderInformation.SellersEscrowWallet.PublicAddress
	so.AssistedTradeOrderInformation.SellersEscrowWallet.PrivateKey = awr.AssistedSellOrderInformation.SellersEscrowWallet.PrivateKey
	so.AssistedTradeOrderInformation.TradeAsset = awr.AssistedSellOrderInformation.TradeAsset
	so.AssistedTradeOrderInformation.Currency = awr.AssistedSellOrderInformation.Currency
	so.AssistedTradeOrderInformation.Amount = awr.AssistedSellOrderInformation.Amount
	so.AssistedTradeOrderInformation.Price = awr.AssistedSellOrderInformation.Price
	// fetch current orders from the db
	co, err := e.fetchOrdersFromDB()
	if err != nil {
		e.logger.Errorw("failed to fetch the orders from the db", err)
		return err
	}
	// append the new order to the list
	co = append(co, *so)
	// update the db
	if err := e.updateOrdersInDB(co); err != nil {
		e.logger.Errorw("failed to update the orders in the db", err)
		return err
	}
	return nil
}

// closeOrder closes the order with the specified OrderID
// this is called after a sucessfull transaction
func (e *ExchangeServer) closeOrder(co *CompletedOrder) {
	// fetch current orders from the db
	cco, err := e.fetchCompleteOrdersFromDB()
	if err != nil {
		e.logger.Errorw("failed to fetch the complete orders from the db", err)
		return
	}

	for i, order := range cco {
		if cco[i].OrderID == co.OrderID {
			e.logger.Info("order: " + co.OrderID + " is closing")
			order.Stage = 3
			cco[i] = order
			if err := e.updateCompleteOrdersInDBWithIndex(cco, i); err != nil {
				e.logger.Error("error updating the orders in the db: " + err.Error())
			}
			return
		}
		e.logger.Info("order: " + co.OrderID + " not found")
	}
}

// closeFailedOrder closes the order with the specified OrderID
// this is called after a failed transaction
func (e *ExchangeServer) closeFailedOrder(co *CompletedOrder) {
	// fetch current orders from the db
	cso, err := e.fetchOrdersFromDB()
	if err != nil {
		e.logger.Errorw("failed to fetch the complete orders from the db", err)
		return
	}

	for i, order := range cso {
		if order.TXID == co.OrderID {
			e.logger.Info("order: " + co.OrderID + " is closing")
			cso[i] = order
			err := e.updateOrdersInDB(cso)
			if err != nil {
				e.logger.Error("error updating the orders in the db: " + err.Error())
			}

			// notify the buyer and seller that the order has failed
			go e.notifyBothPartiesOfTradeCancelation(*co)
			return
		}
		e.logger.Errorw("order: " + co.OrderID + " not found")
	}
}

// cancelOrder cancels the order with the specified OrderID by refunding the buyer and seller
// via the escrow wallets
func (e *ExchangeServer) cancelOrder(OrderID string) {
	// fetch current orders from the db
	coo, err := e.fetchOrdersFromDB()
	if err != nil {
		e.logger.Errorw("failed to fetch the complete orders from the db", err)
		return
	}

	cco, err := e.fetchCompleteOrdersFromDB()
	if err != nil {
		e.logger.Errorw("failed to fetch the complete orders from the db", err)
		return
	}

	for i, order := range coo {
		if order.TXID == OrderID {
			if err := e.refundSellerViaEscrowWallet(cco[i]); err != nil {
				e.logger.Error("refund seller via escrow wallet failed, pushing the transaction to state")
			}

			if err := e.refundBuyerViaEscrowWallet(cco[i]); err != nil {
				e.logger.Error("refund seller via escrow wallet failed, pushing the transaction to state")
				// update the failed orders in the db
				if err := e.updateFailedOrdersInDB(cco[i], fmt.Sprintf("refund buyer via escrow wallet failed: %s", err.Error())); err != nil {
					e.logger.Errorf("error updating the failed orders in the db, could not store the failed order: %+v . With error %s", cco[i], err.Error())
				}
			}

			e.logger.Info("order: " + OrderID + " is canceling")
			coo = append(coo[:i], coo[i+1:]...)

			// update the orders in the db
			err := e.updateOrdersInDB(coo)
			if err != nil {
				e.logger.Error("error updating the orders in the db: " + err.Error())
			}

			return
		}
		e.logger.Info("order: " + OrderID + " not found")
	}
	return
}

func (e *ExchangeServer) completeOrder(order CompletedOrder) {
	e.logger.Info("order: " + order.OrderID + " is complete.. now finalizing the order")
	// if the buyer elected to be paid off chain, send the buyer the private key for the escrow wallet
	// via NKN
	if !order.BuyerToFinalizeOnChain {
		e.logger.Info("sending the buyer the private key for the escrow wallet")
		if err := e.sendSellerEscrowWalletPrivateKeyToBuyer(order); err != nil {
			e.logger.Error("error sending the buyer the private key for the escrow wallet: " + err.Error())
			e.logger.Info("attempting to send the buyer funds on chain")
			if err := e.sendFundsToBuyer(order, 100); err != nil {
				e.logger.Error("error sending the funds to the buyer on chain: " + err.Error())
				if err != nil {
					e.logger.Error("attempting to refund the buyer to the provided refund address")
					if err := e.sendREFUNDToBuyer(order); err != nil {
						e.logger.Error("error sending the refund to the buyer on chain: " + err.Error())
						MetricsFailedTradesIncrement("could not send to buyer on or off chain")
						e.logger.Info("attempting to store this failed transaction in the db")
						if err2 := e.updateFailedOrdersInDB(order, fmt.Sprintf("error sending the funds to the buyer on chain: %s", err.Error())); err2 != nil {
							e.logger.Error("error updating the failed orders in the db: " + err2.Error())
							e.logger.Errorf("could not store the failed transaction: %+v", order)

						}
					}
				}
			}
			// if the buyer elected to be paid on chain, send the funds onchain
		}
	} else {
		e.logger.Info("sending the funds to the buyer on chain")
		retry := 0
		backoff := 1
		var fee int64 = 100
		for retry < 5 {
			if err := e.sendFundsToBuyer(order, fee); err != nil {
				e.logger.Error("error sending the funds to the buyer on chain: " + err.Error())
				e.logger.Info("incrementing the retry counter to: " + strconv.Itoa(retry))
				retry++
				if retry == 2 {
					e.logger.Info("removing fee from transaction")
					fee = 0
				}
				e.logger.Info("backing off for " + strconv.Itoa(backoff) + " seconds")
				time.Sleep(time.Duration(backoff) * time.Second)
				backoff *= 2
			} else {
				e.logger.Info("successfully sent the funds to the buyer on chain")
				break
			}

			if retry == 5 {
				e.logger.Error("sending the funds to the buyer on chain failed after 5 retries")
				e.logger.Info("attempting to refund the buyer to the provided refund address")
				if err := e.sendREFUNDToBuyer(order); err != nil {
					e.logger.Error("error sending the refund to the buyer on chain: " + err.Error())
					e.logger.Info("attempting to store this failed transaction in the db")
					if err2 := e.updateFailedOrdersInDB(order, fmt.Sprintf("error sending the funds to the buyer on chain: %s", err.Error())); err2 != nil {
						e.logger.Error("error updating the failed orders in the db: " + err2.Error())
						e.logger.Errorf("could not store the failed transaction: %+v", order)
						e.logger.Info("attempting to write the failed transaction to the failed transactions file")
						if err3 := e.writeFailedTransactionToFile(order); err3 != nil {
							e.logger.Error("error writing the failed transaction to the failed transactions file: " + err3.Error())
						}
					}
				}
			}
		}
	}

	// if the seller elected to be paid off chain, send the seller the private key for the escrow wallet
	// via NKN
	if !order.SellerToFinalizeOnChain {
		e.logger.Info("sending the seller the private key for the escrow wallet")
		if err2 := e.sendBuyerEscrowWalletPrivateKeyToSeller(order); err2 != nil {
			e.logger.Error("error sending the seller the private key for the escrow wallet: " + err2.Error())
			e.logger.Info("attempting to send the seller funds on chain")
			if err := e.sendFundsToSeller(order, int64(100)); err != nil {
				e.logger.Error("error sending the funds to the seller on chain: " + err.Error())
				e.logger.Info("attempting to refund the seller to the provided refund address")
				if err := e.sendREFUNDToSeller(order); err != nil {
					e.logger.Error("error sending the refund to the seller on chain: " + err.Error())
					e.logger.Info("attempting to store this failed transaction in the db")
					if err2 := e.updateFailedOrdersInDB(order, fmt.Sprintf("error sending the funds to the seller on chain: %s", err.Error())); err2 != nil {
						e.logger.Error("error updating the failed orders in the db: " + err2.Error())
						e.logger.Errorf("could not store the failed transaction: %+v", order)
					}
				}
			}
		}
		// if the seller elected to be paid on chain, send the funds onchain
	} else {
		e.logger.Info("sending the funds to the seller on chain")
		// create a retry and backoff mechanism for sending the funds to the seller
		retry := 0
		backoff := 1
		var fee int64 = 100
		// if the funds are not sent successfully, retry 5 times
		for retry < 5 {
			if err := e.sendFundsToSeller(order, fee); err != nil {
				e.logger.Error("error sending the funds to the seller on chain: " + err.Error())
				e.logger.Info("incrementing the retry counter to: " + strconv.Itoa(retry))
				retry++
				if retry == 2 {
					e.logger.Info("removing fee from transaction")
					fee = 0
				}
				e.logger.Info("backing off for " + strconv.Itoa(backoff) + " seconds")
				time.Sleep(time.Duration(backoff) * time.Second)
				backoff *= 2
			} else {
				e.logger.Info("successfully sent the funds to the seller on chain")
				break
			}
			if retry == 5 {
				e.logger.Error("error sending the funds to the seller on chain: Max retries reached")
				e.logger.Info("attempting to refund the seller to the provided refund address")
				if err := e.sendREFUNDToSeller(order); err != nil {
					e.logger.Error("error sending the refund to the seller on chain: " + err.Error())
					e.logger.Info("attempting to store this failed transaction in the db")
					if err2 := e.updateFailedOrdersInDB(order, fmt.Sprintf("error sending the funds to the seller on chain: %s", err.Error())); err2 != nil {
						e.logger.Error("error updating the failed orders in the db: " + err2.Error())
						e.logger.Errorf("could not store the failed transaction: %+v", order)
					}
				}
			}
		}
	}

	MetricsCompletedTradesIncrement()
	e.closeOrder(&order)
}

func (e *ExchangeServer) sendFundsToBuyer(order CompletedOrder, fee int64) error {
	e.logger.Info("sending funds to buyer for order : %+v", order)
	// currenty we do not support sending `kaspa` on chain
	if order.Currency == KAS || order.Currency == BTC || order.Currency == RXD {
		return fmt.Errorf("sending kaspa on chain is not supported")
	}

	var curencyFee *big.Int
	var assetFee *big.Int

	if fee != 0 {
		curencyFee = new(big.Int).Div(order.Amount, big.NewInt(fee))
		order.Amount.Sub(order.Amount, curencyFee)
		assetFee = new(big.Int).Div(order.Price, big.NewInt(fee))
		order.Price.Sub(order.Price, assetFee)
	}

	// send the funds to the buyer
	switch order.Currency {
	case SOL:
		if err := e.SendCoreSOLAsset(order.SellerEscrowWallet.PrivateKey, order.BuyerShippingAddress, order.OrderID, order.Amount); err != nil {
			e.logger.Error("sending funds to buyer on SOL: " + err.Error())
			return err
		}
	case ACC:
		if err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.BuyerShippingAddress, order.Amount, order.OrderID, e.altcoinchain.rpcClient); err != nil {
			e.logger.Error("sending funds to buyer on ACC: " + err.Error())
			return err
		}
	case CEL:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.BuyerShippingAddress, order.Amount, order.OrderID, e.celoNode.rpcClient)
		if err != nil {
			e.logger.Error("sending funds to buyer on CEL: " + err.Error())
			return err
		}
	case FLO:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.BuyerShippingAddress, order.Amount, order.OrderID, e.floNode.rpcClient)
		if err != nil {
			e.logger.Error("sending funds to buyer on FLO: " + err.Error())
			return err
		}

	case GRAMS:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.BuyerShippingAddress, order.Amount, order.OrderID, e.partyChain.rpcClient)
		if err != nil {
			e.logger.Error("sending funds to buyer on GRAMS: " + err.Error())
			return err
		}

	case POL:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.BuyerShippingAddress, order.Amount, order.OrderID, e.polygonNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to buyer on POL: " + err.Error())
			return err
		}

	case ETH:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.BuyerShippingAddress, order.Amount, order.OrderID, e.ethNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to buyer on ETH: " + err.Error())
			return err
		}
	case ETC:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.BuyerShippingAddress, order.Amount, order.OrderID, e.etcNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to buyer on ETC: " + err.Error())
			return err
		}
	case OCT:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.BuyerShippingAddress, order.Amount, order.OrderID, e.octNode.rpcClient)
		if err != nil {
			return err
		}
	case CANTO:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.BuyerShippingAddress, order.Amount, order.OrderID, e.cantoNode.rpcClient)
		if err != nil {
			return err
		}
	case ETHO:
		if err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.BuyerShippingAddress, order.Amount, order.OrderID, e.ethONode.rpcClient); err != nil {
			return err
		}
	case CFXE:
		if err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.BuyerShippingAddress, order.Amount, order.OrderID, e.cfxEspaceNode.rpcClient); err != nil {
			return err
		}
	}
	return nil
}

func (e *ExchangeServer) sendREFUNDToBuyer(order CompletedOrder) error {
	// currenty we do not support sending `kaspa` on chain
	if order.Currency == KAS || order.Currency == BTC || order.Currency == RXD {
		return fmt.Errorf("sending kaspa on chain is not supported")
	}

	curencyFee := new(big.Int).Div(order.Amount, big.NewInt(100))
	order.Amount.Sub(order.Amount, curencyFee)
	assetFee := new(big.Int).Div(order.Price, big.NewInt(100))
	order.Price.Sub(order.Price, assetFee)

	// send the funds to the buyer
	switch order.Currency {
	case SOL:
		if err := e.SendCoreSOLAsset(order.BuyerEscrowWallet.PrivateKey, order.BuyerRefundAddress, order.OrderID, order.Amount); err != nil {
			e.logger.Error("sending funds to buyer on SOL: " + err.Error())
			return err
		}
	case ACC:
		if err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.BuyerRefundAddress, order.Amount, order.OrderID, e.altcoinchain.rpcClient); err != nil {
			e.logger.Error("sending funds to buyer on ACC: " + err.Error())
			return err
		}
	case CEL:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.BuyerRefundAddress, order.Price, order.OrderID, e.celoNode.rpcClient)
		if err != nil {
			e.logger.Error("sending funds to buyer on CEL: " + err.Error())
			return err
		}
	case FLO:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.BuyerRefundAddress, order.Price, order.OrderID, e.floNode.rpcClient)
		if err != nil {
			e.logger.Error("sending funds to buyer on FLO: " + err.Error())
			return err
		}
	case GRAMS:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.BuyerRefundAddress, order.Price, order.OrderID, e.partyChain.rpcClient)
		if err != nil {
			e.logger.Error("sending funds to buyer on GRAMS: " + err.Error())
			return err
		}

	case POL:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.BuyerRefundAddress, order.Price, order.OrderID, e.polygonNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to buyer on POL: " + err.Error())
			return err
		}

	case ETH:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.BuyerRefundAddress, order.Price, order.OrderID, e.ethNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to buyer on ETH: " + err.Error())
			return err
		}
	case ETC:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.BuyerRefundAddress, order.Price, order.OrderID, e.etcNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to buyer on ETC: " + err.Error())
			return err
		}
	case OCT:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.BuyerRefundAddress, order.Price, order.OrderID, e.octNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to buyer on OCT: " + err.Error())
			return err
		}
	case CANTO:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.BuyerRefundAddress, order.Price, order.OrderID, e.cantoNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to buyer on CANTO: " + err.Error())
			return err
		}
	case ETHO:
		if err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.BuyerRefundAddress, order.Price, order.OrderID, e.ethONode.rpcClient); err != nil {
			e.logger.Error("error sending funds to buyer on ETHO: " + err.Error())
			return err
		}
	case CFXE:
		if err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.BuyerRefundAddress, order.Price, order.OrderID, e.cfxEspaceNode.rpcClient); err != nil {
			e.logger.Error("error sending funds to buyer on CFXE: " + err.Error())
			return err
		}
	}
	return nil
}

// sendFundsToSeller provides functionality to send funds to the seller
func (e *ExchangeServer) sendFundsToSeller(order CompletedOrder, fee int64) error {
	if order.TradeAsset == KAS || order.TradeAsset == BTC || order.TradeAsset == RXD {
		return fmt.Errorf("sending kaspa on chain is not supported")
	}

	var curencyFee *big.Int
	var assetFee *big.Int

	if fee != 0 {
		curencyFee = new(big.Int).Div(order.Amount, big.NewInt(fee))
		order.Amount.Sub(order.Amount, curencyFee)
		assetFee = new(big.Int).Div(order.Price, big.NewInt(fee))
		order.Price.Sub(order.Price, assetFee)
	}

	switch order.TradeAsset {
	case SOL:
		if err := e.SendCoreSOLAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.SellerShippingAddress, order.Price); err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case ACC:
		if err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.SellerShippingAddress, order.Amount, order.OrderID, e.altcoinchain.rpcClient); err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case CEL:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.SellerShippingAddress, order.Price, order.OrderID, e.celoNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}

	case GRAMS:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.SellerShippingAddress, order.Price, order.OrderID, e.partyChain.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}

	case POL:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.SellerShippingAddress, order.Price, order.OrderID, e.polygonNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}

	case ETH:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.SellerShippingAddress, order.Price, order.OrderID, e.ethNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case ETC:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.SellerShippingAddress, order.Price, order.OrderID, e.etcNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case OCT:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.SellerShippingAddress, order.Price, order.OrderID, e.octNode.rpcClient)
		if err != nil {
			// wait for 10 seconds and try again
			time.Sleep(10 * time.Second)
			err = e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.SellerShippingAddress, order.Price, order.OrderID, e.octNode.rpcClient)
			if err != nil {
				e.logger.Error("error sending funds to seller: " + err.Error())
				return err
			}
		}

	case FLO:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.SellerShippingAddress, order.Price, order.OrderID, e.floNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case CANTO:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.SellerShippingAddress, order.Price, order.OrderID, e.cantoNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case ETHO:
		err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.SellerShippingAddress, order.Price, order.OrderID, e.ethONode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case CFXE:
		if err := e.sendCoreEVMAsset(order.BuyerEscrowWallet.PublicAddress, order.BuyerEscrowWallet.PrivateKey, order.SellerShippingAddress, order.Price, order.OrderID, e.cfxEspaceNode.rpcClient); err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	}
	return nil
}

// sendREFUNDToSeller provides functionality to send funds to the seller
func (e *ExchangeServer) sendREFUNDToSeller(order CompletedOrder) error {
	if order.TradeAsset == KAS || order.TradeAsset == BTC || order.TradeAsset == RXD {
		return fmt.Errorf("sending kaspa on chain is not supported")
	}

	curencyFee := new(big.Int).Div(order.Amount, big.NewInt(100))
	order.Amount.Sub(order.Amount, curencyFee)
	assetFee := new(big.Int).Div(order.Price, big.NewInt(100))
	order.Price.Sub(order.Price, assetFee)

	switch order.TradeAsset {
	case SOL:
		if err := e.SendCoreSOLAsset(order.SellerEscrowWallet.PrivateKey, order.SellerRefundAddress, order.OrderID, order.Amount); err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}

	case CEL:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.SellerRefundAddress, order.Amount, order.OrderID, e.celoNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case FLO:
		if err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.SellerRefundAddress, order.Amount, order.OrderID, e.floNode.rpcClient); err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case ACC:
		if err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.SellerRefundAddress, order.Amount, order.OrderID, e.altcoinchain.rpcClient); err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case GRAMS:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.SellerRefundAddress, order.Amount, order.OrderID, e.partyChain.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}

	case POL:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.SellerRefundAddress, order.Amount, order.OrderID, e.polygonNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case ETH:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.SellerRefundAddress, order.Amount, order.OrderID, e.ethNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case ETC:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.SellerRefundAddress, order.Amount, order.OrderID, e.etcNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case OCT:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.SellerRefundAddress, order.Amount, order.OrderID, e.octNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case CANTO:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.SellerRefundAddress, order.Amount, order.OrderID, e.cantoNode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case ETHO:
		err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.SellerRefundAddress, order.Amount, order.OrderID, e.ethONode.rpcClient)
		if err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	case CFXE:
		if err := e.sendCoreEVMAsset(order.SellerEscrowWallet.PublicAddress, order.SellerEscrowWallet.PrivateKey, order.SellerRefundAddress, order.Amount, order.OrderID, e.cfxEspaceNode.rpcClient); err != nil {
			e.logger.Error("error sending funds to seller: " + err.Error())
			return err
		}
	}
	return nil
}

func (e *ExchangeServer) watchAccount(awr *AccountWatchRequest) {
	e.logger.Info("watching account: " + awr.Account)
	awr.Locked = true
	awr.LockedBy = os.Getenv("POD_NAME")
	awr.LockedTime = time.Now()
	// tell the database that this instance of the exchange is watching this account
	if err := e.updateAccountWatchRequestInDB(*awr); err != nil {
		e.logger.Error("error updating account watch request in db: " + err.Error())
	}

	switch awr.Chain {
	case SOL:
		e.waitAndVerifySOLChain(*awr)
	case ACC:
		e.waitAndVerifyEVMChain(context.Background(), e.altcoinchain.rpcClient, e.altcoinchain.rpcClientTwo, *awr)
	case CANTO:
		e.waitAndVerifyEVMChain(context.Background(), e.cantoNode.rpcClient, e.cantoNode.rpcClientTwo, *awr)
	case FLO:
		e.waitAndVerifyEVMChain(context.Background(), e.floNode.rpcClient, e.floNode.rpcClientTwo, *awr)
	case CEL:
		e.waitAndVerifyEVMChain(context.Background(), e.celoNode.rpcClient, e.celoNode.rpcClientTwo, *awr)
	case BTC:
		e.waitAndVerifyBTCChain(*e.btcNode.rpcConfig, *e.btcNode.rpcConfigTwo, *awr)
	case LTC:
		e.waitAndVerifyBTCChain(*e.ltcNode.rpcConfig, *e.ltcNode.rpcConfigTwo, *awr)
	case RXD:
		e.waitAndVerifyBTCChain(*e.radiantNode.rpcConfig, *e.radiantNode.rpcConfigTwo, *awr)
	case ETH:
		e.waitAndVerifyEVMChain(context.Background(), e.ethNode.rpcClient, e.ethNode.rpcClientTwo, *awr)
	case ETC:
		e.waitAndVerifyEVMChain(context.Background(), e.etcNode.rpcClient, e.etcNode.rpcClientTwo, *awr)
	case ETHONE:
		e.waitAndVerifyEVMChain(context.Background(), e.ethOneNode.rpcClient, e.ethOneNode.rpcClientTwo, *awr)
	case GRAMS:
		e.waitAndVerifyEVMChain(context.Background(), e.partyChain.rpcClient, e.partyChain.rpcClientTwo, *awr)
	case POL:
		e.waitAndVerifyEVMChain(context.Background(), e.polygonNode.rpcClient, e.polygonNode.rpcClientTwo, *awr)
	case OCT:
		e.waitAndVerifyEVMChain(context.Background(), e.octNode.rpcClient, e.octNode.rpcClientTwo, *awr)
	case KAS:
		e.waitAndVerifyKASChain(context.Background(), e.kaspaNode.rpcClient, e.kaspaNode.rpcClientTwo, *awr)
	case ETHO:
		e.waitAndVerifyEVMChain(context.Background(), e.ethONode.rpcClient, e.ethONode.rpcClientTwo, *awr)
	case CFXE:
		e.waitAndVerifyEVMChain(context.Background(), e.cfxEspaceNode.rpcClient, e.cfxEspaceNode.rpcClientTwo, *awr)
	case MiningGame:
		e.waitAndVerifyThatNFTIsAvalibleOnEVMChain(*awr)
	// TOKENS
	case BSCUSDT:
		e.waitAndVerifyBSCUSDT(*awr)
	default:
		e.logger.Error("unknown chain: " + awr.Chain)
	}
}
