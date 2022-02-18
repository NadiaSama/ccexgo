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
	"strings"
	"time"

	"github.com/pkg/errors"
)

type (
	RestClient struct {
		key        string
		secret     string
		subAccount string
		prefix     string
	}

	Wrap struct {
		Success bool            `json:"success"`
		Result  json.RawMessage `json:"result"`
		Error   string          `json:"error"`
	}
)

const (
	ftxExchange = "ftx"
	ftxRSAddr   = "https://ftx.com/api"
)

func NewRestClient(key, secret string) *RestClient {
	return &RestClient{
		key:    key,
		secret: secret,
		prefix: ftxRSAddr,
	}
}

func NewClientWithSubAccount(key, secret, subAccount string) *RestClient {
	return &RestClient{
		key:        key,
		secret:     secret,
		subAccount: subAccount,
		prefix:     ftxRSAddr,
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

	var r Wrap
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	if !r.Success {
		return errors.Errorf("response error %v", r)
	}

	if err := json.Unmarshal(r.Result, &dst); err != nil {
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
		encStr := fmt.Sprintf("%d%s%s", ts, method, u.Path)
		if u.RawQuery != "" {
			encStr += "?" + u.RawQuery
		}

		if body != nil {
			data, err := ioutil.ReadAll(body)
			if err != nil {
				return nil, err
			}
			encStr += string(data)
			body = bytes.NewBuffer(data)
		}

		signature := signature(rc.secret, encStr)
		req, err = http.NewRequestWithContext(ctx, method, uStr, body)
		req.Header.Add("FTX-KEY", rc.key)
		req.Header.Add("FTX-SIGN", signature)
		req.Header.Add("FTX-TS", fmt.Sprintf("%d", ts))
		req.Header.Add("Content-Type", "application/json")

	} else {
		req, err = http.NewRequestWithContext(ctx, method, uStr, body)
	}

	// 子账户,子账户用户名跟在secret后方
	if rc.subAccount != "" {
		req.Header.Add("FTX-SUBACCOUNT", rc.subAccount)
	}
	return req, err
}

func signature(secret string, param string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(param))
	return fmt.Sprintf("%x", h.Sum(nil))
}
