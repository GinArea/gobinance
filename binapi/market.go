package binapi

import (
	"fmt"
	"strconv"

	"github.com/msw-x/moon/ujson"
)

type ExchangeInfo struct {
	Timezone        string          `json:"timezone"`
	ServerTime      ujson.Int64     `json:"serverTime"`
	FuturesType     string          `json:"futuresType"`
	RateLimits      []RateLimitResp `json:"rateLimits"`
	ExchangeFilters []interface{}   `json:"exchangeFilters"`
	Assets          []Asset         `json:"assets"`
	Symbols         []Symbol        `json:"symbols"`
}

type RateLimitResp struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

type Asset struct {
	Asset             string        `json:"asset"`
	MarginAvailable   bool          `json:"marginAvailable"`
	AutoAssetExchange ujson.Float64 `json:"autoAssetExchange"`
}

type Symbol struct {
	Symbol                string        `json:"symbol"`
	Pair                  string        `json:"pair"`
	ContractType          string        `json:"contractType"`
	DeliveryDate          int64         `json:"deliveryDate"`
	OnboardDate           int64         `json:"onboardDate"`
	MaintMarginPercent    ujson.Float64 `json:"maintMarginPercent"`
	RequiredMarginPercent ujson.Float64 `json:"requiredMarginPercent"`
	BaseAsset             string        `json:"baseAsset"`
	QuoteAsset            string        `json:"quoteAsset"`
	MarginAsset           string        `json:"marginAsset"`
	PricePrecision        int           `json:"pricePrecision"`
	QuantityPrecision     int           `json:"quantityPrecision"`
	BaseAssetPrecision    int           `json:"baseAssetPrecision"`
	QuotePrecision        int           `json:"quotePrecision"`
	UnderlyingType        string        `json:"underlyingType"`
	UnderlyingSubType     []string      `json:"underlyingSubType"`
	TriggerProtect        ujson.Float64 `json:"triggerProtect"`
	LiquidationFee        ujson.Float64 `json:"liquidationFee"`
	MarketTakeBound       ujson.Float64 `json:"marketTakeBound"`
	MaxMoveOrderLimit     int           `json:"maxMoveOrderLimit"`

	Filters        []Filter      `json:"filters"`
	OrderTypes     []OrderType   `json:"orderTypes"`
	TimeInForce    []TimeInForce `json:"timeInForce"`
	PermissionSets []interface{} `json:"permissionSets"`

	// only USDT-Margin/Spot fields:
	Status string `json:"status"`

	// only COIN-Margin fields:
	ContractStatus    string `json:"contractStatus"`
	ContractSize      int    `json:"contractSize"`
	EqualQtyPrecision int    `json:"equalQtyPrecision"`

	// only Spot
	QuoteAssetPrecision             int      `json:"quoteAssetPrecision"`
	BaseCommissionPrecision         int      `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision        int      `json:"quoteCommissionPrecision"`
	IcebergAllowed                  bool     `json:"icebergAllowed"`
	OcoAllowed                      bool     `json:"ocoAllowed"`
	OtoAllowed                      bool     `json:"otoAllowed"`
	QuoteOrderQtyMarketAllowed      bool     `json:"quoteOrderQtyMarketAllowed"`
	AllowTrailingStop               bool     `json:"allowTrailingStop"`
	CancelReplaceAllowed            bool     `json:"cancelReplaceAllowed"`
	AmendAllowed                    bool     `json:"amendAllowed"`
	IsSpotTradingAllowed            bool     `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed          bool     `json:"isMarginTradingAllowed"`
	DefaultSelfTradePreventionMode  string   `json:"defaultSelfTradePreventionMode"`
	Permissions                     []string `json:"permissions"`
	AllowedSelfTradePreventionModes []string `json:"allowedSelfTradePreventionModes"`
}

