package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Token struct {
	Total      int    `json:"total"`
	PageNumber int    `json:"pageNumber"`
	PageSize   int    `json:"pageSize"`
	Cursor     any    `json:"cursor"`
	Network    string `json:"network"`
	Owners     []struct {
		TokenAddress      string `json:"tokenAddress"`
		TokenID           string `json:"tokenId"`
		Amount            string `json:"amount"`
		OwnerOf           string `json:"ownerOf"`
		TokenHash         string `json:"tokenHash"`
		BlockNumberMinted string `json:"blockNumberMinted"`
		BlockNumber       string `json:"blockNumber"`
		ContractType      string `json:"contractType"`
		Name              string `json:"name"`
		Symbol            string `json:"symbol"`
		Metadata          any    `json:"metadata"`
		MinterAddress     string `json:"minterAddress"`
	} `json:"owners"`
}

func main() {
	url := "https://nft.api.infura.io/networks/137/nfts/0x970A8b10147E3459D3CBF56329B76aC18D329728/4/owners"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic NjE5Nzk3OTdhOGJiNGJmZTlkZGRkNGZmOTY3NWRiN2U6ODBjM2ZjZjk2ZmFhNDlmN2JkNmIzMmNjYzQxOTAyZDg=")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(body))

	var token Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", token.Owners[0].OwnerOf)

	fmt.Printf("%+v", strings.Contains(token.Owners[0].OwnerOf, "0xf5784fe0650999bb17ddff0b095430d6a86ba714"))
}
