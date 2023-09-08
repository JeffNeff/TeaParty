package be

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	nkn "github.com/nknorg/nkn-sdk-go"
)

type NKNNotification struct {
	Address    string `json:"address"`
	Amount     string `json:"amount"`
	Network    string `json:"network"`
	PrivateKey string `json:"privateKey"`
	Chain      string `json:"chain"`
	Error      string `json:"error"`
}

// notifySellerOfBuyer notifies the seller that the buyer has appeared for the trade
// and presents the buyer with the buyer's escrow wallet address.
func (e *ExchangeServer) notifySellerOfBuyer(co CompletedOrder) error {
	account, err := nkn.NewAccount(nil)
	if err != nil {
		e.logger.Error("creating a new nkn acount failed")
		return err
	}

	// initialize the nkn client.
	nknClient, err := nkn.NewMultiClient(account, "", 4, false, nil)
	if err != nil {
		e.logger.Error("creating a new nkn client failed")
		return err
	}

	defer nknClient.Close()

	sn := &NKNNotification{
		Address: co.SellerEscrowWallet.PublicAddress,
		Amount:  co.Amount.String(),
		Network: co.Currency,
	}

	bytes, err := json.Marshal(sn)
	if err != nil {
		e.logger.Error("error marshalling seller notification: " + err.Error())
		return err
	}

	<-nknClient.OnConnect.C
	onReply, err := nknClient.Send(nkn.NewStringArray(co.SellerNKNAddress), bytes, nil)
	if err != nil {
		return err
	}

	// create a timeout of 2 minutes
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

	defer cancel()

	select {
	case reply := <-onReply.C:
		switch string(reply.Data) {
		case "ok":
			e.logger.Info("seller has received notification of new buyer for order: " + co.OrderID)
			return nil
		default:
			e.logger.Error("seller has encountered an error for order: " + co.OrderID)
			return fmt.Errorf("seller has encountered an error for order: " + co.OrderID)
		}
	case <-ctx.Done():
		e.logger.Error("seller has not responded to notification of new buyer for order: " + co.OrderID)
		return fmt.Errorf("seller has not responded to notification of new buyer for order: " + co.OrderID)
	}

}

// sendBuyerPayInfo sends the buyer the escrow information
func (e *ExchangeServer) sendBuyerPayInfo(co CompletedOrder) error {
	account, err := nkn.NewAccount(nil)
	if err != nil {
		e.logger.Error("creating a new nkn acount failed")
		return err
	}

	// initialize the nkn client.
	nknClient, err := nkn.NewMultiClient(account, "", 4, false, nil)
	if err != nil {
		e.logger.Error("creating a new nkn client failed")
		return err
	}
	defer nknClient.Close()

	sn := &NKNNotification{
		Address: co.BuyerEscrowWallet.PublicAddress,
		Amount:  co.Price.String(),
		Network: co.TradeAsset,
	}

	bytes, err := json.Marshal(sn)
	if err != nil {
		e.logger.Error("error marshalling buyer notification: " + err.Error())
		return err
	}

	<-nknClient.OnConnect.C
	onReply, err := nknClient.Send(nkn.NewStringArray(co.BuyerNKNAddress), bytes, nil)
	if err != nil {
		return err
	}

	// create a timeout of 2 minutes
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

	defer cancel()

	select {
	case reply := <-onReply.C:
		switch string(reply.Data) {
		case "ok":
			e.logger.Info("buyer has received escrow information for order: " + co.OrderID)
			return nil
		default:
			e.logger.Error("buyer has encountered an error for order: " + co.OrderID)
			return fmt.Errorf("buyer has encountered an error for order: " + co.OrderID)
		}
	case <-ctx.Done():
		e.logger.Error("buyer has not responded to notification of escrow information for order: " + co.OrderID)
		return fmt.Errorf("buyer has not responded to notification of escrow information for order: " + co.OrderID)
	}

}

