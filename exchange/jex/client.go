package jex

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

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	Client struct {
		*exchange.Client
		apiHost string
	}
)

func NewClient(key, secret string) *Client {
	return &Client{
		Client: &exchange.Client{
			Key:    key,
			Secret: secret,
		},
		apiHost: "testnet.jex.com",
	}
}

func (c *Client) Exchange() string {
	return "jex"
}

//request send a jex rest request
func (c *Client) request(ctx context.Context, method string, uri string, param map[string]string,
	dest interface{}, sign bool) (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "request %s %s", method, uri)
		}
	}()

	var req *http.Request
	req, err = c.buildRequest(ctx, method, uri, param, sign)
	if err != nil {
		return err
	}
	echan := make(chan error)
	httpc := http.DefaultClient
	go func() {
		defer close(echan)
		resp, err := httpc.Do(req)
		if err != nil {
			echan <- errors.WithMessage(err, "Do request error")
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			echan <- errors.WithMessage(err, "read body error")
			return
		}
		if resp.StatusCode != 200 {
			echan <- errors.Errorf("bad response %s", string(body))
			return
		}
		if err := json.Unmarshal(body, dest); err != nil {
			echan <- errors.WithMessagef(err, "unmarshal response error %s", string(body))
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err, ok := <-echan:
		if ok {
			return nil
		}
		return err
	}
}

func (c *Client) buildURL(uri string) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   c.apiHost,
		Path:   uri,
	}
}

func (c *Client) buildRequest(ctx context.Context, method string, uri string, param map[string]string,
	sign bool) (*http.Request, error) {

	url := c.buildURL(uri)
	fields := []string{}
	for key, val := range param {
		fields = append(fields, fmt.Sprintf("%s=%s", key, val))
	}
	if sign {
		fields = append(fields, fmt.Sprintf("timestamp=%d", (time.Now().UnixNano()/1e6)))
	}
	query := strings.Join(fields, "&")
	if sign {
		h := hmac.New(sha256.New, []byte(c.Secret))
		x := h.Sum([]byte(query))
		query = fmt.Sprintf("%s&signature=%x", query, x)
	}
	var body io.Reader
	if method == "GET" {
		url.RawQuery = query
		body = nil
	} else if method == "POST" || method == "DELETE" {
		body = bytes.NewReader([]byte(query))
	} else {
		return nil, errors.Errorf("unsupport http method %s", method)
	}

	req, err := http.NewRequestWithContext(ctx, method, url.String(), body)
	if err != nil {
		return nil, errors.WithMessagef(err, "build request fail")
	}
	req.Header.Add("X-JEX-APIKEY", c.Key)
	return req, nil
}
