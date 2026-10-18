package binapi

import (
	"encoding/json"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/uws"
)

// Base URLs of the Binance USDⓈ-M Futures WebSocket endpoints.
// The legacy /ws and /stream endpoints were disabled by Binance
// on 2026-04-23; only /public and /market are used.
const (
	WsPublicBaseUrl = "wss://fstream.binance.com/public" // bookTicker (realtime)
	WsMarketBaseUrl = "wss://fstream.binance.com/market" // aggTrade (100 ms)

	// WsDefaultReconnectEvery is the default planned reconnect period:
	// a Binance connection lives at most 24 hours, the client
	// reconnects after 23 hours to stay ahead of the server.
	WsDefaultReconnectEvery = 23 * time.Hour

	// wsMaxReconnectJitter is the cap of the planned reconnect jitter:
	// the effective delay is WsDefaultReconnectEvery ± jitter, where
	// jitter is min(wsMaxReconnectJitter, reconnectEvery/20).
	wsMaxReconnectJitter = 60 * time.Second
)

// wsStreamWrapper is the combined stream message envelope:
// {"stream":"<name>","data":{...}}.
type wsStreamWrapper struct {
	Stream string          `json:"stream"`
	Data   json.RawMessage `json:"data"`
}

// WsStreamClient is a reconnecting client of Binance public WebSocket
// streams in the combined stream mode: streams are fixed in the URL as
// <base>/stream?streams=<s1>/<s2>/..., every message arrives wrapped as
// {"stream":...,"data":...}. The wrapper is removed by the client, the
// consumer receives StreamEvent with a reception timestamp TsRecv.
//
// Streams are fixed before Run; dynamic SUBSCRIBE/UNSUBSCRIBE is not
// supported. The event callback is executed synchronously in the read
// flow: a slow consumer blocks message reception, buffering into a
// channel is the responsibility of the consumer.
//
// Liveness: WebSocket ping frames keep the connection alive - by
// default the uws core sends one every 24 s (Binance answers with
// pong); server ping frames are answered with pong by gorilla/websocket
// automatically. After an abnormal disconnect the core reconnects with
// a fixed re-dial interval. The connection also has a limited lifetime
// on the Binance side (24 h), so the client performs a planned
// reconnect after 23 h by default (jitter ±60 s): a timer sends a
// close frame, the uws loop then re-establishes the same URL, restoring
// all stream subscriptions.
//
// Implementation note: while the planned reconnect is enabled the
// client disables the built-in uws ping job and sends both the keepalive
// pings and the planned close frame from its own timer goroutine - the
// uws ping job of moon v0.7.47 is not race-clean (its send book-keeping
// and the stop of the ping job are unsynchronized). All client sends
// stay confined to a single goroutine, the ping cadence is captured
// from Transport().PingInterval at Run (24 s by default).
type WsStreamClient struct {
	c *uws.Client

	streams              []string
	onStreamEvent        func(StreamEvent)
	reconnectEvery       time.Duration
	timeUnitMicroseconds bool

	onConnectedCallback    func()
	onDisconnectedCallback func()
	onDialErrorCallback    func(error) bool

	mutex       sync.Mutex
	closed      bool
	timer       *time.Timer
	reconnectAt time.Time
	pingEvery   time.Duration
	// jitterFn returns the planned reconnect jitter offset;
	// overridable in tests for determinism.
	jitterFn func() time.Duration
}

// NewWsStreamClient creates a client for the given base URL
// (WsPublicBaseUrl or WsMarketBaseUrl). The planned reconnect is
// enabled with the default period of 23 h.
func NewWsStreamClient(base string) *WsStreamClient {
	o := new(WsStreamClient)
	o.c = uws.NewClient(base)
	o.reconnectEvery = WsDefaultReconnectEvery
	return o
}

// Close stops the client and the planned reconnect timer.
func (o *WsStreamClient) Close() {
	o.mutex.Lock()
	o.closed = true
	o.stopTimer()
	o.mutex.Unlock()
	o.c.Close()
}

