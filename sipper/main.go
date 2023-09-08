package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"

	uuid "github.com/google/uuid"

	redis "github.com/go-redis/redis/v9"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	ETH = "ethereum"
	MO  = "mineonlium"
	POL = "polygon"
	KAS = "kaspa"
	ALP = "alephium"
	BTC = "bitcoin"
)

func main() {
	// accept two payment transaction ID's for the MO burn.

	// accept a server address as an enviorment variable
	// for the redis server (teabarrel)
	teaBarrel := os.Getenv("REDIS_ADDR")
	if teaBarrel == "" {
		teaBarrel = "192.168.50.7:6379"
	}

	// accept a server address as an enviorment variable
	// for the `tea` server
	teaServer := os.Getenv("TEA_ADDR")
	if teaServer == "" {
		teaServer = "http://0.0.0.0:8081"
	}

	// // open a browser to the tea server
	// if err := openbrowser(teaServer); err != nil {
	// 	(err)
	// }
	teaWS := os.Getenv("TEA_WS")
	if teaWS == "" {
		teaWS = "ws://0.0.0.0:8081/ws"
	}

	adamsServer := os.Getenv("ADAMS_ADDR")
	if adamsServer == "" {
		// adamsServer = "http://10.1.243.121:8080"
		// adamsServer = "http://0.0.0.0:8080"
		adamsServer = "http://192.168.50.5:8080"
	}

	// create new redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: teaBarrel,
		DB:   0,
	})

	// create a new eth private key
	ethPK := generateEthAccount()

	// create a new mo private key
	moPK := generateEthAccount()

	sipper := Sipper{
		teaServer:   teaServer,
		teaBarrel:   teaBarrel,
		adamsServer: adamsServer,
		moPK:        moPK,
		ethPK:       ethPK,
		redisClient: redisClient,
	}

	// create a websocket connection with tea backend to simulate a user
	c, _, err := websocket.DefaultDialer.Dial(teaWS, nil)
	if err != nil {

	}
	defer c.Close()

	// go sipper.TeaWSConnectionHandler()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Err:", err)
				return
			}
			log.Println("")
			log.Printf("recv: %s", message)
			log.Println("")
		}
	}()

	// sleep for 1 second to allow the websocket connection to be established
	time.Sleep(1 * time.Second)

	// create and start a bunch of transactions and let them fail.
	// while checking the redis server to make sure the data is being
	// stored correctly
	// sipper.CreateTestTrade("mineonlium", "mineonlium")
	// create 100 trades with 10 goroutines

	// wg := sync.WaitGroup{}
	// startTime := time.Now()
	// for i := 0; i < 100000; i++ {
	// 	wg.Add(1)
	// 	sipper.CreateTestTrade("grams", "grams")
	// 	wg.Done()
	// }
	// endTime := time.Now()

	// fmt.Println("time to create 100000 trades: ", endTime.Sub(startTime))

	go func() {
		for {
			// if err := sipper.BuyAllTheThings(); err != nil {
			// 	fmt.Println(err)
			// }
			go sipper.BuyAllTheThings()
		}
	}()
	// sipper.CreateTestTrade("celo", "celo")
	// sipper.CreateTestTrade("ethereum", "ethereum")
	// sipper.CreateTestTrade("octa", "octa")
	// sipper.CreateTestTrade("bscUSDT", "bscUSDT")
	// sipper.CreateTestTrade("solana", "solana")
	// sipper.CreateTestTrade("ethOne", "ethOne")

	// // start a trade for all the things and verify that escrow accounts are created
	if err := sipper.BuyAllTheThings(); err != nil {
		fmt.Println(err)
	}

	// create a bunch of transactions quickly but do not test redis

	// create transactions and fail on reciving the order via NKN from both the buyer and seller
	// if err := sipper.CreateTradeBuyAndFailOnBuyerNKNRecieve("ethereum", "ethereum"); err != nil {
	// 	(err)
	// }

	// create transactions and fail on payment to the escrow wallet from the buyer

	// create transactions and fail on payment to the escrow wallet from the seller

	// create transactions and succeed on payment to the escrow wallet from the buyer

	// create transactions and succeed on payment to the escrow wallet from the seller

	waitforever := make(chan struct{})
	<-waitforever

}

