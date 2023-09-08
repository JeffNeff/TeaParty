package be

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

// tradesInProgress is a map containing the trades that are currently in progress on the platform
// the key is the order id
var tradesInProgress = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "trades_in_progress",
	Help: "The total number of trades in progress",
}, []string{"currency", "amount", "trade_asset", "price", "on_chain", "private"})

// Sell is a http route handler that accepts a sell order
// sell orders are stored in an on prem MongoDB database
func (e *ExchangeServer) Sell(w http.ResponseWriter, r *http.Request) {
	// add the CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// parse the request body into a sell order
	sellOrder := &SellOrder{}
	err := json.NewDecoder(r.Body).Decode(&sellOrder)
	if err != nil {
		er := fmt.Sprintf("error decoding sell order: %s", err.Error())
		e.logger.Error(er)
		MetricsFailedSellRequestIncrement(er)
		// respond with an error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid sell order, could not decode body")
		return
	}

	// var found bool
	// if sellOrder.SellerNKNAddress != "" {
	// 	if e.addressesTimedOut != nil {
	// 		// check if the IP address is in the timeout list
	// 		for addr, timeouttime := range e.addressesTimedOut {
	// 			e.logger.Info("checking " + addr + " against " + sellOrder.SellerNKNAddress)
	// 			if addr == sellOrder.SellerNKNAddress {
	// 				found = true
	// 				// check if the timeout time has passed
	// 				if time.Now().After(timeouttime) {
	// 					delete(e.addressesTimedOut, addr)
	// 				} else {
	// 					w.Header().Set("Content-Type", "application/json")
	// 					w.WriteHeader(http.StatusTooManyRequests)
	// 					e.logger.Error("too many requests from " + sellOrder.SellerNKNAddress)
	// 					json.NewEncoder(w).Encode("too many requests")
	// 					return
	// 				}
	// 			}
	// 		}
	// 	}

	// 	// if we are still here and the IP address was not found in the timeout list, add it
	// 	if !found {
	// 		e.logger.Info("adding " + sellOrder.SellerNKNAddress + " to the timeout list")
	// 		e.addressesTimedOut[sellOrder.SellerNKNAddress] = time.Now().Add(time.Minute * 1)
	// 	}
	// }

	var supportedAssets = []string{ETH, GRAMS, POL, KAS, RXD, CEL, SOL, OCT, BSCUSDT, ACC, FLO, CANTO, ETC, BTC, ETHO, CFXE, LTC, MiningGame, "ANY"}

	// if !sellOrder.OnChain {
	// 	adr, err := common.NewMixedcaseAddressFromString(sellOrder.PaymentTransactionID)
	// 	if err != nil {
	// 		er := fmt.Sprintf("error decoding sell order: %s", err.Error())
	// 		e.logger.Error(er)
	// 		MetricsFailedSellRequestIncrement(er)
	// 		// respond with an error
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		json.NewEncoder(w).Encode("invalid sell order")
	// 		return
	// 	}
	// 	verification, err := e.VerifyPaymentWithMoPartyContract(adr.Address())
	// 	if err != nil {
	// 		er := fmt.Sprintf("error verifying sell order: %s", err.Error())
	// 		e.logger.Error(er)
	// 		MetricsFailedSellRequestIncrement(er)
	// 		// respond with an error
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		json.NewEncoder(w).Encode("invalid sell order, could not verify fee payment")
	// 		return
	// 	}

	// 	if !verification {
	// 		er := fmt.Sprintf("error verifying sell order: %s", err.Error())
	// 		e.logger.Error(er)
	// 		MetricsFailedSellRequestIncrement(er)
	// 		// respond with an error
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		json.NewEncoder(w).Encode("invalid sell order, fee is unpaid. Please pay the fee and try again")
	// 		return
	// 	}
	// }

	var supported bool
	for _, asset := range supportedAssets {
		if asset == sellOrder.Currency {
			supported = true
		}
	}
	if !supported {
		er := fmt.Sprintf("unsupported currency: %s", sellOrder.Currency)
		e.logger.Error(er)
		MetricsFailedSellRequestIncrement(er)
		// respond with an error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("unsupported currency")
		return
	}

	for _, asset := range supportedAssets {
		if asset == sellOrder.TradeAsset {
			supported = true
		}
	}

	if !supported {
		er := fmt.Sprintf("unsupported trade asset: %s", sellOrder.TradeAsset)
		e.logger.Error(er)
		MetricsFailedSellRequestIncrement(er)
		// respond with an error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("unsupported trade asset")
		return
	}

	// TODO: verify that sellOrder.Amount != nil

	if sellOrder.Currency == "" || sellOrder.TradeAsset == "" {
		er := fmt.Sprintf("sell order is missing a required field")
		e.logger.Error(er)
		MetricsFailedSellRequestIncrement(er)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid sell order")
		return
	}

	// fetch the current orders from the database
	orders, err := e.fetchOrdersFromDB()
	if err != nil {
		er := fmt.Sprintf("error fetching orders from database: %s", err.Error())
		e.logger.Error(er)
		MetricsFailedSellRequestIncrement(er)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("error fetching orders from database..please try again later ")
		return
	}

	for _, order := range orders {
		if order.TXID == sellOrder.TXID {
			er := fmt.Sprintf("duplicate sell order found")
			e.logger.Error(er)
			MetricsFailedSellRequestIncrement(er)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("duplicate sell order found")
			return
		}
	}

	// check to see if seller has a pre-existing order in the database
	// if they do, return an error
	for _, order := range orders {
		if order.SellerNKNAddress == sellOrder.SellerNKNAddress {
			er := fmt.Sprintf("seller already has an order in the database")
			e.logger.Error(er)
			MetricsFailedSellRequestIncrement(er)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("seller already has an order in the database")
			return
		}
	}

	orders = append(orders, *sellOrder)
	e.updateOrdersInDB(orders)

	if !sellOrder.OnChain {
		// remove the order from the mo-payment contract
		// adr, err := common.NewMixedcaseAddressFromString(sellOrder.PaymentTransactionID)
		// if err != nil {
		// 	er := fmt.Sprintf("error decoding sell order: %s", err.Error())
		// 	e.logger.Error(er)
		// 	MetricsFailedSellRequestIncrement(er)
		// 	// respond with an error
		// 	w.Header().Set("Content-Type", "application/json")
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	json.NewEncoder(w).Encode("invalid sell order")
		// 	return
		// }
		// err2 := e.RemovePaymentFromMoPartyContract(adr.Address())
		// if err2 != nil {
		// 	er := fmt.Sprintf("error removing payment from mo-party contract: %s", err2.Error())
		// 	e.logger.Error(er)
		// 	MetricsFailedSellRequestIncrement(er)
		// 	w.Header().Set("Content-Type", "application/json")
		// 	json.NewEncoder(w).Encode("error removing payment from mo-party contract.. please try again later")
		// }
	}

	// increment the sell request counter in prometheus
	MetricsSellRequestIncrement()
	fmt.Println("sell order added to metrics")

	// respond with a success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	// send the txid in plain text
	w.Write([]byte(sellOrder.TXID))
}