// Transport returns the uws options: timeouts, proxy and message
// logging can be tuned there before Run.
func (o *WsStreamClient) Transport() *uws.Options {
	return &o.c.Options
}

// Connected reports whether the WebSocket connection is established.
func (o *WsStreamClient) Connected() bool {
	return o.c.Connected()
}

// WithLog sets the log of the underlying client.
func (o *WsStreamClient) WithLog(log *ulog.Log) *WsStreamClient {
	o.c.WithLog(log)
	return o
}

// WithBaseUrl overrides the base URL (for tests or alternative regions).
func (o *WsStreamClient) WithBaseUrl(url string) *WsStreamClient {
	o.c.WithBase(url)
	return o
}

// WithProxy sets the proxy URL.
func (o *WsStreamClient) WithProxy(proxy string) *WsStreamClient {
	o.c.WithProxy(proxy)
	return o
}

// WithStreams sets the stream subscriptions: full stream names, e.g.
// "btcusdt@bookTicker". Symbol parts are converted to lower case.
// Must be called before Run; the set is fixed afterwards.
func (o *WsStreamClient) WithStreams(streams ...string) *WsStreamClient {
	o.streams = make([]string, 0, len(streams))
	for _, s := range streams {
		o.streams = append(o.streams, normalizeStreamName(s))
	}
	return o
}

// WithOnStreamEvent sets the callback receiving every stream message.
func (o *WsStreamClient) WithOnStreamEvent(f func(StreamEvent)) *WsStreamClient {
	o.onStreamEvent = f
	return o
}

// WithTimeUnitMicroseconds adds &timeUnit=MICROSECOND to the URL:
// Binance then sends event timestamps (E, T) in microseconds instead
// of milliseconds.
func (o *WsStreamClient) WithTimeUnitMicroseconds(enable bool) *WsStreamClient {
	o.timeUnitMicroseconds = enable
	return o
}

// WithReconnectEvery sets the planned reconnect period. Zero disables
// the planned reconnect. Default: WsDefaultReconnectEvery (23 h).
// While enabled, the keepalive ping cadence is captured from
// Transport().PingInterval at Run and the pings are sent by the client
// itself.
func (o *WsStreamClient) WithReconnectEvery(v time.Duration) *WsStreamClient {
	o.reconnectEvery = v
	return o
}

// WithOnConnected sets the callback invoked after every successful
// (re)connection.
func (o *WsStreamClient) WithOnConnected(f func()) *WsStreamClient {
	o.onConnectedCallback = f
	return o
}

// WithOnDisconnected sets the callback invoked after every disconnect.
func (o *WsStreamClient) WithOnDisconnected(f func()) *WsStreamClient {
	o.onDisconnectedCallback = f
	return o
}

// WithOnDialError sets the callback invoked on connection failures;
// returning true skips the re-dial delay.
func (o *WsStreamClient) WithOnDialError(f func(error) bool) *WsStreamClient {
	o.onDialErrorCallback = f
	return o
}

// Run connects and keeps the connection alive until Close: abnormal
// disconnects are handled by the uws re-dial loop, planned reconnects
// by the internal timer. Stream subscriptions are part of the URL and
// are restored automatically on every reconnect.
func (o *WsStreamClient) Run() {
	o.c.WithPath(o.streamPath())
	o.c.WithOnMessage(o.onMessage)
	o.c.WithOnConnected(o.onConnected)
	o.c.WithOnDisconnected(o.onDisconnected)
	if o.onDialErrorCallback != nil {
		o.c.WithOnDialError(o.onDialErrorCallback)
	}
	if o.reconnectEvery > 0 {
		// The planned reconnect owns the keepalive: capture the ping
		// cadence and disable the built-in uws ping job - the timer of
		// this client then sends pings and the planned close frame from
		// a single goroutine (see the implementation note above).
		o.pingEvery = o.c.Options.PingInterval
		o.c.Options.PingInterval = 0
	}
	o.c.Run()
}

// streamPath builds the combined stream path:
// stream?streams=<s1>/<s2>/...[&timeUnit=MICROSECOND].
func (o *WsStreamClient) streamPath() string {
	path := "stream?streams=" + strings.Join(o.streams, "/")
	if o.timeUnitMicroseconds {
		path += "&timeUnit=MICROSECOND"
	}
	return path
}