type Filter struct {
	FilterType          string        `json:"filterType"`
	MaxPrice            ujson.Float64 `json:"maxPrice,omitempty"`
	MinPrice            ujson.Float64 `json:"minPrice,omitempty"`
	TickSize            ujson.Float64 `json:"tickSize,omitempty"`
	MaxQty              ujson.Float64 `json:"maxQty,omitempty"`
	MinQty              ujson.Float64 `json:"minQty,omitempty"`
	StepSize            ujson.Float64 `json:"stepSize,omitempty"`
	Limit               ujson.Float64 `json:"limit,omitempty"`
	Notional            ujson.Float64 `json:"notional,omitempty"`
	MultiplierUp        ujson.Float64 `json:"multiplierUp,omitempty"`
	MultiplierDown      ujson.Float64 `json:"multiplierDown,omitempty"`
	MultiplierDecimal   ujson.Float64 `json:"multiplierDecimal,omitempty"`
	PositionControlSide *string       `json:"positionControlSide,omitempty"`

	// spot only
	MinTrailingAboveDelta ujson.Float64 `json:"minTrailingAboveDelta,omitempty"`
	MaxTrailingAboveDelta ujson.Float64 `json:"maxTrailingAboveDelta,omitempty"`
	MinTrailingBelowDelta ujson.Float64 `json:"minTrailingBelowDelta,omitempty"`
	MaxTrailingBelowDelta ujson.Float64 `json:"maxTrailingBelowDelta,omitempty"`
	BidMultiplierUp       ujson.Float64 `json:"bidMultiplierUp,omitempty"`
	BidMultiplierDown     ujson.Float64 `json:"bidMultiplierDown,omitempty"`
	AskMultiplierUp       ujson.Float64 `json:"askMultiplierUp,omitempty"`
	AskMultiplierDown     ujson.Float64 `json:"askMultiplierDown,omitempty"`
	AvgPriceMins          ujson.Float64 `json:"avgPriceMins,omitempty"`
	MinNotional           ujson.Float64 `json:"minNotional,omitempty"`
	ApplyMinToMarket      bool          `json:"applyMinToMarket,omitempty"`
	MaxNotional           ujson.Float64 `json:"maxNotional,omitempty"`
	ApplyMaxToMarket      bool          `json:"applyMaxToMarket,omitempty"`
	MaxNumOrders          ujson.Float64 `json:"maxNumOrders,omitempty"`
	MaxNumAlgoOrders      ujson.Float64 `json:"maxNumAlgoOrders,omitempty"`
}

func (c *Client) GetExchangeInfo() Response[ExchangeInfo] {
	return ExchangeInfoRequest{}.Do(c)
}

type ExchangeInfoRequest struct{}

func (r ExchangeInfoRequest) Do(c *Client) Response[ExchangeInfo] {
	path := GetPath(c.category, "exchangeInfo")
	return GetPub(c, path, r, identity[ExchangeInfo])
}

type BookTicker struct {
	Symbol       string        `json:"symbol"`
	Pair         string        `json:"pair"`
	BidPrice     ujson.Float64 `json:"bidPrice"`
	BidQty       ujson.Float64 `json:"bidQty"`
	AskPrice     ujson.Float64 `json:"askPrice"`
	AskQty       ujson.Float64 `json:"askQty"`
	Time         int64         `json:"time"`
	LastUpdateId int64         `json:"lastUpdateId"`
}

type BookTickerRequest struct {
	Symbol string `url:",omitempty"`
}

func (c *Client) GetBookTickers(req BookTickerRequest) Response[[]BookTicker] {
	if req.Symbol == "" || c.category == CoinMargin {
		return req.Do(c)
	} else {
		respSingle := req.DoSingle(c)
		respArray := Response[[]BookTicker]{
			Data:       []BookTicker{respSingle.Data},
			Limit:      respSingle.Limit,
			Error:      respSingle.Error,
			StatusCode: respSingle.StatusCode,
			NetError:   respSingle.NetError,
		}
		return respArray
	}

}

func (r BookTickerRequest) Do(c *Client) Response[[]BookTicker] {
	path := GetPath(c.category, "ticker/bookTicker")
	return GetPub(c, path, r, identity[[]BookTicker])
}