func (e *ExchangeServer) FetchSellOrders(w http.ResponseWriter, r *http.Request) {
	// add the CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// update the state of the orders
	// from redis

	// var err error
	ordr, err := e.fetchOrdersFromDB()
	if err != nil {
		e.logger.Error("error fetching orders from db")
	}

	// filter the private orders
	// from the public orders
	// and return the public orders
	var publicOrders []SellOrder
	for _, order := range ordr {
		fmt.Printf("order: %v", order.Private)
		if !order.Private {
			publicOrders = append(publicOrders, order)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(publicOrders)
}

// Buy is a http route handler that accepts a buy order
func (e *ExchangeServer) Buy(w http.ResponseWriter, r *http.Request) {
	// add the CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// parse the request body into a buy order
	buyOrder := &BuyOrder{}
	err := json.NewDecoder(r.Body).Decode(buyOrder)
	if err != nil {
		er := fmt.Sprintf("error decoding buy order: %s", err.Error())
		e.logger.Error(er)
		MetricsFailedBuyRequestIncrement(er)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid buy order, err decoding")
		return
	}
	// var found bool
	// if buyOrder.BuyerNKNAddress != "" {
	// 	if e.addressesTimedOut != nil {
	// 		for addr, timeouttime := range e.addressesTimedOut {
	// 			e.logger.Info("checking " + addr + " against " + buyOrder.BuyerNKNAddress)
	// 			if addr == buyOrder.BuyerNKNAddress {
	// 				found = true
	// 				// check if the timeout time has passed
	// 				if time.Now().After(timeouttime) {
	// 					delete(e.addressesTimedOut, addr)
	// 				} else {
	// 					w.Header().Set("Content-Type", "application/json")
	// 					w.WriteHeader(http.StatusTooManyRequests)
	// 					e.logger.Error("too many requests from " + buyOrder.BuyerNKNAddress)
	// 					json.NewEncoder(w).Encode("too many requests")
	// 					return
	// 				}
	// 			}
	// 		}
	// 	}

	// 	if !found {
	// 		e.logger.Info("adding " + buyOrder.BuyerNKNAddress + " to the timeout list")
	// 		e.addressesTimedOut[buyOrder.BuyerNKNAddress] = time.Now().Add(time.Minute * 1)
	// 	}
	// }

	// verify no nill values
	if buyOrder.TXID == "" {
		er := fmt.Sprintf("buy order is missing a TXID field")
		e.logger.Error(er)
		MetricsFailedBuyRequestIncrement(er)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("buy order is missing a TXID field")
		return
	}

	// fetch the open orders
	coo, err := e.fetchOrdersFromDB()
	if err != nil {
		er := fmt.Sprintf("error fetching orders from database: %s", err.Error())
		e.logger.Error(er)
		MetricsFailedBuyRequestIncrement(er)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("error fetching orders from database..please try again later ")
		return
	}

	// verify that the order exists in the local memory store
	// if it does not exist return an error
	var sellOrderData SellOrder
	for i, o := range coo {
		if o.TXID == buyOrder.TXID {
			if !o.Locked {
				sellOrderData = o
				// lock the order
				o.Locked = true
				// update the order in the local memory store
				coo[i] = o
				// update the order in the redis store
				e.updateOrdersInDB(coo)
			} else {
				er := fmt.Sprintf("order is already locked")
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("order is locked")
				return
			}
		}
	}

	// if !buyOrder.OnChain {
	// 	adr, err := common.NewMixedcaseAddressFromString(buyOrder.PaymentTransactionID)
	// 	if err != nil {
	// 		er := fmt.Sprintf("error decoding buy order: %s", err.Error())
	// 		e.logger.Error(er)
	// 		MetricsFailedBuyRequestIncrement(er)
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		json.NewEncoder(w).Encode("invalid buy order, could not decode payment transaction id")
	// 		return
	// 	}
	// 	verification, err := e.VerifyPaymentWithMoPartyContract(adr.Address())
	// 	if err != nil {
	// 		er := fmt.Sprintf("error verifying buy order: %s", err.Error())
	// 		e.logger.Error(er)
	// 		MetricsFailedBuyRequestIncrement(er)
	// 		// respond with an error
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		json.NewEncoder(w).Encode("invalid sell order, could not verify fee payment")
	// 		return
	// 	}

	// 	if !verification {
	// 		er := fmt.Sprintf("error verifying buy order: %s", err.Error())
	// 		e.logger.Error(er)
	// 		MetricsFailedBuyRequestIncrement(er)
	// 		// respond with an error
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		json.NewEncoder(w).Encode("invalid sell order, fee is unpaid. Please pay the fee and try again")
	// 		return
	// 	}
	// }

	ta := buyOrder.TradeAsset
	if ta == "" {
		ta = sellOrderData.TradeAsset
	}

	if ta == "" || ta == "ANY" {
		er := fmt.Sprintf("no trade asset found")
		e.logger.Error(er)
		MetricsFailedBuyRequestIncrement(er)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("no trade asset found")
		return
	}

	// TODO:: move a static timeout to a Seller defined parameter
	// passed at order creation time.
	const productionTimeLimit = 7200 // 2 hours
	const devTimelimit = 300         // 300 second
	var timeLimit int64
	if e.dev {
		timeLimit = devTimelimit
	} else {
		timeLimit = productionTimeLimit
	}

	taAmount := sellOrderData.Price

	// Broken ie Experimental
	if sellOrderData.TradeAsset == "ANY" {
		// fetch the market price of the trade asset
		marketPrice, err := FetchMarketPriceInUSD(ta)
		if err != nil {
			er := fmt.Sprintf("error fetching market price: %s", err.Error())
			e.logger.Error(er)
			MetricsFailedBuyRequestIncrement(er)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("error fetching market price")
			return
		}

		e.logger.Infof("market price of %s is: %s", ta, marketPrice)
		// calcuate the amount to send to the buyer

		pgto := big.NewFloat(0).SetInt(sellOrderData.Price)
		bito := big.NewFloat(0).SetInt(marketPrice)

		// convert to big.int

		fl, _ := pgto.Quo(pgto, bito).Float64()

		taAmount = FloatToBigInt(fl)
		// taAmount = big.NewInt(int64(fl * 100000000))
		e.logger.Infof("calculated amount to send to buyer: %s", fl)
	}

	co := &CompletedOrder{
		OrderID:                       buyOrder.TXID,
		BuyerShippingAddress:          buyOrder.BuyerShippingAddress,
		SellerShippingAddress:         sellOrderData.SellerShippingAddress,
		BuyerToFinalizeOnChain:        buyOrder.OnChain,
		SellerToFinalizeOnChain:       sellOrderData.OnChain,
		TradeAsset:                    ta,
		Price:                         taAmount,
		Currency:                      sellOrderData.Currency,
		Amount:                        sellOrderData.Amount,
		Timeout:                       timeLimit,
		SellerNKNAddress:              sellOrderData.SellerNKNAddress,
		BuyerNKNAddress:               buyOrder.BuyerNKNAddress,
		Assisted:                      sellOrderData.Assisted,
		AssistedTradeOrderInformation: &sellOrderData.AssistedTradeOrderInformation,
		NFTID:                         sellOrderData.NFTID,
	}

	buyersAccountWatchRequest := &AccountWatchRequest{}
	// if this is an assisted order, we do not need to create an escrow account
	// as we already have a funded one
	sellersAccountWatchRequest := &AccountWatchRequest{}
	// if this is an assisted order, we do not need to create an escrow account
	// as we already have a funded one
	if !sellOrderData.Assisted {
		switch sellOrderData.Currency {
		case SOL:
			// generate a new solana account for the seller
			acc := e.CreateSolanaAccount()
			co.SellerEscrowWallet = EscrowWallet{
				PublicAddress: acc.PublicKey,
				PrivateKey:    acc.PrivateKey,
				Chain:         SOL,
			}
			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				w.WriteHeader(http.StatusBadRequest)
				e.closeFailedOrder(co)
				return
			}
			sellersAccountWatchRequest = &AccountWatchRequest{
				Account:         co.SellerEscrowWallet.PublicAddress,
				TimeOut:         co.Timeout,
				Chain:           SOL,
				Amount:          co.Amount,
				TransactionID:   co.OrderID,
				Seller:          true,
				FinalizeOnChain: co.SellerToFinalizeOnChain,
			}

		case BTC:
			acc, err := e.generateBTCAccount(co.OrderID, *e.btcNode.rpcConfig, *e.btcNode.rpcClient)
			if err != nil {
				er := fmt.Sprintf("error generating bitcoin account: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("error generating bitcoin account..please try again later")
				e.closeFailedOrder(co)
				return
			}
			co.SellerEscrowWallet = EscrowWallet{
				PublicAddress: acc.PublicKey,
				PrivateKey:    acc.PrivateKey,
				Chain:         BTC,
			}
			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				w.WriteHeader(http.StatusBadRequest)
				e.closeFailedOrder(co)
				return
			}
			sellersAccountWatchRequest = &AccountWatchRequest{
				Account:         co.SellerEscrowWallet.PublicAddress,
				TimeOut:         co.Timeout,
				Chain:           BTC,
				Amount:          co.Amount,
				TransactionID:   co.OrderID,
				Seller:          true,
				FinalizeOnChain: co.SellerToFinalizeOnChain,
			}
		case LTC:
			acc, err := e.generateBTCAccount(co.OrderID, *e.ltcNode.rpcConfig, *e.ltcNode.rpcClient)
			if err != nil {
				er := fmt.Sprintf("error generating litecoin account: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("error generating bitcoin account..please try again later")
				e.closeFailedOrder(co)
				return
			}
			co.SellerEscrowWallet = EscrowWallet{
				PublicAddress: acc.PublicKey,
				PrivateKey:    acc.PrivateKey,
				Chain:         LTC,
			}
			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				w.WriteHeader(http.StatusBadRequest)
				e.closeFailedOrder(co)
				return
			}
			sellersAccountWatchRequest = &AccountWatchRequest{
				Account:         co.SellerEscrowWallet.PublicAddress,
				TimeOut:         co.Timeout,
				Chain:           LTC,
				Amount:          co.Amount,
				TransactionID:   co.OrderID,
				Seller:          true,
				FinalizeOnChain: co.SellerToFinalizeOnChain,
			}

		case RXD:
			acc, err := e.generateBTCAccount(co.OrderID, *e.radiantNode.rpcConfig, *e.radiantNode.rpcClient)
			if err != nil {
				er := fmt.Sprintf("error generating radiant account: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("error generating radiant account..please try again later")
				return
			}

			co.SellerEscrowWallet = EscrowWallet{
				PublicAddress: acc.PublicKey,
				PrivateKey:    acc.PrivateKey,
				Chain:         RXD,
			}

			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				// TODO: this should be handled better...maybe a retry?
				w.WriteHeader(http.StatusBadRequest)
				e.closeFailedOrder(co)
				return
			}

			sellersAccountWatchRequest = &AccountWatchRequest{
				Account:         co.SellerEscrowWallet.PublicAddress,
				TimeOut:         co.Timeout,
				Chain:           RXD,
				Amount:          co.Amount,
				TransactionID:   co.OrderID,
				Seller:          true,
				FinalizeOnChain: co.SellerToFinalizeOnChain,
			}

		case ACC:
			acc := e.generateEVMAccount(ACC)
			co.SellerEscrowWallet = EscrowWallet{
				PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
				PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
				Chain:         ACC,
			}

			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				// TODO: this should be handled better...maybe a retry?
				w.WriteHeader(http.StatusBadRequest)
				e.closeFailedOrder(co)
				return
			}

			sellersAccountWatchRequest = &AccountWatchRequest{
				Account:         co.SellerEscrowWallet.PublicAddress,
				TimeOut:         co.Timeout,
				Chain:           ACC,
				Amount:          co.Amount,
				TransactionID:   co.OrderID,
				Seller:          true,
				FinalizeOnChain: co.SellerToFinalizeOnChain,
			}
		case CFXE:
			acc := e.generateEVMAccount(CFXE)
			co.SellerEscrowWallet = EscrowWallet{
				PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
				PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
				Chain:         CFXE,
			}

			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				// TODO: this should be handled better...maybe a retry?
				w.WriteHeader(http.StatusBadRequest)
				e.closeFailedOrder(co)
				return
			}
		case ETHO:
			acc := e.generateEVMAccount(ETHO)
			co.SellerEscrowWallet = EscrowWallet{
				PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
				PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
				Chain:         ETHO,
			}

			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				// TODO: this should be handled better...maybe a retry?
				w.WriteHeader(http.StatusBadRequest)
				e.closeFailedOrder(co)
				return
			}

		case CANTO:
			acc := e.generateEVMAccount(CANTO)
			co.SellerEscrowWallet = EscrowWallet{
				PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
				PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
				Chain:         CANTO,
			}

			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				// TODO: this should be handled better...maybe a retry?
				w.WriteHeader(http.StatusBadRequest)
				e.closeFailedOrder(co)
				return
			}

			sellersAccountWatchRequest = &AccountWatchRequest{
				Account:         co.SellerEscrowWallet.PublicAddress,
				TimeOut:         co.Timeout,
				Chain:           CANTO,
				Amount:          co.Amount,
				TransactionID:   co.OrderID,
				Seller:          true,
				FinalizeOnChain: co.SellerToFinalizeOnChain,
			}
		case FLO:
			acc := e.generateEVMAccount(FLO)
			co.SellerEscrowWallet = EscrowWallet{
				PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
				PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
				Chain:         FLO,
			}

			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				// TODO: this should be handled better...maybe a retry?
				w.WriteHeader(http.StatusBadRequest)
				e.closeFailedOrder(co)
				return
			}

			sellersAccountWatchRequest = &AccountWatchRequest{
				Account:         co.SellerEscrowWallet.PublicAddress,
				TimeOut:         co.Timeout,
				Chain:           FLO,
				Amount:          co.Amount,
				TransactionID:   co.OrderID,
				Seller:          true,
				FinalizeOnChain: co.SellerToFinalizeOnChain,
			}
		case CEL:
			acc := e.generateEVMAccount(CEL)
			co.SellerEscrowWallet = EscrowWallet{
				PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
				PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
				Chain:         CEL,
			}

			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				// TODO: this should be handled better...maybe a retry?
				w.WriteHeader(http.StatusBadRequest)
				e.closeFailedOrder(co)
				return
			}

			sellersAccountWatchRequest = &AccountWatchRequest{
				Account:         co.SellerEscrowWallet.PublicAddress,
				TimeOut:         co.Timeout,
				Chain:           CEL,
				Amount:          co.Amount,
				TransactionID:   co.OrderID,
				Seller:          true,
				FinalizeOnChain: co.SellerToFinalizeOnChain,
			}

		case ETH:
			acc := e.generateEVMAccount(ETH)
			co.SellerEscrowWallet = EscrowWallet{

				PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
				PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
				Chain:         ETH,
			}

			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				// TODO: this should be handled better...maybe a retry?
				w.WriteHeader(http.StatusBadRequest)
				e.closeFailedOrder(co)
				return
			}

			sellersAccountWatchRequest = &AccountWatchRequest{
				Account:         co.SellerEscrowWallet.PublicAddress,
				TimeOut:         co.Timeout,
				Chain:           ETH,
				Amount:          co.Amount,
				TransactionID:   co.OrderID,
				Seller:          true,
				FinalizeOnChain: co.SellerToFinalizeOnChain,
			}
		case ETC:
			acc := e.generateEVMAccount(ETC)
			co.SellerEscrowWallet = EscrowWallet{

				PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
				PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
				Chain:         ETC,
			}

			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				// TODO: this should be handled better...maybe a retry?
				w.WriteHeader(http.StatusBadRequest)
				e.closeFailedOrder(co)
				return
			}

			sellersAccountWatchRequest = &AccountWatchRequest{
				Account:         co.SellerEscrowWallet.PublicAddress,
				TimeOut:         co.Timeout,
				Chain:           ETC,
				Amount:          co.Amount,
				TransactionID:   co.OrderID,
				Seller:          true,
				FinalizeOnChain: co.SellerToFinalizeOnChain,
			}
		case GRAMS:
			acc := e.generateEVMAccount(GRAMS)
			co.SellerEscrowWallet = EscrowWallet{

				PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
				PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
				Chain:         GRAMS,
			}

			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				e.closeFailedOrder(co)
				return
			}

			sellersAccountWatchRequest = &AccountWatchRequest{
				Account:         co.SellerEscrowWallet.PublicAddress,
				TimeOut:         co.Timeout,
				Chain:           GRAMS,
				Amount:          co.Amount,
				TransactionID:   co.OrderID,
				Seller:          true,
				FinalizeOnChain: co.SellerToFinalizeOnChain,
			}

		case POL:
			acc := e.generateEVMAccount(POL)
			co.SellerEscrowWallet = EscrowWallet{

				PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
				PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
				Chain:         POL,
			}

			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				e.closeFailedOrder(co)
				return
			}

			sellersAccountWatchRequest = &AccountWatchRequest{
				Account:         co.SellerEscrowWallet.PublicAddress,
				TimeOut:         co.Timeout,
				Chain:           POL,
				Amount:          co.Amount,
				TransactionID:   co.OrderID,
				Seller:          true,
				FinalizeOnChain: co.SellerToFinalizeOnChain,
			}

		case OCT:
			acc := e.generateEVMAccount(OCT)
			co.SellerEscrowWallet = EscrowWallet{

				PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
				PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
				Chain:         OCT,
			}

			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				w.WriteHeader(http.StatusBadRequest)
				e.closeFailedOrder(co)
				return
			}

			sellersAccountWatchRequest = &AccountWatchRequest{
				Account:         co.SellerEscrowWallet.PublicAddress,
				TimeOut:         co.Timeout,
				Chain:           OCT,
				Amount:          co.Amount,
				TransactionID:   co.OrderID,
				Seller:          true,
				FinalizeOnChain: co.SellerToFinalizeOnChain,
			}

		// TOKENS
		case BSCUSDT:
			acc := e.generateEVMAccount(BSCUSDT)
			co.SellerEscrowWallet = EscrowWallet{

				PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
				PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
				Chain:         BSCUSDT,
			}

			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				e.closeFailedOrder(co)
				return
			}

			sellersAccountWatchRequest = &AccountWatchRequest{
				Account:         co.SellerEscrowWallet.PublicAddress,
				TimeOut:         co.Timeout,
				Chain:           BSCUSDT,
				Amount:          co.Amount,
				TransactionID:   co.OrderID,
				Seller:          true,
				FinalizeOnChain: co.SellerToFinalizeOnChain,
			}

		case MiningGame:
			acc := e.generateEVMAccount(MiningGame)
			co.SellerEscrowWallet = EscrowWallet{

				PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
				PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
				Chain:         MiningGame,
			}

			if err := e.notifySellerOfBuyer(*co); err != nil {
				er := fmt.Sprintf("error notifying seller of buyer: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("error notifying seller of buyer.. closing the order")
				e.closeFailedOrder(co)
				return
			}

			sellersAccountWatchRequest = &AccountWatchRequest{
				Account:         co.SellerEscrowWallet.PublicAddress,
				TimeOut:         co.Timeout,
				Chain:           MiningGame,
				Amount:          co.Amount,
				TransactionID:   co.OrderID,
				Seller:          true,
				FinalizeOnChain: co.SellerToFinalizeOnChain,
				NFTID:           co.NFTID,
			}

		default:
			er := fmt.Sprintf("error generating seller escrow wallet: %s", err.Error())
			e.logger.Error(er)
			MetricsFailedBuyRequestIncrement(er)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("error generating seller escrow wallet..please try again later")
			e.closeFailedOrder(co)
			return
		}
	}

	switch ta {
	case SOL:
		acc := e.CreateSolanaAccount()
		co.BuyerEscrowWallet = EscrowWallet{
			PublicAddress: acc.PublicKey,
			PrivateKey:    acc.PrivateKey,
			Chain:         SOL,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}

		buyersAccountWatchRequest = &AccountWatchRequest{
			Account:         co.BuyerEscrowWallet.PublicAddress,
			TimeOut:         co.Timeout,
			Chain:           SOL,
			Amount:          co.Price,
			TransactionID:   co.OrderID,
			FinalizeOnChain: co.BuyerToFinalizeOnChain,
		}

	case LTC:
		acc, err := e.generateBTCAccount(co.OrderID, *e.ltcNode.rpcConfig, *e.ltcNode.rpcClient)
		if err != nil {
			er := fmt.Sprintf("error generating litecoin account: %s", err.Error())
			e.logger.Error(er)
			MetricsFailedBuyRequestIncrement(er)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("error generating litecoin account..please try again later")
			return
		}

		co.BuyerEscrowWallet = EscrowWallet{
			PublicAddress: acc.PublicKey,
			PrivateKey:    acc.PrivateKey,
			Chain:         LTC,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}

		buyersAccountWatchRequest = &AccountWatchRequest{
			Account:         co.BuyerEscrowWallet.PublicAddress,
			TimeOut:         co.Timeout,
			Chain:           LTC,
			Amount:          co.Price,
			TransactionID:   co.OrderID,
			FinalizeOnChain: co.BuyerToFinalizeOnChain,
		}
	case BTC:
		acc, err := e.generateBTCAccount(co.OrderID, *e.btcNode.rpcConfig, *e.btcNode.rpcClient)
		if err != nil {
			er := fmt.Sprintf("error generating bitcoin account: %s", err.Error())
			e.logger.Error(er)
			MetricsFailedBuyRequestIncrement(er)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("error generating bitcoin account..please try again later")
			return
		}

		co.BuyerEscrowWallet = EscrowWallet{
			PublicAddress: acc.PublicKey,
			PrivateKey:    acc.PrivateKey,
			Chain:         BTC,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}

		buyersAccountWatchRequest = &AccountWatchRequest{
			Account:         co.BuyerEscrowWallet.PublicAddress,
			TimeOut:         co.Timeout,
			Chain:           BTC,
			Amount:          co.Price,
			TransactionID:   co.OrderID,
			FinalizeOnChain: co.BuyerToFinalizeOnChain,
		}

	case RXD:
		acc, err := e.generateBTCAccount(co.OrderID, *e.radiantNode.rpcConfig, *e.radiantNode.rpcClient)
		if err != nil {
			er := fmt.Sprintf("error generating radiant account: %s", err.Error())
			e.logger.Error(er)
			MetricsFailedBuyRequestIncrement(er)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("error generating radiant account..please try again later")
			return
		}

		co.BuyerEscrowWallet = EscrowWallet{
			PublicAddress: acc.PublicKey,
			PrivateKey:    acc.PrivateKey,
			Chain:         RXD,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}

		buyersAccountWatchRequest = &AccountWatchRequest{
			Account:         co.BuyerEscrowWallet.PublicAddress,
			TimeOut:         co.Timeout,
			Chain:           RXD,
			Amount:          co.Price,
			TransactionID:   co.OrderID,
			FinalizeOnChain: co.BuyerToFinalizeOnChain,
		}

	case CFXE:
		acc := e.generateEVMAccount(CFXE)
		co.BuyerEscrowWallet = EscrowWallet{
			PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
			PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
			Chain:         CFXE,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}
	case ETHO:
		acc := e.generateEVMAccount(ETHO)
		co.BuyerEscrowWallet = EscrowWallet{
			PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
			PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
			Chain:         ETHO,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}
	case ACC:
		acc := e.generateEVMAccount(ACC)
		co.BuyerEscrowWallet = EscrowWallet{

			PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
			PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
			Chain:         ACC,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}
	case CANTO:
		acc := e.generateEVMAccount(CANTO)
		co.BuyerEscrowWallet = EscrowWallet{

			PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
			PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
			Chain:         CANTO,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}
	case FLO:
		acc := e.generateEVMAccount(FLO)
		co.BuyerEscrowWallet = EscrowWallet{

			PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
			PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
			Chain:         FLO,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}

		buyersAccountWatchRequest = &AccountWatchRequest{
			Account:         co.BuyerEscrowWallet.PublicAddress,
			TimeOut:         co.Timeout,
			Chain:           FLO,
			Amount:          co.Price,
			TransactionID:   co.OrderID,
			FinalizeOnChain: co.BuyerToFinalizeOnChain,
		}

	case GRAMS:
		acc := e.generateEVMAccount(GRAMS)
		co.BuyerEscrowWallet = EscrowWallet{

			PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
			PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
			Chain:         GRAMS,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}

		buyersAccountWatchRequest = &AccountWatchRequest{
			Account:         co.BuyerEscrowWallet.PublicAddress,
			TimeOut:         co.Timeout,
			Chain:           GRAMS,
			Amount:          co.Price,
			TransactionID:   co.OrderID,
			FinalizeOnChain: co.BuyerToFinalizeOnChain,
		}

	case ETH:
		acc := e.generateEVMAccount(ETH)
		co.BuyerEscrowWallet = EscrowWallet{

			PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
			PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
			Chain:         ETH,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}

		buyersAccountWatchRequest = &AccountWatchRequest{
			Account:         co.BuyerEscrowWallet.PublicAddress,
			TimeOut:         co.Timeout,
			Chain:           ETH,
			Amount:          co.Price,
			TransactionID:   co.OrderID,
			FinalizeOnChain: co.BuyerToFinalizeOnChain,
		}
	case ETC:
		acc := e.generateEVMAccount(ETC)
		co.BuyerEscrowWallet = EscrowWallet{

			PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
			PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
			Chain:         ETC,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}

		buyersAccountWatchRequest = &AccountWatchRequest{
			Account:         co.BuyerEscrowWallet.PublicAddress,
			TimeOut:         co.Timeout,
			Chain:           ETC,
			Amount:          co.Price,
			TransactionID:   co.OrderID,
			FinalizeOnChain: co.BuyerToFinalizeOnChain,
		}
	case OCT:
		acc := e.generateEVMAccount(OCT)
		co.BuyerEscrowWallet = EscrowWallet{

			PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
			PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
			Chain:         OCT,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}

		buyersAccountWatchRequest = &AccountWatchRequest{
			Account:         co.BuyerEscrowWallet.PublicAddress,
			TimeOut:         co.Timeout,
			Chain:           OCT,
			Amount:          co.Price,
			TransactionID:   co.OrderID,
			FinalizeOnChain: co.BuyerToFinalizeOnChain,
		}

	case CEL:
		acc := e.generateEVMAccount(CEL)
		co.BuyerEscrowWallet = EscrowWallet{

			PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
			PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
			Chain:         CEL,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}

		buyersAccountWatchRequest = &AccountWatchRequest{
			Account:         co.BuyerEscrowWallet.PublicAddress,
			TimeOut:         co.Timeout,
			Chain:           CEL,
			Amount:          co.Price,
			TransactionID:   co.OrderID,
			FinalizeOnChain: co.BuyerToFinalizeOnChain,
		}

	case POL:
		acc := e.generateEVMAccount(POL)
		co.BuyerEscrowWallet = EscrowWallet{

			PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
			PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
			Chain:         POL,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}

		buyersAccountWatchRequest = &AccountWatchRequest{
			Account:         co.BuyerEscrowWallet.PublicAddress,
			TimeOut:         co.Timeout,
			Chain:           POL,
			Amount:          co.Price,
			TransactionID:   co.OrderID,
			FinalizeOnChain: co.BuyerToFinalizeOnChain,
		}

	// TOKENS
	case BSCUSDT:
		acc := e.generateEVMAccount(BSCUSDT)
		co.BuyerEscrowWallet = EscrowWallet{

			PublicAddress: crypto.PubkeyToAddress(acc.PublicKey).String(),
			PrivateKey:    hex.EncodeToString(acc.D.Bytes()),
			Chain:         BSCUSDT,
		}

		// if the buyer elected to finalize off chain we need to send them the buyer pay info via NKN
		if !co.BuyerToFinalizeOnChain {
			if err := e.sendBuyerPayInfo(*co); err != nil {
				er := fmt.Sprintf("sending buyer pay info: %s", err.Error())
				e.logger.Error(er)
				MetricsFailedBuyRequestIncrement(er)
				// respond with an error and close the order
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("sending buyer pay info.. closing the order")
				e.closeFailedOrder(co)
				return
			}
		}

		buyersAccountWatchRequest = &AccountWatchRequest{
			Account:         co.BuyerEscrowWallet.PublicAddress,
			TimeOut:         co.Timeout,
			Chain:           BSCUSDT,
			Amount:          co.Price,
			TransactionID:   co.OrderID,
			FinalizeOnChain: co.BuyerToFinalizeOnChain,
		}
	default:
		e.logger.Error("invalid chain: " + co.TradeAsset)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid chain")
		return
	}
	if sellOrderData.Assisted {
		// We need to mark the sellers payment as complete in the order
		co.SellerPaymentComplete = true
		co.SellerEscrowWallet = sellOrderData.AssistedTradeOrderInformation.SellersEscrowWallet
		co.SellerToFinalizeOnChain = true
		co.SellerRefundAddress = sellOrderData.AssistedTradeOrderInformation.SellerRefundAddress
		co.SellerShippingAddress = sellOrderData.AssistedTradeOrderInformation.SellerShippingAddress
	}

	cco, err := e.fetchCompleteOrdersFromDB()
	if err != nil {
		er := fmt.Sprintf("error fetching complete orders from db: %s", err.Error())
		MetricsFailedBuyRequestIncrement(er)
		e.logger.Error(er)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("error fetching complete orders from db.. please try again later")
		return
	}

	cco = append(cco, *co)
	if err := e.updateCompleteOrdersInDB(cco); err != nil {
		e.logger.Error("error updating complete orders in db: " + err.Error())
		er := fmt.Sprintf("error updating complete orders in db: %s", err.Error())
		MetricsFailedBuyRequestIncrement(er)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("error updating complete orders in db.. please try again later")
	}

	// remove the order from the open orders
	oo, err := e.fetchOrdersFromDB()
	if err != nil {
		er := fmt.Sprintf("error fetching orders from db: %s", err.Error())
		MetricsFailedBuyRequestIncrement(er)
		e.logger.Error(er)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("error fetching orders from db.. please try again later")
		return
	}

	var newoo []SellOrder
	for _, o := range oo {
		if o.TXID != co.OrderID {
			newoo = append(newoo, o)
		}
	}

	if err := e.setOrders(newoo); err != nil {
		e.logger.Error("error updating orders in db: " + err.Error())
		er := fmt.Sprintf("error updating orders in db: %s", err.Error())
		MetricsFailedBuyRequestIncrement(er)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("error updating orders in db.. please try again later")
	}

	// if the buyer is not to finalize on chain we need to remove the fee payment from the mo-party contract
	if !co.BuyerToFinalizeOnChain {
		// adr, err := common.NewMixedcaseAddressFromString(buyOrder.PaymentTransactionID)
		// if err != nil {
		// 	er := fmt.Sprintf("invalid buy order, could not decode payment transaction id: %s", err.Error())
		// 	e.logger.Error(er)
		// 	MetricsFailedBuyRequestIncrement(er)

		// 	w.Header().Set("Content-Type", "application/json")
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	json.NewEncoder(w).Encode("invalid buy order, could not decode payment transaction id")
		// 	return
		// }
		// remove the order from the mo-payment contract
		// err2 := e.RemovePaymentFromMoPartyContract(adr.Address())
		// if err2 != nil {
		// 	er := fmt.Sprintf("removing payment from mo-party contract: %s", err2.Error())
		// 	e.logger.Error(er)
		// 	MetricsFailedBuyRequestIncrement(er)
		// 	// w.Header().Set("Content-Type", "application/json")
		// 	// json.NewEncoder(w).Encode("error removing payment from mo-party contract.. please try again later")
		// }
	} else {
		// we need to respond with the buyer's escrow information to the buyer
		w.Header().Set("Content-Type", "application/json")
		x := map[string]interface{}{
			"publicAddress": co.BuyerEscrowWallet.PublicAddress,
			"chain":         co.BuyerEscrowWallet.Chain,
			"amount":        co.Price,
			"transactionID": co.OrderID,
		}
		e.logger.Info("order bought successfully")
		MetricsBuyRequestIncrement()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(x)
	}

	// add the account watch requests to the database
	buyersAccountWatchRequest.AWRID = uuid.New().String()
	sellersAccountWatchRequest.AWRID = uuid.New().String()
	if err := e.updateAccountWatchRequestInDB(*buyersAccountWatchRequest); err != nil {
		e.logger.Error("error updating account watch request in db: " + err.Error())
	}
	if !sellOrderData.Assisted {
		if err := e.updateAccountWatchRequestInDB(*sellersAccountWatchRequest); err != nil {
			e.logger.Error("error updating account watch request in db: " + err.Error())
		}
	}

	MetricsAddTradeInProgress()
}

