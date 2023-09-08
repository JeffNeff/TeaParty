package be

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	bridge "github.com/teapartycrypto/TeaParty/adams/pkg/contract/bridge"
)

func TestWatchContractBalance(t *testing.T) {
	node, err := ethclient.Dial("https://rpc.octa.space/")
	if err != nil {
		t.Fatal(err)
	}
	// TODO: this should be pulled in as a env variable into the exchange server
	contractAddress := common.HexToAddress("0x0eeAaF074B23942CD660175dEaE6e1A5849d6614")
	contract, err := bridge.NewBe(contractAddress, node)
	if err != nil {
		t.Fatal(err)
	}

	balance, err := contract.BalanceOf(nil, common.HexToAddress("0xe2dE4EAE45225052f8F1b782D48076DA00d6A4b0"))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(balance)
}