// notifyBothPartiesOfTradeCancelation notifies the seller and the buyer that the trade has been canceled
func (e *ExchangeServer) notifyBothPartiesOfTradeCancelation(co CompletedOrder) error {
	sellerAccount, err := nkn.NewAccount(nil)
	if err != nil {
		e.logger.Error("creating a new nkn acount failed")
		return err
	}

	// initialize the nkn client.
	sellerNKNClient, err := nkn.NewMultiClient(sellerAccount, "", 4, false, nil)
	if err != nil {
		e.logger.Error("creating a new nkn client failed")
		return err
	}
	defer sellerNKNClient.Close()

	buyerAccount, err := nkn.NewAccount(nil)
	if err != nil {
		e.logger.Error("creating a new nkn acount failed")
		return err
	}

	// initialize the nkn client.
	buyerNKNClient, err := nkn.NewMultiClient(buyerAccount, "", 4, false, nil)
	if err != nil {
		e.logger.Error("creating a new nkn client failed")
		return err
	}
	defer buyerNKNClient.Close()

	// create a timeout of 2 minutes
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	sn := &NKNNotification{
		Error: "trade" + co.OrderID + " canceled",
	}

	bytes, err := json.Marshal(sn)
	if err != nil {
		e.logger.Error("error marshalling buyer notification: " + err.Error())
		return err
	}

	<-sellerNKNClient.OnConnect.C
	onSellerReply, err := sellerNKNClient.Send(nkn.NewStringArray(co.SellerNKNAddress), bytes, nil)
	if err != nil {
		return err
	}

	var sendingError error

	select {
	case reply := <-onSellerReply.C:
		switch string(reply.Data) {
		case "ok":
			e.logger.Info("seller has received notification")
			return nil
		default:
			e.logger.Error("seller has encountered an error")
			sendingError = fmt.Errorf("seller has encountered an error")
		}
	case <-ctx.Done():
		e.logger.Error("seller has not responded to notification of cancelation for order: " + co.OrderID)
		return fmt.Errorf("seller has not responded to notification of cancelation for order: " + co.OrderID)
	}

	onBuyerReply, err := buyerNKNClient.Send(nkn.NewStringArray(co.BuyerNKNAddress), bytes, nil)
	if err != nil {
		return err
	}

	select {
	case buyerReply := <-onBuyerReply.C:
		switch string(buyerReply.Data) {
		case "ok":
			e.logger.Info("buyer has received notification")
			return nil
		default:
			e.logger.Error("buyer has encountered an error recieving close order notification")
			sendingError = fmt.Errorf("buyer has encountered an error recieving close order notification")
		}
	case <-ctx.Done():
		e.logger.Error("buyer has not responded to notification of cancelation for order: " + co.OrderID)
		return fmt.Errorf("buyer has not responded to notification of cancelation for order: " + co.OrderID)
	}

	return sendingError
}

// sendSellerEscrowWalletPrivateKeyToBuyer is called to send the private key of the sellers
// escrow wallet to the buyer
func (e *ExchangeServer) sendSellerEscrowWalletPrivateKeyToBuyer(order CompletedOrder) error {
	account, err := nkn.NewAccount(nil)
	if err != nil {
		e.logger.Error("creating a new nkn acount failed")
		return err
	}

	// initialize the nkn client.
	nknClient, err := nkn.NewMultiClient(account, "", 4, false, nil)
	if err != nil {
		e.logger.Error("creating a new nkn client failed")
		return err
	}

	ac := &NKNNotification{
		PrivateKey: order.SellerEscrowWallet.PrivateKey,
		Address:    order.SellerEscrowWallet.PublicAddress,
		Chain:      order.Currency,
	}

	bytes, err := json.Marshal(ac)
	if err != nil {
		e.logger.Error("error marshalling seller notification: " + err.Error())
		return err
	}

	// send the private key to the seller
	e.logger.Infof("sending private key: %s to buyer: %s", ac.PrivateKey, order.BuyerNKNAddress)

	<-nknClient.OnConnect.C
	onReply, err := nknClient.Send(nkn.NewStringArray(order.BuyerNKNAddress), bytes, nil)
	if err != nil {
		er := fmt.Sprintf("sending private key: %s to buyer: %s failed.", order.SellerEscrowWallet.PrivateKey, order.BuyerNKNAddress)
		e.logger.Infof(er)
		// emit new error event
		return err
	}

	// create a timeout of 2 minutes
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

	defer cancel()

	select {
	case reply := <-onReply.C:
		switch string(reply.Data) {
		case "ok":
			e.logger.Info("buyer has received seller escrow wallet private key for order: " + order.OrderID)
			return nil
		default:
			e.logger.Error("buyer has encountered an error receiving sellers escrow wallet private key for order: " + order.OrderID)
			return fmt.Errorf("buyer has encountered an error receiving sellers escrow wallet private key for order: " + order.OrderID)
		}
	case <-ctx.Done():
		e.logger.Error("buyer has not responded to delivery of the PK: " + order.OrderID)
		return fmt.Errorf("buyer has not responded to notification of new PK for order: " + order.OrderID)
	}
	return nil
}

