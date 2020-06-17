package jex

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
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

func NewClient(ctx context.Context, key, secret string) *Client {
	return &Client{
		Client: &exchange.Client{
			Ctx:     ctx,
			Timeout: time.Second * 2,
			Key:     key,
			Secret:  secret,
		},
		apiHost: "testnet.jex.com",
	}
}

//Request send a jex rest request
func (c *Client) Request(method string, uri string, param map[string]string, sign bool) ([]byte, error) {
	ctx, cancel := context.WithTimeout(c.Ctx, c.Timeout)
	defer cancel()
	req, err := c.buildRequest(ctx, method, uri, param, sign)
	if err != nil {
		return nil, err
	}
	echan := make(chan error)
	result := make(chan []byte)
	httpc := http.DefaultClient
	go func() {
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
		result <- body
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-echan:
		return nil, err
	case ret := <-result:
		return ret, nil
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
	} else if method == "POST" {
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
