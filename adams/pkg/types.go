package be

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	// party "github.com/teapartycrypto/TeaParty/adams/pkg/contract"

	btcRPC "github.com/btcsuite/btcd/rpcclient"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/ethereum/go-ethereum/ethclient"
	solRPC "github.com/gagliardetto/solana-go/rpc"
	"github.com/go-redis/redis/v9"
	kasrpc "github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"go.uber.org/zap"
	pkgadapter "knative.dev/eventing/pkg/adapter/v2"
)

const (
	ACC    = "altcoinchain"
	FLO    = "flora"
	ETH    = "ethereum"
	ETC    = "ethereumclassic"
	ETHONE = "ethOne"
	GRAMS  = "grams"
	POL    = "polygon"
	KAS    = "kaspa"
	RXD    = "radiant"
	CEL    = "celo"
	SOL    = "solana"
	OCT    = "octa"
	CANTO  = "canto"
	LTC    = "litecoin"
	BTC    = "bitcoin"
	ETHO   = "etho"
	CFXE   = "confluxEspace"

	// tokens
	BSCUSDT = "bscusdt"
	WGRAMS  = "wgrams"

	// unsupported
	ALP  = "alephium"
	NEAR = "near"

	// NFT's
	MiningGame = "miningGame"
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

// ErrorEvent represents the expected information in an emitted error event
// "tea.party.error"| ERROREVENT
type ErrorEvent struct {
	Err     string
	Context string
	Data    interface{}
}

// AccountWatchRequest is the information we need to watch a new account
// this type is associated with the "tea.party.watch.account" | IOWATCHACCOUNTREQUEST event type
type AccountWatchRequest struct {
	Seller                       bool                          `json:"seller"`
	Account                      string                        `json:"account"`
	Chain                        string                        `json:"chain"`
	Amount                       *big.Int                      `json:"amount"`
	NFTID                        int64                         `json:"nft_id"`
	TransactionID                string                        `json:"transaction_id"`
	TimeOut                      int64                         `json:"timeout"`
	FinalizeOnChain              bool                          `json:"finalizeOnChain"`
	AssistedSellOrderInformation AssistedTradeOrderInformation `json:"assistedSellOrderInformation"`
	Locked                       bool                          `json:"locked"`
	LockedTime                   time.Time                     `json:"lockedTime"`
	LockedBy                     string                        `json:"lockedBy"`
	AWRID                        string                        `json:"awrid"`
}

// AccountWatchRequestResult is the result of the watch request
// this type is associated with the "tea.party.watch.result" | IOWATCHRESULT event type
type AccountWatchRequestResult struct {
	AccountWatchRequest AccountWatchRequest `json:"account_watch_request"`
	Result              string              `json:"result"`
}

type envAccessor struct {
	pkgadapter.EnvConfig
	Development bool `envconfig:"DEV" default:"false"`
	Watch       bool `envconfig:"WATCH" default:"true"`

	ExchangeServerID string `envconfig:"EXCHANGE_SERVER_ID" required:"true"`

	SMARTCONTRACTPRIVATEKEY string `envconfig:"PRIVATE_KEY" required:"true"`

	CELORPC1 string `envconfig:"CELO_RPC_1" default:"" required:"true"`
	CELORPC2 string `envconfig:"CELO_RPC_2" default:"" required:"true"`

	AltCoinRPC1 string `envconfig:"ALTCOIN_RPC_1" default:"https://rpc0.altcoinchain.org/rpc" required:"true"`
	AltCoinRPC2 string `envconfig:"ALTCOIN_RPC_2" default:"https://rpc0.altcoinchain.org/rpc" required:"true"`

	FloraRPC1 string `envconfig:"FLORA_RPC_1" default:"https://rpc.florascan.io" required:"true"`
	FloraRPC2 string `envconfig:"FLORA_RPC_2" default:"https://rpc.florascan.io" required:"true"`

	LTCRPC1 string `envconfig:"LTC_RPC_1" default:"" required:"true"`
	LTCRPC2 string `envconfig:"LTC_RPC_2" default:"" required:"true"`

	ETHORPC1 string `envconfig:"ETHORPC1" default:"https://rpc.ethoprotocol.com" required:"true"`
	ETHORPC2 string `envconfig:"ETHORPC2" default:"https://rpc.ethoprotocol.com" required:"true"`

	CFXEspaceRPC1 string `envconfig:"CFXEspaceRPC1" default:"https://confluxscan.net" required:"true"`
	CFXEspaceRPC2 string `envconfig:"CFXEspaceRPC2" default:"https://confluxscan.net" required:"true"`

	PartyChainRPC1 string `envconfig:"PARTY_CHAIN_1" required:"true"`
	PartyChainRPC2 string `envconfig:"PARTY_CHAIN_2" required:"true"`

	ETHRPC1 string `envconfig:"ETH_RPC_1" default:"" required:"true"`
	ETHRPC2 string `envconfig:"ETH_RPC_2" default:"" required:"true"`

	ETHONERPC1 string `envconfig:"ETHONE_RPC_1" default:"" required:"true"`
	ETHONERPC2 string `envconfig:"ETHONE_RPC_2" default:"" required:"true"`

	ETCRPC1 string `envconfig:"ETC_RPC_1" default:"https://www.ethercluster.com/etc" required:"true"`
	ETCRPC2 string `envconfig:"ETC_RPC_2" default:"https://www.ethercluster.com/etc" required:"true"`

	POLYRPC1 string `envconfig:"POLY_RPC_1" default:"" required:"true"`
	POLYRPC2 string `envconfig:"POLY_RPC_2" default:"" required:"true"`

	KASRPC1 string `envconfig:"KAS_RPC_1" default:"" `
	KASRPC2 string `envconfig:"KAS_RPC_2" default:"" `

	RXDRPC1 string `envconfig:"RXD_RPC_1" default:"" `
	RXDRPC2 string `envconfig:"RXD_RPC_2" default:"" `

	OCTRPC1 string `envconfig:"OCT_RPC_1" default:"" required:"true"`
	OCTRPC2 string `envconfig:"OCT_RPC_2" default:"" required:"true"`

	BTCRPC1 string `envconfig:"BTC_RPC_1" default:"" required:"true"`
	BTCRPC2 string `envconfig:"BTC_RPC_2" default:"" required:"true"`

	SOLRPC1 string `envconfig:"SOL_RPC_1" default:"" required:"true"`
	SOLRPC2 string `envconfig:"SOL_RPC_2" default:"" required:"true"`

	CANTORPC1 string `envconfig:"CANTORPC1" default:"https://canto.gravitychain.io" required:"true"`
	CANTORPC2 string `envconfig:"CANTORPC2" default:"https://canto.gravitychain.io" required:"true"`

	NOPARTYFEEADDRESSES []string `envconfig:"NO_PARTY_FEE_ADDRESSES" default:"0x2cc906ee4E8A648c1C5fFf5209155C4edF4D1A40"`

	PartyChainBlockExplorer string `envconfig:"BLOCK_EXPLORER" default:"" required:"true"`

	// redis server
	RedisAddress  string `envconfig:"REDIS_ADDRESS" required:"true"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`
	RedisDB       int    `envconfig:"REDIS_DB" default:"0"`
}

// // SellerNotification represents the information that is to be sent to the seller
// // once a buyer appears for a posted order
// type SellerNotification struct {
// 	Address string `json:"address"` // the address of the seller
// 	Amount  string `json:"amount"`  // the amount of the order
// 	Network string `json:"network"` // the network the order is on
// }

// ExchangeServer holds the state of the exchange server.
type ExchangeServer struct {
	ctx              context.Context
	exchangeServerID string
	btcNode          BTCNode
	ltcNode          BTCNode
	radiantNode      BTCNode
	altcoinchain     EthereumNode
	celoNode         EthereumNode
	cantoNode        EthereumNode
	floNode          EthereumNode
	ethNode          EthereumNode
	etcNode          EthereumNode
	ethOneNode       EthereumNode
	ethONode         EthereumNode
	cfxEspaceNode    EthereumNode
	kaspaNode        KaspaNode
	partyChain       EthereumNode
	// nearNode    EthereumNode
	octNode     EthereumNode
	polygonNode PolygonNode

	solNode SOLNode

	// ltcNode     BTCNode

	addressesTimedOut map[string]time.Time
	blockExplorer     string

	// nknClient is the client used to interact with the NKN network.
	// nknClient *nkn.MultiClient

	redisClient *redis.Client

	moContractTransactOpts *bind.TransactOpts
	// partyContract          *party.TeaParty
	noPartyFeeAddresses []string

	warrenChan chan AccountWatchRequest
	warrenWG   *sync.WaitGroup

	ceClient cloudevents.Client
	logger   *zap.SugaredLogger
	dev      bool
	watch    bool
}

// BTCNode hold all the information and interfaces we need to interact
// with the a bitcoin node.
type BTCNode struct {
	rpcClient    *btcRPC.Client
	rpcClientTwo *btcRPC.Client
	rpcConfig    *btcRPC.ConnConfig
	rpcConfigTwo *btcRPC.ConnConfig
}

type SOLNode struct {
	rpcClient    *solRPC.Client
	rpcClientTwo *solRPC.Client
}

type EthereumNode struct {
	rpcClient    *ethclient.Client
	rpcClientTwo *ethclient.Client
}

type PolygonNode struct {
	rpcClient    *ethclient.Client
	rpcClientTwo *ethclient.Client
}

type KaspaNode struct {
	rpcClient    *kasrpc.RPCClient
	rpcClientTwo *kasrpc.RPCClient
}

type AccountGenResponse struct {
	PrivateKey string `json:"privateKey"`
	PubKey     string `json:"publicKey"`
	Address    string `json:"address"`
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
	// TradeAsset reflects the asset the buyer elected to trade for (mineonlium, bitcoin, USDT, etc).
	// this is an optional field. only avalible when the seller lists "ANY" as the trade asset.
	TradeAsset string `json:"tradeAsset"`
	// OnChain reflects if the trade order is to be finalized on-chain or not.
	OnChain bool `json:"onChain"`
}

type CloseOpenMarketOrderRequest struct {
	OrderID string `json:"orderID"`
}

type CancleAssistedSellOrderRequest struct {
	OrderID          string `json:"orderID"`
	CancelationToken string `json:"cancelationToken"`
}

// SellOrder contains the information expected in a sell order.
type SellOrder struct {
	// Currency reflects the currency that the SELLER wishes to trade. (bitcoin, mineonlium, USDT, etc).
	Currency string `json:"currency"`
	// Amount reflects the ammount of Currency the SELLER wishes to trade.
	Amount *big.Int `json:"amount"`
	// TradeAsset reflects the asset that the SELLER wishes to obtain. (bitcoin, mineonlium, USDT, etc).
	TradeAsset string `json:"tradeAsset"`
	// Price reflects the ammount of TradeAsset the SELLER requires.
	Price *big.Int `json:"price"`
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
	// OnChain reflects if the trade order is to be finalized on-chain or not.
	OnChain bool `json:"onChain"`
	// Assisted reflects if the trade order is to be assisted by the exchange or not.
	Assisted bool `json:"assisted"`
	// AssistedTradeOrderInformation reflects the information required to assist the trade order.
	AssistedTradeOrderInformation AssistedTradeOrderInformation `json:"assistedTradeOrderInformation"`
	// NFTID reflects the NFT ID of the NFT that is being traded.
	NFTID int64 `json:"nftID"`
}

type AssistedTradeOrderInformation struct {
	// SellersEscrowWallet represents the wallet that the seller has already funded
	// with the currency they wish to trade.
	SellersEscrowWallet   EscrowWallet `json:"sellersEscrowWallet"`
	SellerRefundAddress   string       `json:"sellerRefundAddress"`
	SellerShippingAddress string       `json:"sellerShippingAddress"`
	// TradeAsset reflects the asset that the SELLER wishes to obtain. (bitcoin, mineonlium, USDT, etc). Or Bridge.
	TradeAsset string `json:"tradeAsset"`
	// Price reflects the ammount of TradeAsset the SELLER requires.
	Price *big.Int `json:"price"`
	// Currency reflects the currency that the SELLER wishes to trade. (bitcoin, mineonlium, USDT, etc).
	Currency string `json:"currency"`
	// Amount reflects the ammount of Currency the SELLER wishes to trade.
	Amount *big.Int `json:"amount"`
	// BridgeTo reflects the blockchain that the seller wishes to bridge to.
	BridgeTo string `json:"bridgeTo"`
}

// CompletedTransactionInformation represents data expected when
// describing a transaction that has been completed on-chain.
type CompletedTransactionInformation struct {
	// the transaction id of the completed transaction
	TXID string `json:"txid"`
	// the amount of the transaction
	Amount *big.Int `json:"amount"`
	// the blockchain the transaction was completed on
	Blockchain string `json:"blockchain"`
}

type EscrowWallet struct {
	PublicAddress string `json:"publicAddress"`
	PrivateKey    string `json:"privateKey"`
	Chain         string `json:"chain"`
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
	// BuyerToFinalizeOnChain is a boolean that tells us if the buyer has elected
	// to finalize the transaction on-chain.
	BuyerToFinalizeOnChain bool `json:"buyerToFinalizeOnChain"`
	// SellerRefundAddress
	SellerRefundAddress string `json:"sellerRefundAddress"`
	// SellerShippingAddress the public key of the account the seller wants to receive on
	SellerShippingAddress string `json:"sellerShippingAddress"`
	// SellerToFinalizeOnChain is a boolean that tells us if the seller has elected
	// to finalize the transaction on-chain.
	SellerToFinalizeOnChain bool `json:"sellerToFinalizeOnChain"`
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
	//FailureReason reflects the reason for the failure of the order, if any.
	FailureReason string `json:"failureReason"`
	// Assisted reflects if the trade order is to be assisted by the exchange or not.
	Assisted bool `json:"assisted"`
	// AssistedTradeOrderInformation reflects the information required to assist the trade order.
	AssistedTradeOrderInformation *AssistedTradeOrderInformation `json:"assistedTradeOrderInformation"`
	// NFTID reflects the NFT ID of the NFT that is being traded.
	NFTID int64 `json:"nftID"`
}

// Query contains the information expected in a transaction query
type Query struct {
	TXID string `json:"txid"`
}

type BlockScoutTxQueryResponse struct {
	Message string `json:"message"`
	Result  struct {
		BlockNumber    string        `json:"blockNumber"`
		Confirmations  string        `json:"confirmations"`
		From           string        `json:"from"`
		GasLimit       string        `json:"gasLimit"`
		GasPrice       string        `json:"gasPrice"`
		GasUsed        string        `json:"gasUsed"`
		Hash           string        `json:"hash"`
		Input          string        `json:"input"`
		Logs           []interface{} `json:"logs"`
		NextPageParams interface{}   `json:"next_page_params"`
		RevertReason   string        `json:"revertReason"`
		Success        bool          `json:"success"`
		TimeStamp      string        `json:"timeStamp"`
		To             string        `json:"to"`
		Value          string        `json:"value"`
	} `json:"result"`
	Status string `json:"status"`
}

type AssetTradeData struct {
	Asset             string      `json:"asset"`
	NumberOfTrades    int         `json:"numberOfTrades"`
	TotalAmountTraded *big.Int    `json:"totalAmountTraded"`
	TradesForAsset    []TradeData `json:"tradesForAsset"`
}

type TradeData struct {
	TradeID       string   `json:"tradeID"`
	TradeAsset    string   `json:"tradeAsset"`
	TradeAmount   *big.Int `json:"tradeAmount"`
	TradeCurrency string   `json:"tradeCurrency"`
	TradePrice    *big.Int `json:"tradePrice"`
	TradeTime     int64    `json:"tradeTime"`
}

type Metrics struct {
	OpenOrders                    int              `json:"openOrders"`
	CompletedOrders               int              `json:"completedOrders"`
	TotalNumberOfTrades           int              `json:"totalNumberOfTrades"`
	TotalVolume                   *big.Int         `json:"totalVolume"`
	TotalNumberOfAssetsTraded     int              `json:"totalNumberOfAssetsTraded"`
	TotalNumberOfCurrenciesTraded int              `json:"totalNumberOfCurrenciesTraded"`
	AssetsTraded                  []AssetTradeData `json:"assetsTraded"`
	TotalFees                     *big.Int         `json:"totalFees"`
	FailedOrders                  int              `json:"failedOrders"`
	CancledOrders                 int              `json:"cancledOrders"`
	AccountsUnderWatch            int              `json:"accountsUnderWatch"`
}

type BridgeRequest struct {
	// Currency reflects the currency that the SELLER wishes to trade. (bitcoin, mineonlium, USDT, etc).
	Currency string `json:"currency"`
	// Amount reflects the ammount of Currency the SELLER wishes to trade.
	Amount          *big.Int `json:"amount"`
	FromChain       string   `json:"fromChain"`
	BridgeTo        string   `json:"bridgeTo"`
	ShippingAddress string   `json:"shippingAddress"`
}
