package ftx

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type (
	RestClient struct {
		key    string
		secret string
		prefix string
	}
)

func NewRestClient(key, secret string) *RestClient {
	return &RestClient{
		key:    key,
		secret: secret,
		prefix: "https://ftx.com/api",
	}
}

func (rc *RestClient) request(ctx context.Context, method string, endPoint string, params url.Values, body io.Reader, sign bool, dst interface{}) error {
	req, err := rc.buildRequest(ctx, method, endPoint, params, body, sign)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(data, &dst); err != nil {
		return err
	}
	return nil
}

func (rc *RestClient) buildRequest(ctx context.Context, method string, endPoint string, params url.Values, body io.Reader, sign bool) (*http.Request, error) {
	var (
		req *http.Request
		err error
	)

	uStr := fmt.Sprintf("%s%s", rc.prefix, endPoint)
	if params != nil {
		uStr = fmt.Sprintf("%s?%s", uStr, params.Encode())
	}

	u, err := url.Parse(uStr)
	if err != nil {
		return nil, err
	}

	if sign {
		ts := time.Now().UnixNano() / 1e6
		encStr := fmt.Sprintf("%d%s%s%s", ts, method, u.Path, u.RawQuery)
		if body != nil {
			data, err := ioutil.ReadAll(body)
			if err != nil {
				return nil, err
			}
			encStr += string(data)
			body = bytes.NewBuffer(data)
		}

		signature := rc.signature(encStr)
		req, err = http.NewRequestWithContext(ctx, method, uStr, body)
		req.Header.Add("FTX-KEY", rc.key)
		req.Header.Add("FTX-SIGN", signature)
		req.Header.Add("FTX-TS", fmt.Sprintf("%d", ts))
		req.Header.Add("Content-Type", "application/json")

	} else {
		req, err = http.NewRequestWithContext(ctx, method, uStr, body)
	}
	return req, err
}

func (rc *RestClient) signature(param string) string {
	h := hmac.New(sha256.New, []byte(rc.secret))
	h.Write([]byte(param))
	return fmt.Sprintf("%x", h.Sum(nil))
}
