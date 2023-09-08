package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// NewAdapter adapter implementation
func main() {

	// initialize the Party Chain nodes.
	partyclient, err := ethclient.Dial("http://192.168.50.193:8545")
	if err != nil {
		panic(err)
	}

	privateKey, err := crypto.HexToECDSA("8e78aa5c1a7bdc369a7c6b540399904778a90b2efbd1c7db16b68d1085b39098")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := partyclient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println("failed to get nonce: %v", err)
		log.Fatal(err)
	}

	gasPrice, err := partyclient.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Printf("failed to get suggest gas price: %v", err)
		log.Fatal(err)
	}

	address := common.HexToAddress("0xDAD416F84E67d4B37c12c37979DC4E8d07Fc83d2")
	instance, err := NewTeaParty(address, partyclient)
	if err != nil {
		fmt.Printf("failed to instantiate a Token contract: %v", err)
		log.Fatal(err)
	}

	fmt.Println(instance)

	chainID, err := partyclient.ChainID(context.Background())
	if err != nil {
		fmt.Printf("failed to get chainID: %v", err)
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice
}
