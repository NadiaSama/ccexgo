package deribit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type (
	RestClient struct {
		key    string
		secret string
		prefix string
	}
)

func NewRestClient(key, secret string) *RestClient {
	return newRestClientWithPrefix(key, secret, "https://www.deribit.com/api/v2")
}

func NewTestRestClient(key, secret string) *RestClient {
	return newRestClientWithPrefix(key, secret, "https://test.deribit.com/api/v2")
}

func newRestClientWithPrefix(key, secret, prefix string) *RestClient {
	return &RestClient{
		key:    key,
		secret: secret,
		prefix: prefix,
	}
}

func (rc *RestClient) Request(ctx context.Context, method string, endPoint string, params url.Values, body io.Reader, signed bool, dst interface{}) error {
	if signed {
		return errors.Errorf("signed rest request is not support yet")
	}

	url := fmt.Sprintf("%s%s", rc.prefix, endPoint)
	if len(params) != 0 {
		url = fmt.Sprintf("%s?%s", url, params.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return errors.WithMessage(err, "build request fail")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.WithMessage(err, "do http request fail")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.WithMessage(err, "read response fail")
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.Errorf("invalid statusCode %d status %s", resp.StatusCode, resp.Status)
	}

	var r Response
	if err := json.Unmarshal(data, &r); err != nil {
		return errors.WithMessage(err, "unmarshal json error")
	}
	if r.Error.Code != 0 {
		return errors.Errorf("response error code: %d message: %s", r.Error.Code, r.Error.Message)
	}

	if err := json.Unmarshal(r.Result, &dst); err != nil {
		return errors.WithMessage(err, "unmarshal Result error")
	}
	return nil
}