// sendBuyerEscrowWalletPrivateKeyToSeller is called to send the private key of the buyers escrow wallet
// to the seller.
func (e *ExchangeServer) sendBuyerEscrowWalletPrivateKeyToSeller(order CompletedOrder) error {
	account, err := nkn.NewAccount(nil)
	if err != nil {
		e.logger.Error("creating a new nkn acount failed")
		return err
	}

	// initialize the nkn client.
	nknClient, err := nkn.NewMultiClient(account, "", 4, false, nil)
	if err != nil {
		e.logger.Error("creating a new nkn client failed")
		return err
	}
	ac := &NKNNotification{
		PrivateKey: order.BuyerEscrowWallet.PrivateKey,
		Address:    order.BuyerEscrowWallet.PublicAddress,
		Chain:      order.TradeAsset,
	}

	bytes, err := json.Marshal(ac)
	if err != nil {
		e.logger.Error("error marshalling buyer notification: " + err.Error())
		return err
	}

	defer nknClient.Close()
	e.logger.Infof("sending private key: %s to seller: %s", order.BuyerEscrowWallet.PrivateKey, order.SellerNKNAddress)
	<-nknClient.OnConnect.C
	onReply, err := nknClient.Send(nkn.NewStringArray(order.SellerNKNAddress), bytes, nil)
	if err != nil {
		er := fmt.Sprintf("sending private key: %s to seller: %s failed.", order.BuyerEscrowWallet.PrivateKey, order.BuyerNKNAddress)
		e.logger.Infof(er)
		return err
	}

	// create a timeout of 2 minutes
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

	defer cancel()

	select {
	case reply := <-onReply.C:
		switch string(reply.Data) {
		case "ok":
			e.logger.Info("seller has received buyers escrow wallet private key for order: " + order.OrderID)
			nknClient.Close()
			return nil
		default:
			e.logger.Error("seller has encountered an error receiving buyers escrow wallet private key for order: " + order.OrderID)
			nknClient.Close()
			return fmt.Errorf("seller has encountered an error receiving buyers escrow wallet private key for order: " + order.OrderID)
		}
	case <-ctx.Done():
		e.logger.Error("seller has not responded to delivery of the PK: " + order.OrderID)
		nknClient.Close()
		return fmt.Errorf("seller has not responded to notification of new PK for order: " + order.OrderID)
	}
}

// pingNKNAddress is called to ping the NKN address of the buyer or seller to ensure that
// the address is valid.
func (e *ExchangeServer) pingNKNAddress(address string) error {
	account, err := nkn.NewAccount(nil)
	if err != nil {
		e.logger.Error("creating a new nkn acount failed")
		return err
	}

	// initialize the nkn client.
	nknClient, err := nkn.NewMultiClient(account, "", 4, false, nil)
	if err != nil {
		e.logger.Error("creating a new nkn client failed")
		return err
	}

	defer nknClient.Close()
	e.logger.Infof("pinging NKN address: %s", address)

	pn := &PingNotification{
		PingNotification: "ping",
	}

	bytes, err := json.Marshal(pn)
	if err != nil {
		e.logger.Error("error marshalling ping notification: " + err.Error())
		return err
	}

	<-nknClient.OnConnect.C
	onReply, err := nknClient.Send(nkn.NewStringArray(address), bytes, nil)
	if err != nil {
		e.logger.Error("pinging NKN address: " + address + " failed.")
		return err
	}

	// create a timeout of 2 minutes
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	select {
	case reply := <-onReply.C:
		switch string(reply.Data) {
		case "ok":
			e.logger.Info("NKN address: " + address + " is valid")
			return nil
		default:
			e.logger.Error("NKN address: " + address + " is invalid")
			return fmt.Errorf("NKN address: " + address + " is invalid")
		}
	case <-ctx.Done():
		e.logger.Error("NKN address: " + address + " is invalid")
		return fmt.Errorf("NKN address: " + address + " is invalid")
	}

}

type PingNotification struct {
	PingNotification string `json:"ping_notification"`
}