func (o *WsStreamClient) onMessage(messageType int, data []byte) {
	// The reception timestamp is the first action, before any parsing.
	tsRecv := time.Now().UnixNano()
	log := o.c.Log()
	if messageType != websocket.TextMessage {
		log.Warning("invalid message type:", uws.MessageTypeString(messageType))
		return
	}
	var w wsStreamWrapper
	if err := json.Unmarshal(data, &w); err != nil {
		log.Error("invalid stream message:", err)
		return
	}
	if w.Stream == "" || len(w.Data) == 0 {
		log.Error("invalid stream message: empty stream or data")
		return
	}
	if o.onStreamEvent != nil {
		o.onStreamEvent(StreamEvent{
			Stream: w.Stream,
			Data:   w.Data,
			TsRecv: tsRecv,
		})
	}
}

func (o *WsStreamClient) onConnected() {
	o.mutex.Lock()
	o.reconnectAt = time.Now().Add(o.reconnectEvery + o.nextJitter())
	o.mutex.Unlock()
	o.armTimer()
	if o.onConnectedCallback != nil {
		o.onConnectedCallback()
	}
}

func (o *WsStreamClient) onDisconnected() {
	if o.onDisconnectedCallback != nil {
		o.onDisconnectedCallback()
	}
}

// armTimer (re)arms the keepalive/planned-reconnect timer for the
// current connection: it fires at the earlier of the next keepalive
// ping and the planned reconnect deadline.
func (o *WsStreamClient) armTimer() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.closed || o.reconnectEvery <= 0 || o.reconnectAt.IsZero() {
		return
	}
	o.stopTimer()
	delay := time.Until(o.reconnectAt)
	if o.pingEvery > 0 && o.pingEvery < delay {
		delay = o.pingEvery
	}
	if delay <= 0 {
		delay = time.Millisecond
	}
	o.timer = time.AfterFunc(delay, o.onTimer)
}

func (o *WsStreamClient) stopTimer() {
	if o.timer != nil {
		o.timer.Stop()
		o.timer = nil
	}
}

// onTimer runs in the timer goroutine: until the planned reconnect is
// due it sends the keepalive ping, at the deadline a close frame - the
// uws read flow then receives the close response and the built-in loop
// re-dials the same URL.
func (o *WsStreamClient) onTimer() {
	o.mutex.Lock()
	closed := o.closed
	due := !o.reconnectAt.IsZero() && !time.Now().Before(o.reconnectAt)
	o.mutex.Unlock()
	if closed {
		return
	}
	if due {
		o.c.Log().Info("planned reconnect")
		err := o.c.Send(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "planned reconnect"))
		if err != nil {
			o.c.Log().Warning("planned reconnect: send close frame:", err)
		}
		// no re-arm here: the next connection re-arms in onConnected
		return
	}
	if o.pingEvery > 0 {
		if err := o.c.Send(websocket.PingMessage, nil); err != nil {
			o.c.Log().Warning("ping:", err)
		}
	}
	o.armTimer()
}

// nextJitter returns the planned reconnect jitter offset within
// ±min(wsMaxReconnectJitter, reconnectEvery/20).
func (o *WsStreamClient) nextJitter() time.Duration {
	if o.jitterFn != nil {
		return o.jitterFn()
	}
	j := wsMaxReconnectJitter
	if s := o.reconnectEvery / 20; s < j {
		j = s
	}
	if j <= 0 {
		return 0
	}
	return time.Duration(rand.Int63n(int64(2*j))) - j
}

// normalizeStreamName converts the symbol part of a stream name to
// lower case ("BTCUSDT@bookTicker" -> "btcusdt@bookTicker"); a name
// without "@" is lower case entirely.
func normalizeStreamName(s string) string {
	if i := strings.Index(s, "@"); i >= 0 {
		return strings.ToLower(s[:i]) + s[i:]
	}
	return strings.ToLower(s)
}
