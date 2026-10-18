package binapi

import (
	"encoding/json"
)

// StreamEvent is a single message received from a combined stream
// connection. The combined stream wrapper is already removed:
// Data contains the raw payload of the stream named by Stream.
// TsRecv is the reception timestamp taken in nanoseconds
// (time.Now().UnixNano()) as the first action of the message handler,
// before any unmarshalling.
type StreamEvent struct {
	Stream string          // "btcusdt@bookTicker"
	Data   json.RawMessage // payload without the combined stream wrapper
	TsRecv int64           // reception timestamp, nanoseconds
}

// BookTickerEvent is the payload of the "<symbol>@bookTicker" stream
// (endpoint /public, realtime push of the full best bid/offer).
// Prices and quantities are strings: Binance sends them as strings,
// the precision is preserved; conversion is up to the consumer.
// Timestamps are int64; with timeUnit=MICROSECOND only the order of
// magnitude changes, parsing stays the same.
type BookTickerEvent struct {
	Event        string `json:"e"` // "bookTicker"
	UpdateId     int64  `json:"u"`
	EventTime    int64  `json:"E"`
	TransactTime int64  `json:"T"`
	Symbol       string `json:"s"`
	BidPrice     string `json:"b"`
	BidQty       string `json:"B"`
	AskPrice     string `json:"a"`
	AskQty       string `json:"A"`
}

// AggTradeEvent is the payload of the "<symbol>@aggTrade" stream
// (endpoint /market, aggregated trades pushed every 100 ms).
// Prices and quantities are strings: Binance sends them as strings,
// the precision is preserved; conversion is up to the consumer.
type AggTradeEvent struct {
	Event        string `json:"e"` // "aggTrade"
	EventTime    int64  `json:"E"`
	AggTradeId   int64  `json:"a"`
	Symbol       string `json:"s"`
	Price        string `json:"p"`
	Qty          string `json:"q"`
	FirstTradeId int64  `json:"f"`
	LastTradeId  int64  `json:"l"`
	TradeTime    int64  `json:"T"`
	BuyerIsMaker bool   `json:"m"` // true: the buyer is the maker, the aggressor is the seller
}

// ParseBookTickerEvent decodes a raw bookTicker payload.
func ParseBookTickerEvent(data []byte) (r BookTickerEvent, err error) {
	err = json.Unmarshal(data, &r)
	return
}

// ParseAggTradeEvent decodes a raw aggTrade payload.
func ParseAggTradeEvent(data []byte) (r AggTradeEvent, err error) {
	err = json.Unmarshal(data, &r)
	return
}
