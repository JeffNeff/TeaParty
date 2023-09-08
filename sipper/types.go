package main

import (
	"crypto/ecdsa"
	"math/big"

	redis "github.com/go-redis/redis/v9"
)

type Sipper struct {
	teaServer   string
	teaBarrel   string
	adamsServer string
	brokerAddr  string
	moPK        *ecdsa.PrivateKey
	ethPK       *ecdsa.PrivateKey
	redisClient *redis.Client

	// chanel for websocket connection updates
	// teaWSConn chan *websocket.Conn
}

type EscrowWallet struct {
	PublicAddress string            `json:"publicAddress"`
	PrivateKey    string            `json:"privateKey"`
	Chain         string            `json:"chain"`
	ECDSA         *ecdsa.PrivateKey `json:"ecdsa"`
}

// CompletedOrder contains all the required elements to complete an order
type CompletedOrder struct {
	// BuyerEscrowWallet the escrow wallet that the buyer will be inserting the
	// TradeAsset into.
	BuyerEscrowWallet EscrowWallet `json:"buyerEscrowWallet"`
	// SellerEscrowWallet the escrow wallet that the seller will be inserting the
	// Currency into.
	SellerEscrowWallet EscrowWallet `json:"sellerEscrowWallet"`
	// SellerPaymentComplete is a boolean that tells us if the seller has completed
	// the payment.
	SellerPaymentComplete bool `json:"sellerPaymentComplete"`
	// BuyerPaymentComplete is a boolean that tells us if the buyer has completed
	// the payment.
	BuyerPaymentComplete bool `json:"buyerPaymentComplete"`
	// Amount the amount of funds that we are sending to the buyer.
	Amount *big.Int `json:"amount"`
	// OrderID the orderID that we are completing.
	OrderID string `json:"orderID"`
	// BuyerShippingAddress the public key of the account the buyer wants to receive on
	BuyerShippingAddress string `json:"buyerShippingAddress"`
	// BuyerRefundAddress
	BuyerRefundAddress string `json:"buyerRefundAddress"`
	// SellerRefundAddress
	SellerRefundAddress string `json:"sellerRefundAddress"`
	// SellerShippingAddress the public key of the account the seller wants to receive on
	SellerShippingAddress string `json:"sellerShippingAddress"`
	// BuyerNKNAddress the public NKN address of the buyer.
	BuyerNKNAddress string `json:"buyerNKNAddress"`
	// SellerNKNAddress the public NKN address of the seller.
	SellerNKNAddress string `json:"sellerNKNAddress"`
	// TradeAsset is the asset that we are sending to the buyer.
	TradeAsset string `json:"tradeAsset"`
	// Currency the currency that we are sending to the seller.
	Currency string `json:"currency"`
	// Price the price of the trade. (how much of the TradeAsset we are asking
	// from the seller for the Currency)
	Price *big.Int `json:"price"`
	// Timeout the amount of time that we are willing to wait for the transaction to be mined.
	Timeout int64 `json:"timeout"`
	// Stage reflects the stage of the order.
	Stage int `json:"stage"`
}

// BuyOrder is a struct that contains the information expected in a buy order
type BuyOrder struct {
	TXID string `json:"txid"`
	// BuyerShippingAddress represents the public key of the account the buyer wants to receive on
	BuyerShippingAddress string `json:"buyerShippingAddress"`
	// BuyerNKNAddress reflects the  publicly address of the buyer.
	BuyerNKNAddress string `json:"buyerNKNAddress"`
	// PaymentTransactionID reflects the transaction ID of the payment made in MO.
	PaymentTransactionID string `json:"paymentTransactionID"`
	// RefundAddress reflects the address of which the funds will be refunded in case of a failure.
	RefundAddress string `json:"refundAddress"`
}

// SellOrder contains the information expected in a sell order.
type SellOrder struct {
	// TradeAsset reflects the asset that the SELLER wishes to obtain. (bitcoin, mineonlium, USDT, etc).
	TradeAsset string `json:"tradeAsset"`
	// Price reflects the ammount of TradeAsset the SELLER requires.
	Price *big.Int `json:"price"`
	// Currency reflects the currency that the SELLER wishes to trade. (bitcoin, mineonlium, USDT, etc).
	Currency string `json:"currency"`
	// Amount reflects the ammount of Currency the SELLER wishes to trade.
	Amount *big.Int `json:"amount"`
	// TXID reflects the Transaction ID of the SELL order to be created.
	TXID string `json:"txid"`
	// Locked tells us if this transaction is pending/proccessing another payment.
	Locked bool `json:"locked" default:false`
	// SellerShippingAddress reflects the public key of the account the seller wants to receive on
	SellerShippingAddress string `json:"sellerShippingAddress"`
	// SellerNKNAddress reflects the  public NKN address of the seller.
	SellerNKNAddress string `json:"sellerNKNAddress"`
	// PaymentTransactionID reflects the transaction ID of the payment made in MO.
	PaymentTransactionID string `json:"paymentTransactionID"`
	// RefundAddress reflects the address of which the funds will be refunded in case of a failure.
	RefundAddress string `json:"refundAddress"`
	// Private reflects if the trade order is to be private or not. I.E. listed in the public
	// market place or not.
	Private bool `json:"private"`
}