func (s *Sipper) CreateTradeBuyAndFailOnBuyerNKNRecieve(tradeAsset, currency string) error {
	// create a new trade
	txid := s.CreateTestTrade(tradeAsset, currency)
	so, _ := s.ListSellOrdersViaTea()
	fmt.Println("sell orders: ", so)
	// buy the trade
	if err := s.BuySpecificTrade(txid); err != nil {
		return err
	}
	return nil
}

func (s *Sipper) BuySpecificTrade(tradeID string) error {
	fmt.Println("buying trade: ", tradeID)
	// create a new trade
	buyOrder := BuyOrder{
		TXID:                 tradeID,
		BuyerShippingAddress: "0x9cA67FFE69698d963A393E9338aD3BcfD2CEa02e",
		BuyerNKNAddress:      "nknAddr",
		PaymentTransactionID: "0x9cA67FFE69698d963A393E9338aD3BcfD2CEa02e",
		RefundAddress:        "0x9cA67FFE69698d963A393E9338aD3BcfD2CEa02e",
	}
	fmt.Println("buy order: ", buyOrder)

	// initiate the transaction via tea
	err := s.CreateBuyTransactionViaTea(buyOrder)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// BuyAllTheThings is a test function that will initate a trade for all open market orders.
func (s *Sipper) BuyAllTheThings() error {
	buyersRefundAddress := "0x5dd4039c32F6EEF427D6F67600D8920c9631D59D"
	buyersPublicAddress := "0x5dd4039c32F6EEF427D6F67600D8920c9631D59D"
	// retireve all the open market orders
	order, err := s.redisClient.Get(context.Background(), "orders").Result()
	if err != nil {
		return err
	}

	//unmarshal the order
	var o []SellOrder
	err = json.Unmarshal([]byte(order), &o)
	if err != nil {
		return err
	}

	// loop through the orders and create a trade for each one
	for _, order := range o {
		// create a trade order for the order
		// // create a new buy transaction via tea
		nknAddr, err := s.RetrieveNKNAddress()
		if err != nil {
			return err
		}

		buyOrder := BuyOrder{
			TXID:                 order.TXID,
			BuyerShippingAddress: buyersPublicAddress,
			BuyerNKNAddress:      nknAddr,
			PaymentTransactionID: buyersPublicAddress,
			RefundAddress:        buyersRefundAddress,
		}
		fmt.Println("buy order: ", buyOrder)

		// initiate the transaction via tea
		err = s.CreateBuyTransactionViaTea(buyOrder)
		if err != nil {
			fmt.Println("error creating buy transaction via tea: ", err)
			return err
		}
	}
	return nil
}

// EthToMo is a test function to test the creation and execution of a trade order for ETH and MO
func (s *Sipper) CreateTestTrade(tradeAsset, currency string) string {
	// retrieve the NKN address from tea
	nknAddr, err := s.RetrieveNKNAddress()
	if err != nil {
		fmt.Println("error retrieving nkn address: ", err)
	}

	// sellersPublicAddress := crypto.PubkeyToAddress(s.moPK.PublicKey).Hex()
	// buyersPublicAddress := crypto.PubkeyToAddress(s.ethPK.PublicKey).Hex()
	// sellersRefundAddress := crypto.PubkeyToAddress(s.ethPK.PublicKey).Hex()
	// buyersRefundAddress := crypto.PubkeyToAddress(s.moPK.PublicKey).Hex()

	sellersPublicAddress := "0x5dd4039c32F6EEF427D6F67600D8920c9631D59D"
	sellersRefundAddress := "0x5dd4039c32F6EEF427D6F67600D8920c9631D59D"
	// buyersRefundAddress := "0x9cA67FFE69698d963A393E9338aD3BcfD2CEa02e"
	// buyersPublicAddress := "0x9cA67FFE69698d963A393E9338aD3BcfD2CEa02e"

	sellOrder := SellOrder{
		TradeAsset:            tradeAsset,
		Price:                 big.NewInt(1),
		Currency:              currency,
		Amount:                big.NewInt(2),
		SellerShippingAddress: sellersPublicAddress,
		SellerNKNAddress:      nknAddr,
		PaymentTransactionID:  sellersPublicAddress,
		RefundAddress:         sellersRefundAddress,
		Locked:                false,
	}

	// create a new sell order via tea for ETH and MO
	err, orderTxid := s.CreateSellOrderViaTea(sellOrder)
	if err != nil {
		fmt.Println("error creating sell order: ", err)
	}

	fmt.Println("SELL ORDER TXID: ", orderTxid)
	return orderTxid

	// List the orders via tea and verify that the order is in the list
	// orders, err := s.ListSellOrdersViaTea()
	// if err != nil {
	// 	(err)
	// }

	// // verify that the order is in the list
	// found := false
	// fmt.Println("looking for order: ", orderTxid)
	// for _, order := range orders {
	// 	fmt.Println("lookgin at order: ", order.TXID)
	// 	if order.TXID == orderTxid {
	// 		found = true
	// 		// retrieve the txid from the order
	// 		break
	// 	}
	// }

	// if !found {
	// 	("sell order not found in tea")
	// }

	// if err := s.VerifyOrderCreatedInReddis(orderTxid); err != nil {
	// 	(err)
	// }

	// // create a new buy transaction via tea
	// buyOrder := BuyOrder{
	// 	TXID:                 orderTxid,
	// 	BuyerShippingAddress: buyersPublicAddress,
	// 	BuyerNKNAddress:      nknAddr,
	// 	PaymentTransactionID: buyersPublicAddress,
	// 	RefundAddress:        buyersRefundAddress,
	// }

	// fmt.Println("buy order: ", buyOrder)

	// // initiate the transaction via tea
	// err = s.CreateBuyTransactionViaTea(buyOrder)
	// if err != nil {
	// 	(err)
	// 	return
	// }

	// if err := s.VerifyOrderLockedInReddis(orderTxid); err != nil {
	// 	(err)
	// }

	// // List the orders via tea and verify that the order is in the list
	// updatedOrders, err := s.ListSellOrdersViaTea()
	// if err != nil {
	// 	(err)
	// }

	// // verify that the order is in the list
	// found = false
	// for _, order := range updatedOrders {
	// 	if order.PaymentTransactionID == sellOrder.PaymentTransactionID {
	// 		found = true
	// 		// verify that the order is locked
	// 		if !order.Locked {
	// 			("sell order found, but not locked")
	// 			break
	// 		}
	// 		break
	// 	}
	// }

	// if !found {
	// 	("sell order not found in tea")
	// }

	// // List the orders via tea and verify that the order is not in the list
	// orders, err2 := s.ListSellOrdersViaTea()
	// if err2 != nil {
	// 	(err2)
	// }

	// // verify that the order is not in the list
	// found = false
	// for _, order := range orders {
	// 	if order.TXID == sellOrder.TXID {
	// 		found = true
	// 		break
	// 	}
	// }

	// if found {
	// 	("sell order found in tea, but should have been deleted")
	// }
}

func generateEthAccount() *ecdsa.PrivateKey {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println(err)
	}

	// privateKeyBytes := crypto.FromECDSA(privateKey)
	// pk := hexutil.Encode(privateKeyBytes)[2:]
	// fmt.Println("")
	// fmt.Printf("Generated Private Key: " + pk)
	return privateKey
}

