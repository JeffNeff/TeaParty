package main

import (
	"fmt"
	"log"

	nkn "github.com/nknorg/nkn-sdk-go"
)

func main() {
	err := start()
	if err != nil {
		fmt.Println(err)
	}
}

func start() error {
	ca := "467fa1b7b6aa78806198f5fe2b1d6b37eef6748b306feadc9ce41291c9bb8be6"

	account, err := nkn.NewAccount(nil)
	if err != nil {
		return err
	}

	fromClient, err := nkn.NewMultiClient(account, "", 4, false, nil)
	if err != nil {
		return err
	}
	defer fromClient.Close()
	<-fromClient.OnConnect.C

	log.Println("Send message from", fromClient.Address(), "to", ca)
	// []byte("Hello") can be replaced with "Hello" for text payload type
	onReply, err := fromClient.Send(nkn.NewStringArray(ca), []byte(`{"bob":"workd"}`), nil)
	if err != nil {
		return err
	}

	reply := <-onReply.C
	isEncryptedStr := "unencrypted"
	if reply.Encrypted {
		isEncryptedStr = "encrypted"
	}

	log.Println("Got", isEncryptedStr, "reply", "\""+string(reply.Data)+"\"", "from", reply.Src, "after")

	return nil
}
