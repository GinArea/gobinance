package binapi

import (
	"time"

	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/uws"
)

// WsMarket is the facade of WsStreamClient for the Binance endpoint
// /market: regular market streams pushed every 100 ms (aggTrade).
type WsMarket struct {
	c *WsStreamClient
}

// NewWsMarket creates a client of the /market endpoint.
func NewWsMarket() *WsMarket {
	o := new(WsMarket)
	o.c = NewWsStreamClient(WsMarketBaseUrl)
	return o
}

// AggTrade subscribes to the "<symbol>@aggTrade" streams of the given
// symbols: aggregated trades, pushed every 100 ms.
func (o *WsMarket) AggTrade(symbols ...string) *WsMarket {
	o.c.WithStreams(streamNames("aggTrade", symbols...)...)
	return o
}

// Close stops the client.
func (o *WsMarket) Close() {
	o.c.Close()
}

// Transport returns the uws options: timeouts, proxy and message
// logging can be tuned there before Run.
func (o *WsMarket) Transport() *uws.Options {
	return o.c.Transport()
}

// Connected reports whether the WebSocket connection is established.
func (o *WsMarket) Connected() bool {
	return o.c.Connected()
}

// WithLog sets the log of the underlying client.
func (o *WsMarket) WithLog(log *ulog.Log) *WsMarket {
	o.c.WithLog(log)
	return o
}

// WithBaseUrl overrides the base URL (for tests or alternative regions).
func (o *WsMarket) WithBaseUrl(url string) *WsMarket {
	o.c.WithBaseUrl(url)
	return o
}

// WithProxy sets the proxy URL.
func (o *WsMarket) WithProxy(proxy string) *WsMarket {
	o.c.WithProxy(proxy)
	return o
}

// WithTimeUnitMicroseconds adds &timeUnit=MICROSECOND to the URL.
func (o *WsMarket) WithTimeUnitMicroseconds(enable bool) *WsMarket {
	o.c.WithTimeUnitMicroseconds(enable)
	return o
}

// WithReconnectEvery sets the planned reconnect period; zero disables.
func (o *WsMarket) WithReconnectEvery(v time.Duration) *WsMarket {
	o.c.WithReconnectEvery(v)
	return o
}

// WithOnStreamEvent sets the callback receiving every stream message.
func (o *WsMarket) WithOnStreamEvent(f func(StreamEvent)) *WsMarket {
	o.c.WithOnStreamEvent(f)
	return o
}

// WithOnConnected sets the callback invoked after every (re)connection.
func (o *WsMarket) WithOnConnected(f func()) *WsMarket {
	o.c.WithOnConnected(f)
	return o
}

// WithOnDisconnected sets the callback invoked after every disconnect.
func (o *WsMarket) WithOnDisconnected(f func()) *WsMarket {
	o.c.WithOnDisconnected(f)
	return o
}

// WithOnDialError sets the callback invoked on connection failures.
func (o *WsMarket) WithOnDialError(f func(error) bool) *WsMarket {
	o.c.WithOnDialError(f)
	return o
}

// Run connects and keeps the connection alive until Close.
func (o *WsMarket) Run() {
	o.c.Run()
}
