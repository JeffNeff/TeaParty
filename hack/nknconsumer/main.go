package main

import (
	"log"

	nkn "github.com/nknorg/nkn-sdk-go"
)

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
}

func start() error {
	account, err := nkn.NewAccount(nil)
	if err != nil {
		return err
	}

	toClient, err := nkn.NewMultiClient(account, "", 4, false, nil)
	if err != nil {
		return err
	}
	defer toClient.Close()

	log.Printf("listening on %s", toClient.Address())
	<-toClient.OnConnect.C
	for {
		msg := <-toClient.OnMessage.C
		isEncryptedStr := "unencrypted"
		if msg.Encrypted {
			isEncryptedStr = "encrypted"
		}
		log.Println("Receive", isEncryptedStr, "message", "\""+string(msg.Data)+"\"", "from", msg.Src)
		msg.Reply([]byte("ok"))
	}

	return nil
}
