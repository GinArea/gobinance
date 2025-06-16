package binapi

import (
	"net/http"
	"time"

	"github.com/msw-x/moon/uhttp"
)

type Client struct {
	c                *uhttp.Client
	s                *Sign
	category         Category // TODO may be remove?
	onTransportError OnTransportError
}

func NewClient(c Category) *Client {
	o := new(Client)
	o.c = uhttp.NewClient()
	o.WithCategory(c)
	return o
}

func (o *Client) WithBaseUrl(url string) *Client {
	o.c.WithBase(url)
	return o
}

func (o *Client) WithCategory(c Category) *Client {
	o.category = c
	switch c {
	case UsdtMargin:
		o.WithBaseUrl(UsdtMarginBaseUrl)
	case CoinMargin:
		o.WithBaseUrl(CoinMarginBaseUrl)
	case Spot:
		o.WithBaseUrl(SpotBaseUrl)
	}
	return o
}

func (o *Client) WithTimeout(timeout time.Duration) *Client {
	o.c.WithTimeout(timeout)
	return o
}

func (o *Client) WithTrace(trace func(uhttp.Response)) *Client {
	o.c.WithTrace(trace)
	return o
}

func (o *Client) WithProxy(proxy string) *Client {
	o.c.WithProxy(proxy)
	return o
}

func (o *Client) WithTransport(transport *http.Transport) *Client {
	o.c.WithTransport(transport)
	return o
}

func (o *Client) Copy() *Client {
	r := new(Client)
	r.c = o.c.Copy()
	r.s = o.s
	r.onTransportError = o.onTransportError
	return r
}

func (o *Client) WithPath(path string) *Client {
	o.c.WithPath(path)
	return o
}

func (o *Client) WithAppendPath(path string) *Client {
	o.c.WithAppendPath(path)
	return o
}

func (o *Client) WithAuth(key, secret string) *Client {
	o.s = NewSign(key, secret)
	return o
}

func (o *Client) WithOnReadBodyError(f uhttp.OnError) *Client {
	o.c.WithOnReadBodyError(f)
	return o
}

func (o *Client) WithOnTransportError(f OnTransportError) *Client {
	o.onTransportError = f
	return o
}

type OnTransportError func(err error, method string, statusCode int, attempt int) bool
