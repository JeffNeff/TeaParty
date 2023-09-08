package be

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	btcRPC "github.com/btcsuite/btcd/rpcclient"
	cloudevents "github.com/cloudevents/sdk-go/v2"

	"github.com/ethereum/go-ethereum/ethclient"

	// party "github.com/teapartycrypto/TeaParty/adams/pkg/contract"

	solRPC "github.com/gagliardetto/solana-go/rpc"

	"github.com/go-redis/redis/v9"
	pkgadapter "knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"
)

// EnvAccessorCtor for configuration parameters
func EnvAccessorCtor() pkgadapter.EnvConfigAccessor {
	return &envAccessor{}
}

var _ pkgadapter.Adapter = (*ExchangeServer)(nil)

// NewAdapter adapter implementation
func NewAdapter(ctx context.Context, envAcc pkgadapter.EnvConfigAccessor, ceClient cloudevents.Client) pkgadapter.Adapter {
	env := envAcc.(*envAccessor)
	e := &ExchangeServer{}
	e.logger = logging.FromContext(ctx)

	exchangeServerID := env.ExchangeServerID

	// initialize the Party Chain nodes.
	partyclient, err := ethclient.Dial(env.PartyChainRPC1)
	if err != nil {
		if !env.Development {
			panic(err)
		}
	}
	partyclientTwo, err := ethclient.Dial(env.PartyChainRPC2)
	if err != nil {
		if !env.Development {
			panic(err)
		}
	}

	// initialize the ethereum nodes.
	ethClient1, err := ethclient.Dial(env.ETHRPC1)
	if err != nil {
		e.logger.Errorw("damm son no connection to eth rpc 1 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}
	ethClient2, err := ethclient.Dial(env.ETHRPC2)
	if err != nil {
		e.logger.Errorw("damm son no connection to eth rpc 2 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	etcClient, ett := ethclient.Dial(env.ETCRPC1)
	if ett != nil {
		e.logger.Errorw("damm son no connection to etc rpc 1 ")
		if !env.Development {
			panic(ett)
			return nil
		}
	}

	etcClient2, ett := ethclient.Dial(env.ETCRPC2)
	if ett != nil {
		e.logger.Errorw("damm son no connection to etc rpc 2 ")
		if !env.Development {
			panic(ett)
			return nil
		}
	}

	ethoClient, err := ethclient.Dial(env.ETHORPC1)
	if err != nil {
		e.logger.Errorw("damm son no connection to etho rpc 1 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	ethoClient2, err := ethclient.Dial(env.ETHORPC2)
	if err != nil {
		e.logger.Errorw("damm son no connection to etho rpc 2 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	acc, err := ethclient.Dial(env.AltCoinRPC1)
	if err != nil {
		e.logger.Errorw("damm son no connection to AltCoinCash rpc 1 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	acc2, err := ethclient.Dial(env.AltCoinRPC2)
	if err != nil {
		e.logger.Errorw("damm son no connection to AltCoinCash rpc 2 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	flo, err := ethclient.Dial(env.FloraRPC1)
	if err != nil {
		e.logger.Errorw("damm son no connection to Flo rpc 1 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	flo2, err := ethclient.Dial(env.FloraRPC2)
	if err != nil {
		e.logger.Errorw("damm son no connection to Flo rpc 2 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	// initialize the polygon nodes.
	polyclient, err := ethclient.Dial(env.POLYRPC1)
	if err != nil {
		e.logger.Errorw("damm son no connection to polygon rpc 1 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}
	polyclientTwo, err := ethclient.Dial(env.POLYRPC2)
	if err != nil {
		e.logger.Errorw("damm son no connection to polygon rpc 2 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	octClient, err := ethclient.Dial(env.OCTRPC1)
	if err != nil {
		e.logger.Errorw("damm son no connection to Octa rpc 1 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	octClient2, err := ethclient.Dial(env.OCTRPC2)
	if err != nil {
		e.logger.Errorw("damm son no connection to Octa rpc 2 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	// initalize the ltc nodes
	var connCfgLTC = &btcRPC.ConnConfig{
		Host:         env.LTCRPC1,
		User:         "dockeruser",
		Pass:         "dockerpass",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	ltc1, err := btcRPC.New(connCfgLTC, nil)
	if err != nil {
		e.logger.Errorw("holy fucknuts batman!! The RPC for ltc is down!!")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	var connCfgLTC2 = &btcRPC.ConnConfig{
		Host:         env.LTCRPC2,
		User:         "dockeruser",
		Pass:         "dockerpass",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	ltc2, err := btcRPC.New(connCfgLTC2, nil)
	if err != nil {
		e.logger.Errorw("holy fucknuts batman!! The RPC for ltc is down!!")
		if !env.Development {
			panic(err)
			return nil
		}
	}
	// initialize the radiant nodes.
	var connCfg = &btcRPC.ConnConfig{
		Host:         env.RXDRPC1,
		User:         "dockeruser",
		Pass:         "dockerpass",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	rxd1, err := btcRPC.New(connCfg, nil)
	if err != nil {
		e.logger.Errorw("holy fucknuts batman!! The RPC for Radiant is down!!")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	var connCfg2 = &btcRPC.ConnConfig{
		Host:         env.RXDRPC2,
		User:         "dockeruser",
		Pass:         "dockerpass",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	rxd2, err := btcRPC.New(connCfg2, nil)
	if err != nil {
		e.logger.Errorw("holy fucknuts batman!! The Second RPC for Radiant is down!!")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	celo, err := ethclient.Dial(env.CELORPC1)
	if err != nil {
		e.logger.Errorw("damm son no connection to celo rpc 1 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}
	celo2, err := ethclient.Dial(env.CELORPC2)
	if err != nil {
		e.logger.Errorw("damm son no connection to celo rpc 2 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	canto, err := ethclient.Dial(env.CANTORPC1)
	if err != nil {
		e.logger.Errorw("damm son no connection to canto rpc 1 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}
	canto2, err := ethclient.Dial(env.CANTORPC2)
	if err != nil {
		e.logger.Errorw("damm son no connection to canto rpc 2 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	// initialize the litecoin nodes.
	var btcConnCfg = &btcRPC.ConnConfig{
		Host:         env.BTCRPC1,
		User:         "dockeruser",
		Pass:         "dockerpass",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	btc1, err := btcRPC.New(btcConnCfg, nil)
	if err != nil {
		e.logger.Errorw("holy fucknuts batman!! The RPC for BTC is down!!")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	var btcConnCfg2 = &btcRPC.ConnConfig{
		Host:         env.BTCRPC2,
		User:         "dockeruser",
		Pass:         "dockerpass",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	btc2, err := btcRPC.New(btcConnCfg2, nil)
	if err != nil {
		e.logger.Errorw("holy fucknuts batman!! The Second RPC for BTC is down!!")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	// initialize the solana nodes.
	solClient := solRPC.New(env.SOLRPC1)

	// test connection
	_, err = solClient.GetRecentBlockhash(context.Background(), solRPC.CommitmentRecent)
	if err != nil {
		e.logger.Errorw("damm son no connection to solana rpc 1 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	solClient2 := solRPC.New(env.SOLRPC2)

	// test connection
	_, err = solClient2.GetRecentBlockhash(context.Background(), solRPC.CommitmentRecent)
	if err != nil {
		e.logger.Errorw("damm son no connection to solana rpc 1 ")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	e.celoNode.rpcClient = celo
	e.celoNode.rpcClientTwo = celo2
	e.partyChain.rpcClient = partyclient
	e.partyChain.rpcClientTwo = partyclientTwo
	e.ethNode.rpcClient = ethClient1
	e.ethNode.rpcClientTwo = ethClient2
	e.polygonNode.rpcClient = polyclient
	e.polygonNode.rpcClientTwo = polyclientTwo
	// e.kaspaNode.rpcClient = kaspaClient
	// e.kaspaNode.rpcClientTwo = kaspaClient2
	e.radiantNode.rpcConfig = connCfg
	e.radiantNode.rpcConfigTwo = connCfg2
	e.radiantNode.rpcClient = rxd1
	e.radiantNode.rpcClientTwo = rxd2
	e.btcNode.rpcConfig = btcConnCfg
	e.btcNode.rpcConfigTwo = btcConnCfg2
	e.btcNode.rpcClient = btc1
	e.btcNode.rpcClientTwo = btc2
	e.solNode.rpcClient = solClient
	e.solNode.rpcClientTwo = solClient2
	e.octNode.rpcClient = octClient
	e.octNode.rpcClientTwo = octClient2
	e.floNode.rpcClient = flo
	e.floNode.rpcClientTwo = flo2
	e.altcoinchain.rpcClient = acc
	e.altcoinchain.rpcClientTwo = acc2
	e.cantoNode.rpcClient = canto
	e.cantoNode.rpcClientTwo = canto2
	e.etcNode.rpcClient = etcClient
	e.etcNode.rpcClientTwo = etcClient2
	e.ethONode.rpcClient = ethoClient
	e.ethONode.rpcClientTwo = ethoClient2
	e.ltcNode.rpcClient = ltc1
	e.ltcNode.rpcClientTwo = ltc2
	e.ltcNode.rpcConfig = connCfgLTC
	e.ltcNode.rpcConfigTwo = connCfgLTC2

	e.watch = env.Watch
	e.dev = env.Development
	e.ceClient = ceClient
	e.blockExplorer = env.PartyChainBlockExplorer
	e.exchangeServerID = exchangeServerID
	// addressesTimedOut is a type of map[string]time.Time
	// populate it with an inital address of 192.168.1.1 and a time of 0
	e.addressesTimedOut = make(map[string]time.Time)
	e.addressesTimedOut["192.168.0.1"] = time.Time{}
	// initialize the redis client.
	e.redisClient = redis.NewClient(&redis.Options{
		Addr:     env.RedisAddress,
		Password: env.RedisPassword,
		DB:       env.RedisDB,
	})

	// test the redis connection.
	_, err = e.redisClient.Ping(ctx).Result()
	if err != nil {
		e.logger.Errorw("damm son no connection to redis ")
		if !env.Development {
			panic(err)
			return nil
		}
	}

	cofd, err := e.fetchOrdersFromDB()
	if err != nil {
		// create the orders map.
		cofd = []SellOrder{}
		// save the orders map to redis.
		orders, err := json.Marshal(cofd)
		if err != nil {
			e.logger.Errorw("failed to marshal the orders map")
			panic(err)
			return nil
		}
		err = e.redisClient.Set(ctx, "orders", orders, 0).Err()
		if err != nil {
			e.logger.Errorw("failed to save the orders map to redis")
			panic(err)
			return nil
		}
	}

	co, err := e.fetchCompleteOrdersFromDB()
	if err != nil {
		// create the orders map.
		co = []CompletedOrder{}
		// save the orders map to redis.
		co, err := json.Marshal(co)
		if err != nil {
			e.logger.Errorw("failed to marshal the completedOrders map")
			panic(err)
			return nil
		}
		err = e.redisClient.Set(ctx, "completeorders", co, 0).Err()
		if err != nil {
			e.logger.Errorw("failed to save the completedOrders map to redis")
			panic(err)
			return nil
		}
	}
	fmt.Println("cofd", cofd)
	fmt.Println("co", co)

	// privateKey, err := crypto.HexToECDSA(env.SMARTCONTRACTPRIVATEKEY)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	log.Fatal("error casting public key to ECDSA")
	// }

	// fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// nonce, err := e.partyChain.rpcClient.PendingNonceAt(context.Background(), fromAddress)
	// if err != nil {
	// 	fmt.Println("failed to get nonce: %v", err)
	// 	log.Fatal(err)
	// }

	// gasPrice, err := e.partyChain.rpcClient.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	fmt.Printf("failed to get suggest gas price: %v", err)
	// 	log.Fatal(err)
	// }

	// address := common.HexToAddress("0xDAD416F84E67d4B37c12c37979DC4E8d07Fc83d2")
	// instance, err := party.NewTeaParty(address, e.partyChain.rpcClient)
	// if err != nil {
	// 	fmt.Printf("failed to instantiate a Token contract: %v", err)
	// 	log.Fatal(err)
	// }

	// chainID, err := partyclient.ChainID(context.Background())
	// if err != nil {
	// 	fmt.Printf("failed to get chainID: %v", err)
	// 	if !env.Development {
	// 		panic(err)
	// 	}
	// }

	// auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	// if err != nil {
	// 	log.Fatalf("Failed to create authorized transactor: %v", err)
	// }

	// auth.Nonce = big.NewInt(int64(nonce))
	// auth.Value = big.NewInt(0)     // in wei
	// auth.GasLimit = uint64(300000) // in units
	// auth.GasPrice = gasPrice

	// e.moContractTransactOpts = auth
	// e.partyContract = instance
	e.noPartyFeeAddresses = env.NOPARTYFEEADDRESSES

	return e
}

func createChannel() (chan os.Signal, func()) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	return stopCh, func() {
		close(stopCh)
	}
}

func (e *ExchangeServer) Start(ctx context.Context) error {
	e.logger.Info("starting warren...")
	go e.StartWarren(ctx)
	e.logger.Info("started warren")
	e.logger.Info("starting http server...")
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
	go e.StartHttpServer(ctx, server)

	go e.verificationWorker()

	stopCh, stop := createChannel()
	defer stop()
	<-stopCh

	e.warrenWG.Wait()

	e.logger.Info("stopping Adams")
	e.shutdown(ctx, server)
	return nil
}

// StartWarren starts the warren account watching service
func (e *ExchangeServer) StartWarren(ctx context.Context) error {
	// log the pod name
	e.logger.Infof("starting warren on pod %s", os.Getenv("POD_NAME"))
	// start the account watch service.
	e.ctx = ctx
	e.warrenWG = &sync.WaitGroup{}
	numWorkers := runtime.NumCPU()
	e.warrenWG.Add(numWorkers)
	e.warrenChan = make(chan AccountWatchRequest, numWorkers)

	for i := 0; i < numWorkers; i++ {
		go e.warrenWorker()
	}

	// create a timer that ticks every 30 seconds.
	// create a ticker that ticks every 30 seconds
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	if e.dev {
		ticker = time.NewTicker(time.Second * 10)
	}

	canIlive := true
	for canIlive {
		select {
		case <-ticker.C:
			e.logger.Info("clearing the cawr and updating it with the awr from the database")
			// clear the cawr
			cawr := make([]AccountWatchRequest, 0)
			// update the account watch requests
			// get the account watch requests from the database.
			awr, err := e.retrieveAccountWatchRequestsFromDB()
			if err != nil {
				e.logger.Errorw("error retrieving account watch requests from database", "error", err)
				// do not return the error, continue running the loop
				continue
			}
			if awr != nil {
				if len(awr) > 0 {
					cawr = append(cawr, awr...)
				}
			}
			// send the account watch requests to the warren service.
			e.Warren(cawr)
		case <-ctx.Done():
			// context is canceled, stop the loop
			e.logger.Info("context is canceled, stopping the warren loop")
			canIlive = false
		}
	}
	return nil
}

func (e *ExchangeServer) Warren(awr []AccountWatchRequest) {
	if len(awr) == 0 || awr == nil {
		e.logger.Info("no account watch requests to watch")
		return
	}

	// illerate over the awr's, verify that the watch is not locked, and start the watch.
	// if the watch is locked, verify that the watch is not older then 2.2 hours old.
	// if it is older then 2.2 hours old, then restart the watch.
	for _, request := range awr {
		if request.Locked {
			// check if the watch is older then 2.2 hours old.
			// if it is, then restart the watch.
			// if time.Since(request.LockedTime) > time.Hour*2 {
			// 	e.logger.Infof("watch for account %s is older then 2.2 hours, restarting the watch", request.Account)
			// 	// restart the watch.
			// 	request.Locked = true
			// 	request.LockedTime = time.Now()
			// 	request.LockedBy = os.Getenv("POD_NAME")
			// 	e.warrenChan <- request
			// }
			continue
		}
		e.warrenChan <- request
	}
}

func (e *ExchangeServer) warrenWorker() {
	defer e.warrenWG.Done()

	for {
		select {
		case request := <-e.warrenChan:
			if !request.Locked {
				e.logger.Infof("starting watch for account %s on chain %s", request.Account, request.Chain)
				// add the watch to the metrics counter
				MetricsAddWarrenWatcher()
				// start the watch.
				go e.watchAccount(&request)
			}
		case <-e.ctx.Done():
			// remove the watch from the metrics counter
			MetricsRemoveWarrenWatcher()
			// Exit the loop when the context is canceled
			return
		}
	}
}

func (e *ExchangeServer) StartHttpServer(ctx context.Context, server *http.Server) {
	e.logger.Info("Starting Adams")

	// Metrics
	http.Handle("/metrics", promhttp.Handler())

	// REVIEW
	http.HandleFunc("/sell", e.Sell)
	http.HandleFunc("/buy", e.Buy)
	http.HandleFunc("/listorders", e.FetchSellOrders)

	// Experimental
	http.HandleFunc("/closeMarketOrder", e.CloseOpenMarketOrder)

	http.HandleFunc("/assistedSell", e.AssistedSell)

	http.HandleFunc("/cancleAssistedSell", e.CancleAssistedSell)

	e.logger.Info("starting server on port :8080")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.logger.Errorw("failed to start server", "error", err)
		panic(err)
	} else {
		e.logger.Info("application stopped gracefully")
	}
	<-ctx.Done()
	e.logger.Info("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(shutdownCtx)
}

// verificationWorker is a worker that is responsible for verifying the orders in the marketplace.
// every 60 seconds the worker will run the verification process.
func (e *ExchangeServer) verificationWorker() {

	// create a timer that ticks every 60 seconds.
	// create a ticker that ticks every 60 seconds
	ticker := time.NewTicker(time.Second * 360)
	defer ticker.Stop()

	canIlive := true
	for canIlive {
		select {
		case <-ticker.C:
			e.logger.Info("verifying orders in the marketplace")
			// verify the orders in the marketplace.
			e.verifyOrdersInMarketplace()
		}
	}
}

func (e *ExchangeServer) shutdown(ctx context.Context, server *http.Server) {
	if err := e.updateAccountWatchRequestsOnCrash(); err != nil {
		e.logger.Errorw("error updating account watch requests on crash", "error", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	} else {
		log.Println("application shutdown")
	}
}

// verifyOrdersInMarketplace looks at all of the current avalible orders in the marketplace and pings the associated NKN address
// to verify if the order is still valid. If the order is not valid, then the order is removed from the marketplace.
func (e *ExchangeServer) verifyOrdersInMarketplace() {
	// get all of the orders in the marketplace.
	orders, err := e.fetchOrdersFromDB()
	if err != nil {
		e.logger.Errorw("error fetching orders from database", "error", err)
		return
	}

	if len(orders) == 0 {
		e.logger.Info("no orders in the marketplace")
		return
	}

	// illerate over the orders and verify that the order is still valid  by pinging the associated NKN address.
	for _, order := range orders {
		// ping the associated NKN address.
		// if the ping fails, then remove the order from the marketplace.
		if err := e.pingNKNAddress(order.SellerNKNAddress); err != nil {
			e.logger.Errorw("error pinging nkn address", "error", err)
			// remove the order from the marketplace.
			e.logger.Infof("removing order %s from the marketplace", order.TXID)

			newOrders := make([]SellOrder, 0)
			for _, o := range orders {
				if o.TXID != order.TXID {
					newOrders = append(newOrders, o)
				}
			}

			if err := e.updateOrdersInDB(newOrders); err != nil {
				e.logger.Errorw("error updating orders in database", "error", err)
				continue
			}

			// remove the order from the marketplace.
			e.logger.Infof("removing order %s from the marketplace", order.TXID)
			continue
		}
	}
}
