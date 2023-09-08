package be

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (e *ExchangeServer) verifyPaymentTransactionID(txid string) error {
	url := e.blockExplorer + "/api?module=transaction&action=gettxinfo&txhash=" + txid
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	bstqr := &BlockScoutTxQueryResponse{}
	err = json.Unmarshal(body, bstqr)
	if err != nil {
		return err
	}

	if bstqr.Result.Success == false {
		return fmt.Errorf("transaction not found")
	}

	return nil
}
