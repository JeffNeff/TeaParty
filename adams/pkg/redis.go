package be

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	uuid "github.com/google/uuid"
)

func (e *ExchangeServer) updateOrdersInDB(orders []SellOrder) error {
	// fetch the orders from the database
	dbOrders, err := e.fetchOrdersFromDB()
	if err != nil {
		return err
	}

	// compare the orders in the database with the orders in the request
	// if the orders in the request are not in the database, add them
	var found bool
	for _, oo := range orders {
		found = false
		for _, o := range dbOrders {
			if o.TXID == oo.TXID {
				found = true
				break
			}
		}
		if !found {
			dbOrders = append(dbOrders, oo)
		}
	}

	// store the order in the database
	json, err := json.Marshal(dbOrders)
	if err != nil {
		e.logger.Error(err)
		return err
	}

	err = e.redisClient.Set(context.Background(), "orders", json, 0).Err()
	if err != nil {
		e.logger.Error(err)
		return err
	}

	return nil
}

func (e *ExchangeServer) setOrders(orders []SellOrder) error {
	// store the order in the database
	ojs, err := json.Marshal(orders)
	if err != nil {
		return err
	}

	err = e.redisClient.Set(context.Background(), "orders", ojs, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// func (e *ExchangeServer) completeOrdersInDB(orders []CompletedOrder) error {
// 	// store the order in the database
// 	// ojs, err := json.Marshal(orders)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	err := e.redisClient.Set(context.Background(), "completeorders", orders, 0).Err()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (e *ExchangeServer) fetchCompleteOrdersFromDB() ([]CompletedOrder, error) {
	// fetch the orders from the database
	orders, err := e.redisClient.Get(context.Background(), "completeorders").Result()
	if err != nil {
		return nil, err
	}

	// check if orders is empty or represents an empty array or null value
	if orders == "" {
		return make([]CompletedOrder, 0), nil
	}
	var o []CompletedOrder
	if err := json.Unmarshal([]byte(orders), &o); err != nil {
		if _, ok := err.(*json.SyntaxError); ok {
			return make([]CompletedOrder, 0), nil
		}
		return nil, err
	}
	return o, nil
}

func (e *ExchangeServer) updateCompleteOrdersInDB(orders []CompletedOrder) error {
	// fetch the complete orders from the database
	dbOrders, err := e.fetchCompleteOrdersFromDB()
	if err != nil {
		return err
	}

	// compare the orders in the database with the orders in the request
	// if the orders in the request are not in the database, add them
	var found bool
	for _, oo := range orders {
		found = false
		for _, o := range dbOrders {
			if o.OrderID == oo.OrderID {
				found = true
				break
			}
		}
		if !found {
			dbOrders = append(dbOrders, oo)
		}
	}

	// store the order in the database
	ojs, err := json.Marshal(dbOrders)
	if err != nil {
		return err
	}

	err = e.redisClient.Set(context.Background(), "completeorders", ojs, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// updateCompleteOrdersInDB function
func (e *ExchangeServer) updateCompleteOrdersInDBWithIndex(orders []CompletedOrder, indexToUpdate int) error {
	// fetch the complete orders from the database
	dbOrders, err := e.fetchCompleteOrdersFromDB()
	if err != nil {
		return err
	}

	// compare the orders in the database with the orders in the request
	// if the orders in the request are not in the database, add them
	var found bool
	for _, oo := range orders {
		found = false
		for i, o := range dbOrders {
			if o.OrderID == oo.OrderID {
				found = true
				if i == indexToUpdate {
					dbOrders[i] = oo
				}
				break
			}
		}
		if !found {
			dbOrders = append(dbOrders, oo)
		}
	}

	// store the order in the database
	ojs, err := json.Marshal(dbOrders)
	if err != nil {
		return err
	}

	err = e.redisClient.Set(context.Background(), "completeorders", ojs, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (e *ExchangeServer) fetchFailedOrdersFromDB() ([]CompletedOrder, error) {
	// fetch the orders from the database
	orders, err := e.redisClient.Get(context.Background(), "failedorders").Result()
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return make([]CompletedOrder, 0), nil
	}

	// unmarshal the orders
	var o []CompletedOrder
	err = json.Unmarshal([]byte(orders), &o)
	if err != nil {
		return make([]CompletedOrder, 0), err
	}

	return o, nil
}

func (e *ExchangeServer) updateFailedOrdersInDB(orders CompletedOrder, failureReason string) error {
	// fetch the failed orders from the database
	dbOrders, err := e.fetchFailedOrdersFromDB()
	if err != nil {
		return err
	}

	// compare the orders in the database with the orders in the request
	// if the orders in the request are not in the database, add them
	for _, o := range dbOrders {

		if o.OrderID == orders.OrderID {
			break
		} else {
			dbOrders = append(dbOrders, orders)
		}
	}

	orders.FailureReason = failureReason
	// store the order in the database
	ojs, err := json.Marshal(dbOrders)
	if err != nil {
		return err
	}

	err = e.redisClient.Set(context.Background(), "failedorders", ojs, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (e *ExchangeServer) fetchOrdersFromDB() ([]SellOrder, error) {
	// fetch the orders from the database
	orders, err := e.redisClient.Get(context.Background(), "orders").Result()
	if err != nil {
		return nil, err
	}

	if orders == "" {
		return make([]SellOrder, 0), nil
	}

	// unmarshal the orders
	var o []SellOrder
	if err := json.Unmarshal([]byte(orders), &o); err != nil {
		return nil, err
	}

	return o, nil
}

// updateAccountWatchRequestInDB updates the account watch request in the database
// so that it can be recovered in the event of a crash.
func (e *ExchangeServer) updateAccountWatchRequestInDB(request AccountWatchRequest) error {
	// retrieve the list of current account watch requests
	requests, _ := e.redisClient.Get(context.Background(), "accountwatchrequests").Result()
	// if err != nil {
	// 	return err
	// }``
	// unmarshal the list of account watch requests
	var currentRequests []AccountWatchRequest
	if requests != "" {
		err := json.Unmarshal([]byte(requests), &currentRequests)
		if err != nil {
			return err
		}
	}
	// if the request exists in the list, update it
	var found bool
	for i, r := range currentRequests {
		u1, err := uuid.Parse(r.AWRID)
		if err != nil {
			return err
		}
		u2, err := uuid.Parse(request.AWRID)
		if err != nil {
			return err
		}

		if u1 == u2 {
			currentRequests[i] = request
			found = true
		}
	}
	// if the request does not exist in the list, add it
	if !found {
		currentRequests = append(currentRequests, request)
	}
	// marshal the list of account watch requests
	crjs, err := json.Marshal(currentRequests)
	if err != nil {

		return err
	}
	// store the new list of account watch requests
	err = e.redisClient.Set(context.Background(), "accountwatchrequests", crjs, 0).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// update all account watch requests to reflect that the exchange server has crashed
// so we need to unlock them so that they can be processed again by another exchange server
func (e *ExchangeServer) updateAccountWatchRequestsOnCrash() error {
	// retrieve the list of current account watch requests
	requests, _ := e.redisClient.Get(context.Background(), "accountwatchrequests").Result()
	// if err != nil {
	// 	return err
	// }

	// unmarshal the list of account watch requests
	var currentRequests []AccountWatchRequest
	if requests != "" {
		err := json.Unmarshal([]byte(requests), &currentRequests)
		if err != nil {
			return err
		}
	}

	// compare the exchange server ID of each request to the ID of this exchange server
	// if they match, unlock the request
	for i, r := range currentRequests {
		if r.LockedBy == e.exchangeServerID {
			currentRequests[i].Locked = false
			currentRequests[i].LockedBy = ""
		}
	}

	// marshal the list of account watch requests
	crjs, err := json.Marshal(currentRequests)
	if err != nil {
		return err
	}

	// store the new list of account watch requests
	err = e.redisClient.Set(context.Background(), "accountwatchrequests", crjs, 0).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// removeAccountWatchRequestFromDB removes the account watch request from the database
// after it has been processed.
func (e *ExchangeServer) removeAccountWatchRequestFromDB(request *AccountWatchRequest) error {
	// fetch the list of current account watch requests
	requests, _ := e.redisClient.Get(context.Background(), "accountwatchrequests").Result()
	// unmarshal the list of account watch requests
	var currentRequests []AccountWatchRequest
	if requests != "" {
		err := json.Unmarshal([]byte(requests), &currentRequests)
		if err != nil {
			return err
		}
	}
	// search for the request in the list by the transaction ID
	// if it is found, remove it from the list
	for i, r := range currentRequests {
		if r.TransactionID == request.TransactionID {
			currentRequests = append(currentRequests[:i], currentRequests[i+1:]...)
		}
	}

	// marshal the list of account watch requests
	crjs, err := json.Marshal(currentRequests)
	if err != nil {
		return err
	}

	// store the new list of account watch requests
	err = e.redisClient.Set(context.Background(), "accountwatchrequests", crjs, 0).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// retrieveAccountWatchRequestsFromDB retrieves the account watch requests from the database
// so that they can be processed.
func (e *ExchangeServer) retrieveAccountWatchRequestsFromDB() ([]AccountWatchRequest, error) {
	// fetch the list of current account watch requests
	requests, _ := e.redisClient.Get(context.Background(), "accountwatchrequests").Result()
	// unmarshal the list of account watch requests

	if len(requests) == 0 {
		return make([]AccountWatchRequest, 0), nil
	}

	var currentRequests []AccountWatchRequest
	err := json.Unmarshal([]byte(requests), &currentRequests)
	if err != nil {
		return nil, err
	}

	return currentRequests, nil
}

// addAssistedSellOrderToFailedAssistedSellOrdersDB adds a failed assisted sell order to the database
// so that it can be recovered in the event of a crash.
func (e *ExchangeServer) addAssistedSellOrderToFailedAssistedSellOrdersDB(order *AccountWatchRequestResult) error {
	// retrieve the list of current failed assisted sell orders
	orders, _ := e.redisClient.Get(context.Background(), "failedassistedsellorders").Result()
	// if err != nil {
	// 	return err
	// }
	// unmarshal the list of failed assisted sell orders
	var currentOrders []AccountWatchRequestResult
	if orders != "" {
		err := json.Unmarshal([]byte(orders), &currentOrders)
		if err != nil {
			return err
		}
	}
	// add the order to the list
	currentOrders = append(currentOrders, *order)
	cojs, err := json.Marshal(currentOrders)
	if err != nil {
		return err
	}
	// store the new list of failed assisted sell orders
	err = e.redisClient.Set(context.Background(), "failedassistedsellorders", cojs, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

type BridgeStorage struct {
	Chain      string   `json:"chain"`
	PrivateKey string   `json:"privatekey"`
	Amount     *big.Int `json:"amount"`
	Asset      string   `json:"asset"`
	ID         string   `json:"id"`
}

// storeBridgeAccount we need a function to store the private keys and the amount of each coin held in the database
// so that when users try to bring a wrapped asset back on the chain, we can fund the transaction
// without having to move funds around.
func (e *ExchangeServer) storeBridgeAccount(awrr AccountWatchRequestResult) error {
	// create a bridge storage object out of the account watch request result
	bs := BridgeStorage{
		Chain:      awrr.AccountWatchRequest.Chain,
		PrivateKey: awrr.AccountWatchRequest.AssistedSellOrderInformation.SellersEscrowWallet.PrivateKey,
		Amount:     awrr.AccountWatchRequest.AssistedSellOrderInformation.Amount,
		Asset:      awrr.AccountWatchRequest.AssistedSellOrderInformation.Currency,
		ID:         awrr.AccountWatchRequest.TransactionID,
	}

	// marshal the bridge account
	bsjs, err := json.Marshal(bs)
	if err != nil {
		return err
	}
	// append the values to the database
	err = e.redisClient.HSet(context.Background(), "bridgeaccounts", bs.Chain, bsjs).Err()
	if err != nil {
		return err
	}

	return nil
}

// retrieveBridgeAccount retrieves the bridge account from the database
func (e *ExchangeServer) retrieveBridgeAccount(valueSearch *big.Int) (*BridgeStorage, error) {
	if valueSearch == nil {
		return nil, nil
	}

	// retrieve the list of bridge accounts
	accounts, _ := e.redisClient.Get(context.Background(), "bridgeaccounts").Result()
	// unmarshal the list of bridge accounts
	var currentAccounts []BridgeStorage
	if accounts != "" {
		err := json.Unmarshal([]byte(accounts), &currentAccounts)
		if err != nil {
			return nil, err
		}
	}

	// search for the account in the list and find the one with the closest amount
	var closestAccount BridgeStorage
	for _, a := range currentAccounts {
		if a.Amount.Cmp(valueSearch) == 0 {
			return &a, nil
		}
		if a.Amount.Cmp(valueSearch) == 1 {
			if closestAccount.Amount == nil {
				closestAccount = a
			} else {
				if a.Amount.Cmp(closestAccount.Amount) == -1 {
					closestAccount = a
				}
			}
		}
	}

	return &closestAccount, nil
}
