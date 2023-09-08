package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func main() {
	// Create a new account:
	account := solana.NewWallet()
	fmt.Println("account private key:", account.PrivateKey)
	fmt.Println("account public key:", account.PublicKey())

	// Create a new RPC client:
	client := rpc.New("http://localhost:8899")
	// client := rpc.New(rpc.TestNet_RPC)

	// view the account balance:
	out, err := client.GetBalance(
		context.TODO(),
		account.PublicKey(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("account balance: %v SOL \n", out.Value)

	// Airdrop 5 SOL to the new account:
	out2, err := client.RequestAirdrop(
		context.TODO(),
		account.PublicKey(),
		solana.LAMPORTS_PER_SOL*1,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("airdrop transaction signature:", out2)

	// view the account balance:
	out, err = client.GetBalance(
		context.TODO(),
		account.PublicKey(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		panic(err)
	}

	// sleep for 5 seconds to wait for the airdrop to be confirmed:
	time.Sleep(60 * time.Second)

	fmt.Printf("account balance: %v SOL \n", out.Value)

}