// CreateSellOrderViaTea creates a new sell order via tea
func (s *Sipper) CreateSellOrderViaTea(order SellOrder) (error, string) {
	// prepare sellOrder to send in the http request
	sellOrderJSON, err := json.Marshal(order)
	if err != nil {
		return err, ""
	}

	io := bytes.NewBuffer(sellOrderJSON)
	req, err := http.NewRequest("POST", s.teaServer+"/sell", io)
	if err != nil {
		fmt.Println(err)
		return err, ""
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err, ""
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, ""
	}

	if resp.StatusCode != http.StatusAccepted {
		fmt.Printf("tea response Status:", resp.Status)
		fmt.Printf("tea response Headers:", resp.Header)
		fmt.Printf("response Body:", string(body))
		return errors.New("tea response status not accepted"), ""
	}

	fmt.Printf("tea response Body:", string(body))

	return nil, string(body)
}

// CreateBuyTransactionViaTea creates a new buy transaction via tea
func (s *Sipper) CreateBuyTransactionViaTea(order BuyOrder) error {
	// prepare buyOrder to send in the http request
	buyOrderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	io := bytes.NewBuffer(buyOrderJSON)
	req, err := http.NewRequest("POST", s.teaServer+"/buy", io)
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		fmt.Printf("tea response Status:", resp.Status)
		fmt.Printf("tea response Headers:", resp.Header)
		fmt.Printf("response Body:", string(body))
		return fmt.Errorf("tea server returned status code %d", resp.StatusCode)
	}

	return nil
}

