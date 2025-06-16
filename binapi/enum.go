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