// CloseOpenMarketOrder closes an open order
func (e *ExchangeServer) CloseOpenMarketOrder(w http.ResponseWriter, r *http.Request) {
	var req CloseOpenMarketOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e.logger.Error("error decoding request: " + err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("error decoding request")
		return
	}

	// fetch the current orders from the database
	coo, err := e.fetchOrdersFromDB()
	if err != nil {
		e.logger.Error("error fetching orders from database: " + err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("error fetching orders from database..please try again later ")
		return
	}

	for i, o := range coo {
		if o.TXID == req.OrderID {
			// remove the order from the orders slice
			coo = append(coo[:i], coo[i+1:]...)
			e.updateOrdersInDB(coo)
			// if err := e.updateOrdersInDB(coo); err != nil {
			// 	e.logger.Error("error updating order in database: " + err.Error())
			// 	// respond with an error and close the order
			// 	w.Header().Set("Content-Type", "application/json")
			// 	w.WriteHeader(http.StatusBadRequest)
			// 	json.NewEncoder(w).Encode("error updating order in database..please try again later ")
			// 	return
			// }
			e.logger.Info("order closed successfully")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("ok")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("order not found")
}

func FloatToBigInt(val float64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)
	// Set precision if required.
	// bigval.SetPrec(64)

	coin := new(big.Float)
	coin.SetInt(big.NewInt(1000000000000000000))

	bigval.Mul(bigval, coin)

	result := new(big.Int)
	bigval.Int(result) // store converted number in result

	return result
}

// CancleAssistedSell is a handler that allows a user to cancel an assisted sell order
// and be refunded if the order is not locked ( or in progress)
func (e *ExchangeServer) CancleAssistedSell(w http.ResponseWriter, r *http.Request) {
	// add the CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// parse the request body
	var casor CancleAssistedSellOrderRequest
	err := json.NewDecoder(r.Body).Decode(&casor)
	if err != nil {
		e.logger.Error("error decoding request body: " + err.Error())
		// respond with an error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid request")
		return
	}

	// verify the request
	if casor.OrderID == "" {
		e.logger.Error("invalid request")
		// respond with an error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid request")
		return
	}

	if casor.CancelationToken == "" {
		e.logger.Error("invalid request")
		// respond with an error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid request")
		return
	}

	// fetch the current orders from the database
	coo, err := e.fetchOrdersFromDB()
	if err != nil {
		e.logger.Error("error fetching orders from database: " + err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("error fetching orders from database..please try again later ")
		return
	}

	// verify that the order ID and cancelation token match
	var order *SellOrder
	for _, o := range coo {
		if o.TXID == casor.OrderID {
			order = &o
		}
	}

	if order == nil {
		e.logger.Error("invalid request")
		// respond with an error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid request, could not find order")
		return
	}

	// if order.AssistedTradeOrderInformation.CancelToken != casor.CancelationToken {
	// 	e.logger.Error("invalid request")
	// 	// respond with an error
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode("invalid request, cancelation token does not match")
	// 	return
	// }

	// verify that the order is not locked
	if order.Locked {
		e.logger.Error("invalid request")
		// respond with an error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid request, the order is locked")
		return
	}

	// provide the PrivateKey of the escrow account to the user
	// pk := order.AssistedTradeOrderInformation.SellersEscrowWallet.PrivateKey

}

// AssistedSell is a handler that allows a user to place a sell order
// without the use of Tea/NKN. This is done by first verifying the contents
// of the order, creating an associated escrow account, and then
// responding with the address of the escrow account.
// After responding, we will start listening for a payment to the escrow account
// once the payment is received, we will open a new trade order in the marketplace
// and then if a match is found, we will handle the trade for the user. Resulting
// in the user receiving the trade asset via on-chain transfer.
func (e *ExchangeServer) AssistedSell(w http.ResponseWriter, r *http.Request) {
	// add the CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// parse the request body
	var sellOrder SellOrder
	err := json.NewDecoder(r.Body).Decode(&sellOrder)
	if err != nil {
		e.logger.Error("error decoding sell order: " + err.Error())
		// respond with an error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid sell order")
		return
	}

	// verify that the sell order is valid
	if sellOrder.Amount == nil || sellOrder.Currency == "" || sellOrder.TradeAsset == "" {
		e.logger.Error("sell order is missing a required field")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid sell order")
		return
	}

	// verify that the currency and trade asset are supported
	var supportedAssets = []string{ETH, GRAMS, POL, CEL, SOL, OCT, ACC, BSCUSDT}
	supported := false
	for _, currency := range supportedAssets {
		if currency == sellOrder.Currency {
			supported = true
		}
	}

	if !supported {
		e.logger.Error("unsupported currency")
		// respond with an error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("unsupported currency")
		return
	}

	supported = false
	for _, asset := range supportedAssets {
		if asset == sellOrder.TradeAsset {
			supported = true
		}
	}

	if !supported {
		e.logger.Error("unsupported trade asset")
		// respond with an error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("unsupported trade asset")
		return
	}

	// create a random cancelation token
	// cancelToken := uuid.New().String()

	sellOrder.AssistedTradeOrderInformation = AssistedTradeOrderInformation{}
	// sellOrder.AssistedTradeOrderInformation.CancelToken = cancelToken
	// create a new escrow account
	// and respond with the address
	escrowWallet := EscrowWallet{}
	switch sellOrder.Currency {
	case SOL:
		a := e.CreateSolanaAccount()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(a.PublicKey)
		escrowWallet.PublicAddress = a.PublicKey
		escrowWallet.PrivateKey = a.PrivateKey
		escrowWallet.Chain = SOL
	case ACC:
		a := e.generateEVMAccount(ACC)
		if a == nil {
			e.logger.Error("error generating acc account")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(crypto.PubkeyToAddress(a.PublicKey).String())
		escrowWallet.Chain = CEL
		escrowWallet.PublicAddress = crypto.PubkeyToAddress(a.PublicKey).String()
		escrowWallet.PrivateKey = hex.EncodeToString(a.D.Bytes())
	case CEL:
		a := e.generateEVMAccount(CEL)
		if a == nil {
			e.logger.Error("error generating celo account")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(crypto.PubkeyToAddress(a.PublicKey).String())
		escrowWallet.Chain = CEL
		escrowWallet.PublicAddress = crypto.PubkeyToAddress(a.PublicKey).String()
		escrowWallet.PrivateKey = hex.EncodeToString(a.D.Bytes())
	case GRAMS:
		a := e.generateEVMAccount(GRAMS)
		if a == nil {
			e.logger.Error("error generating partychain account")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(crypto.PubkeyToAddress(a.PublicKey).String())
		escrowWallet.Chain = GRAMS
		escrowWallet.PublicAddress = crypto.PubkeyToAddress(a.PublicKey).String()
		escrowWallet.PrivateKey = hex.EncodeToString(a.D.Bytes())
	case ETH:
		a := e.generateEVMAccount(ETH)
		if a == nil {
			e.logger.Error("error generating eth account")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(crypto.PubkeyToAddress(a.PublicKey).String())
		escrowWallet.Chain = ETH
		escrowWallet.PublicAddress = crypto.PubkeyToAddress(a.PublicKey).String()
		escrowWallet.PrivateKey = hex.EncodeToString(a.D.Bytes())
	case ETC:
		a := e.generateEVMAccount(ETC)
		if a == nil {
			e.logger.Error("error generating ETC account")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(crypto.PubkeyToAddress(a.PublicKey).String())
		escrowWallet.Chain = ETC
		escrowWallet.PublicAddress = crypto.PubkeyToAddress(a.PublicKey).String()
		escrowWallet.PrivateKey = hex.EncodeToString(a.D.Bytes())
	case POL:
		a := e.generateEVMAccount(POL)
		if a == nil {
			e.logger.Error("error generating pol account")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(crypto.PubkeyToAddress(a.PublicKey).String())
		escrowWallet.Chain = POL
		escrowWallet.PublicAddress = crypto.PubkeyToAddress(a.PublicKey).String()
		escrowWallet.PrivateKey = hex.EncodeToString(a.D.Bytes())
	case OCT:
		a := e.generateEVMAccount(OCT)
		if a == nil {
			e.logger.Error("error generating oct account")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(crypto.PubkeyToAddress(a.PublicKey).String())
		escrowWallet.Chain = OCT
		escrowWallet.PublicAddress = crypto.PubkeyToAddress(a.PublicKey).String()
		escrowWallet.PrivateKey = hex.EncodeToString(a.D.Bytes())
	case BSCUSDT:
		a := e.generateEVMAccount(BSCUSDT)
		if a == nil {
			e.logger.Error("error generating bscusdt account")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(crypto.PubkeyToAddress(a.PublicKey).String())
		escrowWallet.Chain = BSCUSDT
		escrowWallet.PublicAddress = crypto.PubkeyToAddress(a.PublicKey).String()
		escrowWallet.PrivateKey = hex.EncodeToString(a.D.Bytes())
	}

	// create an account watch request for the escrow account
	const productionTimeLimit = 7200 // 2 hours
	const devTimelimit = 300         // 300 second
	var timeLimit int64
	if e.dev {
		timeLimit = devTimelimit
	} else {
		timeLimit = productionTimeLimit
	}

	if sellOrder.Currency == SOL {
		sellOrder.AssistedTradeOrderInformation.SellersEscrowWallet.PublicAddress = escrowWallet.PublicAddress
	} else {
		sellOrder.AssistedTradeOrderInformation.SellersEscrowWallet.PublicAddress = escrowWallet.PublicAddress
		// sellOrder.AssistedTradeOrderInformation.SellersEscrowWallet.ECDSA = escrowWallet.ECDSA
	}

	sellOrder.AssistedTradeOrderInformation.SellersEscrowWallet.PrivateKey = escrowWallet.PrivateKey
	sellOrder.AssistedTradeOrderInformation.SellersEscrowWallet.Chain = escrowWallet.Chain
	sellOrder.AssistedTradeOrderInformation.SellerRefundAddress = sellOrder.RefundAddress
	sellOrder.AssistedTradeOrderInformation.SellerShippingAddress = sellOrder.SellerShippingAddress

	assistedSellOrderWatchRequest := AccountWatchRequest{
		Account:         escrowWallet.PublicAddress,
		TimeOut:         timeLimit,
		Chain:           sellOrder.Currency,
		Amount:          sellOrder.Amount,
		TransactionID:   "assisted-sell-order",
		FinalizeOnChain: true,
		AssistedSellOrderInformation: AssistedTradeOrderInformation{
			SellersEscrowWallet:   escrowWallet,
			SellerRefundAddress:   sellOrder.RefundAddress,
			SellerShippingAddress: sellOrder.SellerShippingAddress,
			TradeAsset:            sellOrder.TradeAsset,
			Price:                 sellOrder.Price,
			Currency:              sellOrder.Currency,
			Amount:                sellOrder.Amount,
		},
	}

	e.logger.Infof("adding account watch request for assisted sell order: %+v ", sellOrder)
	// go e.watchAccount(assistedSellOrderWatchRequest)
	e.warrenChan <- assistedSellOrderWatchRequest
}
