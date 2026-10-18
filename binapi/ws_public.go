package binapi

import (
	"strings"
	"time"

	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/uws"
)

// WsPublic is the facade of WsStreamClient for the Binance endpoint
// /public: high-frequency realtime streams (bookTicker).
type WsPublic struct {
	c *WsStreamClient
}

// NewWsPublic creates a client of the /public endpoint.
func NewWsPublic() *WsPublic {
	o := new(WsPublic)
	o.c = NewWsStreamClient(WsPublicBaseUrl)
	return o
}

// BookTicker subscribes to the "<symbol>@bookTicker" streams of the
// given symbols: full best bid/offer, realtime.
func (o *WsPublic) BookTicker(symbols ...string) *WsPublic {
	o.c.WithStreams(streamNames("bookTicker", symbols...)...)
	return o
}

// Close stops the client.
func (o *WsPublic) Close() {
	o.c.Close()
}

// Transport returns the uws options: timeouts, proxy and message
// logging can be tuned there before Run.
func (o *WsPublic) Transport() *uws.Options {
	return o.c.Transport()
}

// Connected reports whether the WebSocket connection is established.
func (o *WsPublic) Connected() bool {
	return o.c.Connected()
}

// WithLog sets the log of the underlying client.
func (o *WsPublic) WithLog(log *ulog.Log) *WsPublic {
	o.c.WithLog(log)
	return o
}

// WithBaseUrl overrides the base URL (for tests or alternative regions).
func (o *WsPublic) WithBaseUrl(url string) *WsPublic {
	o.c.WithBaseUrl(url)
	return o
}

// WithProxy sets the proxy URL.
func (o *WsPublic) WithProxy(proxy string) *WsPublic {
	o.c.WithProxy(proxy)
	return o
}

// WithTimeUnitMicroseconds adds &timeUnit=MICROSECOND to the URL.
func (o *WsPublic) WithTimeUnitMicroseconds(enable bool) *WsPublic {
	o.c.WithTimeUnitMicroseconds(enable)
	return o
}

// WithReconnectEvery sets the planned reconnect period; zero disables.
func (o *WsPublic) WithReconnectEvery(v time.Duration) *WsPublic {
	o.c.WithReconnectEvery(v)
	return o
}

// WithOnStreamEvent sets the callback receiving every stream message.
func (o *WsPublic) WithOnStreamEvent(f func(StreamEvent)) *WsPublic {
	o.c.WithOnStreamEvent(f)
	return o
}

// WithOnConnected sets the callback invoked after every (re)connection.
func (o *WsPublic) WithOnConnected(f func()) *WsPublic {
	o.c.WithOnConnected(f)
	return o
}

// WithOnDisconnected sets the callback invoked after every disconnect.
func (o *WsPublic) WithOnDisconnected(f func()) *WsPublic {
	o.c.WithOnDisconnected(f)
	return o
}

// WithOnDialError sets the callback invoked on connection failures.
func (o *WsPublic) WithOnDialError(f func(error) bool) *WsPublic {
	o.c.WithOnDialError(f)
	return o
}

// Run connects and keeps the connection alive until Close.
func (o *WsPublic) Run() {
	o.c.Run()
}

// streamNames builds full stream names from a stream type and symbols:
// "BTCUSDT" -> "btcusdt@bookTicker".
func streamNames(stream string, symbols ...string) []string {
	l := make([]string, 0, len(symbols))
	for _, s := range symbols {
		l = append(l, strings.ToLower(s)+"@"+stream)
	}
	return l
}
