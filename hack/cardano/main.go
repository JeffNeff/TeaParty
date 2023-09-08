package main

import (
	"fmt"

	// "github.com/echovl/cardano-go"
	cardanocli "github.com/echovl/cardano-go/cardano-cli"
	// "github.com/echovl/cardano-go/wallet"
	"github.com/echovl/cardano-go"
	"github.com/echovl/cardano-go/wallet"
)

func main() {
	// node := blockfrost.NewNode(cardano.Mainnet, "mainnetrNpgKN1W4M6Cv4PVkLtgcDWEjkcFKZpb")

	// pparams, err := node.ProtocolParams()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(pparams)

	// create a new wallet client
	walletClient := wallet.NewClient(&wallet.Options{
		Node: cardanocli.NewNode(cardano.Mainnet),
	})
	// fmt.Println(walletClient)

	// create a new wallet
	wallet, memonic, err := walletClient.CreateWallet("MyWallet", "MyWallet")
	if err != nil {
		panic(err)
	}

	fmt.Println(memonic)

	fmt.Println(fmt.Sprintf("%+v", wallet))

	// balance, err := node.GetBalance(wallet.ID)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(wallet.Balance())

}

type memoryDB struct {
	wm map[string]*wallet.Wallet
}

// type Wallet struct {
// 	ID       string
// 	Name     string
// 	addrKeys []crypto.XPrvKey
// 	stakeKey crypto.XPrvKey
// 	rootKey  crypto.XPrvKey
// 	node     cardano.Node
// 	network  cardano.Network
// }