// ListSellOrdersViaTea lists all the sell orders via tea
func (s *Sipper) ListSellOrdersViaTea() ([]SellOrder, error) {
	req, err := http.NewRequest("GET", s.teaServer+"/list", nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// if resp.StatusCode != http.StatusOK {
	// 	fmt.Printf("tea response Status:", resp.Status)
	// 	fmt.Printf("tea response Headers:", resp.Header)
	// 	fmt.Printf("response Body:", string(body))
	// 	return nil, fmt.Errorf("tea server returned status code %d", resp.StatusCode)
	// }

	var orders []SellOrder
	err = json.Unmarshal(body, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// RetrieveNKNAddress retrieves the NKN address of tea
func (s *Sipper) RetrieveNKNAddress() (string, error) {
	req, err := http.NewRequest("GET", s.teaServer+"/getNKNAddress", nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// if resp.StatusCode != http.StatusOK {
	// 	fmt.Printf("tea response Status:", resp.Status)
	// 	fmt.Printf("tea response Headers:", resp.Header)
	// 	fmt.Printf("response Body:", string(body))
	// 	return "", fmt.Errorf("tea server returned status code %d", resp.StatusCode)
	// }

	return string(body), nil
}

// // // CompleteOrder called to test the outcome of a successful transaction
// func (s *Sipper) CompleteOrder(o CompletedOrder) error {
// 	// for now we just ask the tester to fund the escrow accounts

// }

func sendEth(fromWallet *ecdsa.PrivateKey, toAddress string, amount *big.Int, rpc *ethclient.Client) error {
	// view the current balance of the paying wallet
	account := crypto.PubkeyToAddress(fromWallet.PublicKey)
	// send the currency to the buyer
	// read nonce
	nonce, err := rpc.PendingNonceAt(context.Background(), account)
	if err != nil {
		return err
	}

	// create gas params
	gasLimit := uint64(31000) // in units
	gasPrice, err := rpc.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	// convert the string address to an address
	qualifiedAddress := common.HexToAddress(toAddress)

	// create a transaction
	tx := types.NewTransaction(nonce, qualifiedAddress, amount, gasLimit, gasPrice, nil)

	// fetch chain id
	chainID, err := rpc.NetworkID(context.Background())
	if err != nil {
		return err
	}

	// sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromWallet)
	if err != nil {
		return err
	}

	// send the transaction
	err = rpc.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	fmt.Printf("tx sent: " + signedTx.Hash().Hex() + "txid: " + signedTx.Hash().Hex())
	return nil
}

// BuyerOnlyCondition called to test the outcome of a case when only
// the buyer has funded their respective escrow account

// SellerOnlyCondition called to test the outcome of a case when only
// the seller has funded their respective escrow account

// NoFundsCondition called to test the outcome of a case when neither
// the buyer or seller have funded their respective escrow account

func (s *Sipper) VerifyOrderLockedInReddis(orderID string) error {
	// get the order from redis
	order, err := s.redisClient.Get(context.Background(), "orders").Result()
	if err != nil {
		return err
	}

	// unmarshal the order
	var o []SellOrder
	err = json.Unmarshal([]byte(order), &o)
	if err != nil {
		return err
	}

	// find the order
	for _, v := range o {
		if v.TXID == orderID {
			if v.Locked {
				return nil
			}
			return fmt.Errorf("order status is not locked")
		}
	}

	return fmt.Errorf("order not found")
}

// VerifyOrderCreatedInReddis verifies that a new order is created in redis
func (s *Sipper) VerifyOrderCreatedInReddis(orderID string) error {
	// get the order from redis
	order, err := s.redisClient.Get(context.Background(), "orders").Result()
	if err != nil {
		return err
	}

	//unmarshal the order
	var o []SellOrder
	err = json.Unmarshal([]byte(order), &o)
	if err != nil {
		return err
	}

	// check that the order is in the list
	fmt.Println("checking for order:")
	for _, v := range o {
		if v.TXID == orderID {
			return nil
		}
	}

	return fmt.Errorf("order %s not found in redis", orderID)
}

// VerifyOrderUpdated verifies that a previously created order
// has been updated properly in redis

// VerifyCompleteOrderCreated verifies that a new completeOrder is created in redis

// VerifyCompletedOrderCreated verifies that a new completedOrder is created in redis
// when a transaction has been finalized/closed.

// KillAdamsAndVerifyOrdersStillExist kills the adams server and restarts it.
// then verifies that the orders are properly retrieved from redis.

// func openbrowser(url string) error {
// 	var err error
// 	switch runtime.GOOS {
// 	case "linux":
// 		err = exec.Command("xdg-open", url).Start()
// 	case "windows":
// 		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
// 	case "darwin":
// 		err = exec.Command("open", url).Start()
// 	default:
// 		err = fmt.Errorf("unsupported platform")
// 	}
// 	return err
// }

type AccountWatchRequestResult struct {
	AccountWatchRequest AccountWatchRequest `json:"account_watch_request"`
	Result              string              `json:"result"`
}

// AccountWatchRequest is the information we need to watch a new account
type AccountWatchRequest struct {
	Seller        bool     `json:"seller"`
	Account       string   `json:"account"`
	Chain         string   `json:"chain"`
	Amount        *big.Int `json:"amount"`
	TransactionID string   `json:"transaction_id"`
	TimeOut       int64    `json:"timeout"`
}

// // emulatePaymentEvent emulates a payment event by sending a cloudevent (or http POST request)
// // containing the bits that are expected in a real payment event
// func (s *Sipper) emulatePaymentEvents(orderID, chaina, chainb string) error {
// 	// create a new cloud event client
// 	c, err := cloudevents.NewClientHTTP()
// 	if err != nil {
// 		return err
// 	}

// 	// create a new AccountWatchRequestResult
// 	// and populate it with the information we need
// 	awrr := AccountWatchRequestResult{
// 		AccountWatchRequest: AccountWatchRequest{
// 			Seller:        true,
// 			Account:       "0x0000000000000000000000000000000000000000",
// 			Chain:         chaina,
// 			Amount:        big.NewInt(1000000000000000000),
// 			TransactionID: orderID,
// 			TimeOut:       0,
// 		},
// 		Result: "",
// 	}

// 	// create a new cloud event
// 	event := cloudevents.NewEvent()
// 	event.SetID(uuid.New().String())
// 	event.SetType("tea.party.watch.account.response")
// 	event.SetSource("warren")
// 	event.SetData(cloudevents.ApplicationJSON, awrr)

// 	// send the cloud event
// 	ctx := cloudevents.ContextWithTarget(context.Background(), s.brokerAddr)
// 	if result := c.Send(ctx, event); !cloudevents.IsACK(result) {
// 		return err
// 	}

// 	awrr2 := AccountWatchRequestResult{
// 		AccountWatchRequest: AccountWatchRequest{
// 			Seller:        false,
// 			Account:       "0x0000000000000000000000000000000000000000",
// 			Chain:         chainb,
// 			Amount:        big.NewInt(1000000000000000000),
// 			TransactionID: orderID,
// 			TimeOut:       0,
// 		},
// 		Result: "",
// 	}

// 	event2 := cloudevents.NewEvent()
// 	event2.SetID(uuid.New().String())
// 	event2.SetType("tea.party.watch.account.response")
// 	event2.SetSource("warren")
// 	event2.SetData(cloudevents.ApplicationJSON, awrr2)
// 	// send the cloud eve
// 	if result := c.Send(ctx, event2); !cloudevents.IsUndelivered(result) {
// 		return err
// 	}

// 	return nil
// }

func (s *Sipper) CreateAndLetFailWithNoReddisTesting(tradeAsset, currency string) {
	// retrieve the NKN address from tea
	nknAddr, err := s.RetrieveNKNAddress()
	if err != nil {
		fmt.Println(err)
	}

	sellersPublicAddress := crypto.PubkeyToAddress(s.moPK.PublicKey).Hex()
	buyersPublicAddress := crypto.PubkeyToAddress(s.ethPK.PublicKey).Hex()
	sellersRefundAddress := crypto.PubkeyToAddress(s.ethPK.PublicKey).Hex()
	buyersRefundAddress := crypto.PubkeyToAddress(s.moPK.PublicKey).Hex()

	sellOrder := SellOrder{
		TradeAsset:            tradeAsset,
		Price:                 big.NewInt(10),
		Currency:              currency,
		Amount:                big.NewInt(20),
		SellerShippingAddress: sellersPublicAddress,
		SellerNKNAddress:      nknAddr,
		PaymentTransactionID:  uuid.New().String(),
		RefundAddress:         sellersRefundAddress,
	}

	// create a new sell order via tea for ETH and MO
	err, orderTxid := s.CreateSellOrderViaTea(sellOrder)
	if err != nil {
		fmt.Println(err)
	}

	// List the orders via tea and verify that the order is in the list
	// orders, err := s.ListSellOrdersViaTea()
	// if err != nil {
	// 	(err)
	// }

	// verify that the order is in the list
	// found := false
	// fmt.Println("looking for order: ", orderTxid)
	// for _, order := range orders {
	// 	fmt.Println("looking at order: ", order.TXID)
	// 	if order.TXID == orderTxid {
	// 		fmt.Println("FOUND ORDER:", order)
	// 		found = true
	// 		// retrieve the txid from the order
	// 		break
	// 	}
	// }

	// if !found {
	// 	("sell order not found in tea")
	// }

	// create a new buy transaction via tea
	buyOrder := BuyOrder{
		TXID:                 orderTxid,
		BuyerShippingAddress: buyersPublicAddress,
		BuyerNKNAddress:      nknAddr,
		PaymentTransactionID: buyersPublicAddress,
		RefundAddress:        buyersRefundAddress,
	}

	// initiate the transaction via tea
	err = s.CreateBuyTransactionViaTea(buyOrder)
	if err != nil {
		fmt.Println(err)
		return
	}

}
