package binapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/uhttp"
)

func GetPub[R, T any](c *Client, path string, req any, transform func(R) (T, error)) Response[T] {
	return request(c, http.MethodGet, path, req, transform, false)
}

func Get[R, T any](c *Client, path string, req any, transform func(R) (T, error)) Response[T] {
	return request(c, http.MethodGet, path, req, transform, true)
}

func request[R, T any](c *Client, method string, path string, request any, transform func(R) (T, error), sign bool) (r Response[T]) {
	var attempt int
	for {
		r = req(c, method, path, request, transform, sign)
		if r.StatusCode != http.StatusOK && c.onTransportError != nil {
			if c.onTransportError(r.Error, method, r.StatusCode, attempt) {
				attempt++
				continue
			}
		}
		break
	}
	return
}

func req[R, T any](c *Client, method string, path string, request any, transform func(R) (T, error), sign bool) (r Response[T]) {
	var perf *uhttp.Performer
	switch method {
	case http.MethodGet:
		perf = c.c.Get(path).Params(request)
		if sign && c.s != nil {
			perf.Params(c.s.getSignatureFields(perf.Request.Params))
		}
	case http.MethodPost:
		perf = c.c.Post(path).Json(request)
	default:
		r.Error = fmt.Errorf("forbidden method: %s", method)
		return
	}
	if sign && c.s != nil {
		if perf.Request.Header == nil {
			perf.Request.Header = make(http.Header)
		}
		switch method {
		case http.MethodGet:
			c.s.HeaderGet(perf.Request.Header)
		case http.MethodPost:
			//TODO
			//c.s.HeaderPost(perf.Request.Header, perf.Request.Body, path)
		}
	}
	httpResponse := perf.Do()
	if httpResponse.Error == nil {
		r.StatusCode = httpResponse.StatusCode
		if httpResponse.StatusCode == http.StatusOK {
			if httpResponse.BodyExists() {
				raw := new(response[R])
				r.Error = raw.parseJson(httpResponse)
				if r.Ok() {
					r.Error = raw.Error()
					if r.Ok() {
						r.Data, r.Error = transform(raw.Data)
					}
				}
			}
		} else {
			r.Error = errors.New(ufmt.Join(httpResponse.Status, httpResponse.Text()))
		}
		if sign {
			r.SetErrorIfNil(httpResponse.HeaderTo(&r.Limit))
		}
	} else {
		r.Error = httpResponse.Error
		r.NetError = true
	}
	return
}

func (r *response[T]) parseJson(data uhttp.Response) error {
	result := new(T)
	err := data.Json(result)
	if err == nil {
		r.Data = *result
		return nil
	}
	return err
}
