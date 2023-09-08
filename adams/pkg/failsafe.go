package be

import (
	"encoding/json"
	"os"
)

// writeFailedTransactionToFile writes a failed transaction to a file
// when we cannot connect to redis
func (e *ExchangeServer) writeFailedTransactionToFile(order CompletedOrder) error {
	// store the order in the failed transaction file locally
	file := "failedtransactions" + order.OrderID + ".json"

	// convert the order to []byte
	border, err := json.Marshal(order)
	if err != nil {
		return err
	}
	// write the orders to file
	err2 := os.WriteFile(file, border, 0644)
	if err2 != nil {
		return err2
	}

	return nil
}