// refundSellerViaEscrowWallet is called to refund the seller via the escrow wallet
func (e *ExchangeServer) refundSellerViaEscrowWallet(order CompletedOrder) error {
	account, err := nkn.NewAccount(nil)
	if err != nil {
		e.logger.Error("creating a new nkn acount failed")
		return err
	}

	// initialize the nkn client.
	nknClient, err := nkn.NewMultiClient(account, "", 4, false, nil)
	if err != nil {
		e.logger.Error("creating a new nkn client failed")
		return err
	}

	defer nknClient.Close()
	e.logger.Info("refunding seller via NKN by sending the private key of the seller's escrow wallet")

	// send the private key to the seller
	e.logger.Infof("sending private key: %s to seller: %s", order.SellerEscrowWallet.PrivateKey, order.SellerNKNAddress)

	ac := &NKNNotification{
		PrivateKey: order.SellerEscrowWallet.PrivateKey,
		Chain:      order.Currency,
	}

	bytes, err := json.Marshal(ac)
	if err != nil {
		e.logger.Error("error marshalling  refund seller via escrow wallet: " + err.Error())
		return err
	}

	<-nknClient.OnConnect.C
	onReply, err := nknClient.Send(nkn.NewStringArray(order.SellerNKNAddress), bytes, nil)
	if err != nil {
		er := fmt.Sprintf("sending private key: %s to seller: %s failed.", order.SellerEscrowWallet.PrivateKey, order.SellerNKNAddress)
		e.logger.Infof(er)
		return err
	}

	// create a timeout of 2 minutes
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

	defer cancel()

	select {
	case reply := <-onReply.C:
		switch string(reply.Data) {
		case "ok":
			e.logger.Info("seller has received refund via escrow wallet")
			nknClient.Close()
			return nil
		default:
			e.logger.Error("seller has encountered an error receiving refund via escrow wallet")
			nknClient.Close()
			return fmt.Errorf("seller has encountered an error receiving refund via escrow wallet")
		}
	case <-ctx.Done():
		e.logger.Error("seller has not responded to delivery of refund via PK: " + order.OrderID)
		nknClient.Close()
		return fmt.Errorf("seller has not responded to refund via delivery of new PK for order: " + order.OrderID)
	}
}

// refundBuyerViaEscrowWallet is called to refund the buyer via the escrow wallet
func (e *ExchangeServer) refundBuyerViaEscrowWallet(order CompletedOrder) error {
	account, err := nkn.NewAccount(nil)
	if err != nil {
		e.logger.Error("creating a new nkn acount failed")
		return err
	}

	// initialize the nkn client.
	nknClient, err := nkn.NewMultiClient(account, "", 4, false, nil)
	if err != nil {
		e.logger.Error("creating a new nkn client failed")
		return err
	}

	ac := &NKNNotification{
		PrivateKey: order.BuyerEscrowWallet.PrivateKey,
		Chain:      order.TradeAsset,
	}

	bytes, err := json.Marshal(ac)
	if err != nil {
		e.logger.Error("error marshalling  refund buyer via escrow wallet: " + err.Error())
		return err
	}

	defer nknClient.Close()
	e.logger.Info("refunding buyer via NKN by sending the private key of the buyer's escrow wallet")

	e.logger.Infof("sending private key: %s to buyer: %s", ac.PrivateKey, order.BuyerNKNAddress)

	<-nknClient.OnConnect.C
	onReply, err := nknClient.Send(nkn.NewStringArray(order.BuyerNKNAddress), bytes, nil)
	if err != nil {
		e.logger.Error("refunding buyer via NKN by sending the private key of the buyer's escrow wallet failed")
		return err
	}

	// create a timeout of 2 minutes
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

	defer cancel()

	select {
	case reply := <-onReply.C:
		switch string(reply.Data) {
		case "ok":
			e.logger.Info("buyer has recieved escrow wallet refund")
			nknClient.Close()
			return nil
		default:
			e.logger.Error("buyer has encountered an error reciving  escrow wallet refund")
			nknClient.Close()
			return fmt.Errorf("buyer has encountered an error reciving  escrow wallet refund")
		}
	case <-ctx.Done():
		e.logger.Error("buyer has not responded to refund delivery of the PK: " + order.OrderID)
		nknClient.Close()
		return fmt.Errorf("buyer has not responded to refund notification of new PK for order: " + order.OrderID)
	}
}
