package be

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var sellRequest = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "sell_request_total",
	Help: "The total number of sell requests",
}, []string{"currency", "amount", "trade_asset", "price", "on_chain", "private"})

var failedSellRequest = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "failed_sell_request_total",
	Help: "The total number of failed sell requests",
}, []string{"currency", "amount", "trade_asset", "price", "on_chain", "private", "reason"})

var buyRequest = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "buy_request_total",
	Help: "The total number of buy requests",
}, []string{"currency", "amount", "trade_asset", "price", "on_chain", "private"})

var failedBuyRequest = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "failed_buy_request_total",
	Help: "The total number of failed buy requests",
}, []string{"currency", "amount", "trade_asset", "price", "on_chain", "private", "reason"})

var completedTrades = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "completed_trades_total",
	Help: "The total number of completed trades",
}, []string{"currency", "amount", "trade_asset", "price", "on_chain", "private"})

var failedTrades = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "failed_trades_total",
	Help: "The total number of failed trades",
}, []string{"currency", "amount", "trade_asset", "price", "on_chain", "private", "reason"})

var failedAccountWatchRequests = promauto.NewCounter(prometheus.CounterOpts{
	Name: "failed_account_watch_requests_total",
	Help: "The total number of failed account watch requests",
})

var warrenWatchersInUse = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "warren_watchers_in_use",
	Help: "The number of warren watchers in use",
})

func MetricsSellRequestIncrement() {
	sellRequest.WithLabelValues("currency", "amount", "trade_asset", "price", "on_chain", "private").Add(1)
}

func MetricsFailedSellRequestIncrement(reason string) {
	failedSellRequest.WithLabelValues("currency", "amount", "trade_asset", "price", "on_chain", "private", reason).Add(1)
}

func MetricsBuyRequestIncrement() {
	buyRequest.WithLabelValues("currency", "amount", "trade_asset", "price", "on_chain", "private").Add(1)
}

func MetricsFailedBuyRequestIncrement(reason string) {
	failedBuyRequest.WithLabelValues("currency", "amount", "trade_asset", "price", "on_chain", "private", reason).Add(1)
}

func MetricsCompletedTradesIncrement() {
	completedTrades.WithLabelValues("currency", "amount", "trade_asset", "price", "on_chain", "private").Add(1)
}

func MetricsFailedTradesIncrement(reason string) {
	failedTrades.WithLabelValues("currency", "amount", "trade_asset", "price", "on_chain", "private", reason).Add(1)
}

func MetricsAddTradeInProgress() {
	tradesInProgress.WithLabelValues("currency", "amount", "trade_asset", "price", "on_chain", "private").Add(1)
}

func MetricsRemoveTradeInProgress() {
	tradesInProgress.WithLabelValues("currency", "amount", "trade_asset", "price", "on_chain", "private").Dec()
}

func MetricsAddWarrenWatcher() {
	warrenWatchersInUse.Add(1)
}

func MetricsRemoveWarrenWatcher() {
	warrenWatchersInUse.Dec()
}

func MetricsFailedAccountWatchRequestIncrement() {
	failedAccountWatchRequests.Add(1)
}
