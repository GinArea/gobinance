package binapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"
)

type Sign struct {
	Key    string
	Secret string
}

func NewSign(key, secret string) *Sign {
	o := new(Sign)
	o.Key = key
	o.Secret = secret
	return o
}

func (o *Sign) HeaderGet(h http.Header) {
	h.Set("X-MBX-APIKEY", o.Key)
}

func encodeSortParams(src url.Values) (s string) {
	if len(src) == 0 {
		return
	}
	keys := make([]string, len(src))
	i := 0
	for k := range src {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		s += encodeParam(k, src.Get(k)) + "&"
	}
	s = s[0 : len(s)-1]
	return
}

func encodeParam(name, value string) string {
	params := url.Values{}
	params.Add(name, value)
	return params.Encode()
}

func (o *Sign) computeHMACSHA256(message string) string {
	mac := hmac.New(sha256.New, []byte(o.Secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func (o *Sign) getSignatureFields(params url.Values) SignatureFields {
	encodedParams := encodeSortParams(params)
	timestamp := time.Now().UnixMilli() // milliseconds
	query := fmt.Sprintf("timestamp=%d", timestamp)
	if encodedParams != "" {
		query = fmt.Sprintf("%s&%s", encodedParams, query)
	}
	signature := o.computeHMACSHA256(query)
	return SignatureFields{
		Timestamp: timestamp,
		Signature: signature,
	}
}

type SignatureFields struct {
	Timestamp int64
	Signature string
}
