package binapi

type Category string

const (
	Spot       Category = "SPOT"
	UsdtMargin Category = "USDT-M"
	CoinMargin Category = "COIN-M"
)

type TimeInForce string

const (
	GoodTillCancel    TimeInForce = "GTC"
	ImmediateOrCancel TimeInForce = "IOC"
	FillOrKill        TimeInForce = "FOK"
	GoodTillCrossing  TimeInForce = "GTX"
	GoodTillDate      TimeInForce = "GTD"
)

type OrderType string

const (
	Limit              OrderType = "LIMIT"
	Market             OrderType = "MARKET"
	Stop               OrderType = "STOP"
	StopMarket         OrderType = "STOP_MARKET"
	TakeProfit         OrderType = "TAKE_PROFIT"
	TakeProfitMarket   OrderType = "TAKE_PROFIT_MARKET"
	TrailingStopMarket OrderType = "TRAILING_STOP_MARKET"
)

type Bar string

const (
	Bar1m  Bar = "1m"
	Bar3m  Bar = "3m"
	Bar5m  Bar = "5m"
	Bar15m Bar = "15m"
	Bar30m Bar = "30m"

	Bar1h  Bar = "1h"
	Bar2h  Bar = "2h"
	Bar4h  Bar = "4h"
	Bar6h  Bar = "6h"
	Bar8h  Bar = "8h"
	Bar12h Bar = "12h"

	Bar1d Bar = "1d"
	Bar3d Bar = "3d"

	Bar1w Bar = "1w"

	Bar1M Bar = "1M"
)