func (r BookTickerRequest) DoSingle(c *Client) Response[BookTicker] {
	path := GetPath(c.category, "ticker/bookTicker")
	return GetPub(c, path, r, identity[BookTicker])
}

type Candle struct {
	OpenTime        int64   `json:"-"`
	Open            float64 `json:"-"`
	High            float64 `json:"-"`
	Low             float64 `json:"-"`
	Close           float64 `json:"-"`
	Volume          float64 `json:"-"`
	CloseTime       int64   `json:"-"`
	AssetVolume     float64 `json:"-"` // Quote for linear/spot, Base for inverse
	NumberOfTrades  int     `json:"-"`
	TakerBuyVolume1 float64 `json:"-"` // Taker buy base asset volume for linear/spot. // Taker buy volume - for inverse.
	TakerBuyVolume2 float64 `json:"-"` // Taker buy quote asset volume for linear/spot. // Taker buy base asset volume - for inverse
}

type RawCandle [12]interface{}

func ConvertRawToCandle(raw RawCandle) (Candle, error) {
	var c Candle
	var err error

	// Helpers
	toInt64 := func(v interface{}) (int64, error) {
		switch val := v.(type) {
		case float64:
			return int64(val), nil
		case string:
			return strconv.ParseInt(val, 10, 64)
		default:
			return 0, fmt.Errorf("cannot convert %v (%T) to int64", v, v)
		}
	}

	toFloat64 := func(v interface{}) (float64, error) {
		switch val := v.(type) {
		case float64:
			return val, nil
		case string:
			return strconv.ParseFloat(val, 64)
		default:
			return 0, fmt.Errorf("cannot convert %v (%T) to float64", v, v)
		}
	}

	toInt := func(v interface{}) (int, error) {
		switch val := v.(type) {
		case float64:
			return int(val), nil
		case string:
			n, err := strconv.Atoi(val)
			return n, err
		default:
			return 0, fmt.Errorf("cannot convert %v (%T) to int", v, v)
		}
	}

	// Conversion
	if c.OpenTime, err = toInt64(raw[0]); err != nil {
		return c, err
	}
	if c.Open, err = toFloat64(raw[1]); err != nil {
		return c, err
	}
	if c.High, err = toFloat64(raw[2]); err != nil {
		return c, err
	}
	if c.Low, err = toFloat64(raw[3]); err != nil {
		return c, err
	}
	if c.Close, err = toFloat64(raw[4]); err != nil {
		return c, err
	}
	if c.Volume, err = toFloat64(raw[5]); err != nil {
		return c, err
	}
	if c.CloseTime, err = toInt64(raw[6]); err != nil {
		return c, err
	}
	if c.AssetVolume, err = toFloat64(raw[7]); err != nil {
		return c, err
	}
	if c.NumberOfTrades, err = toInt(raw[8]); err != nil {
		return c, err
	}
	if c.TakerBuyVolume1, err = toFloat64(raw[9]); err != nil {
		return c, err
	}
	if c.TakerBuyVolume2, err = toFloat64(raw[10]); err != nil {
		return c, err
	}
	return c, nil
}

type CandleRequest struct {
	Symbol    string
	Interval  Bar
	StartTime int64 `url:",omitempty"`
	EndTime   int64 `url:",omitempty"`
	Limit     int   `url:",omitempty"`
}

func (o *Client) GetCandle(v CandleRequest) Response[[]Candle] {
	return v.Do(o)
}

func (o CandleRequest) Do(c *Client) Response[[]Candle] {
	path := GetPath(c.category, "klines")
	return GetPub(c, path, o, func(l []RawCandle) (r []Candle, err error) {
		for _, v := range l {
			var s Candle
			s, err = ConvertRawToCandle(v)
			if err != nil {
				break
			}
			r = append(r, s)
		}
		return
	})
}

func GetPath(c Category, path string) string {
	// for public methods
	switch c {
	case Spot:
		path = "/api/v3/" + path
	case CoinMargin:
		path = "/dapi/v1/" + path
	case UsdtMargin:
		path = "/fapi/v1/" + path
	}
	return path
}
