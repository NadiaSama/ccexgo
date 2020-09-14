package okex

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/NadiaSama/ccexgo/misc/request"
	"github.com/pkg/errors"
)

type (
	RestClient struct {
		key        string
		secret     string
		passPhrase string
		apiHost    string
	}
)

func NewRestClient(key, secret, passPhrase, apiHost string) *RestClient {
	return &RestClient{
		key:        key,
		secret:     secret,
		passPhrase: passPhrase,
		apiHost:    apiHost,
	}
}

func (rc *RestClient) request(ctx context.Context, method, endPoint string, param map[string]string, data interface{}, sign bool, dst interface{}) error {
	req, err := rc.buildRequest(ctx, method, endPoint, param, data, sign)
	if err != nil {
		return err
	}

	return request.DoReqWithCtx(req, func(resp *http.Response, ierr error) error {
		if ierr != nil {
			return ierr
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return errors.Errorf("bad response %s", string(body))
		}
		if err := json.Unmarshal(body, dst); err != nil {
			return err
		}
		return nil
	})
}

func (rc *RestClient) buildRequest(ctx context.Context, method, endPoint string, param map[string]string, data interface{}, sign bool) (*http.Request, error) {
	var (
		body string
		req  *http.Request
		err  error
	)

	values := url.Values{}
	for k, v := range param {
		values.Add(k, v)
	}
	u := url.URL{Scheme: "https", Host: rc.apiHost, Path: endPoint, RawQuery: values.Encode()}
	if method == http.MethodPost {
		b, e := json.Marshal(data)
		if e != nil {
			return nil, errors.WithMessagef(err, "build request fail")
		}
		body = string(b)
		req, err = http.NewRequestWithContext(ctx, method, u.String(), bytes.NewBuffer(b))
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", "application/json")
	} else if method == http.MethodGet {
		req, err = http.NewRequestWithContext(ctx, method, u.String(), nil)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.Errorf("unknown method %s", method)
	}

	if sign {
		p := u.Path
		if u.RawQuery != "" {
			p = fmt.Sprintf("%s?%s", u.Path, u.RawQuery)
		}
		ts := time.Now().UTC().Format(time.RFC3339)
		raw := fmt.Sprintf("%s%s%s%s", ts, method, p, body)
		h := hmac.New(sha256.New, []byte(rc.secret))
		h.Write([]byte(raw))
		signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
		req.Header.Add("OK-ACCESS-KEY", rc.key)
		req.Header.Add("OK-ACCESS-SIGN", signature)
		req.Header.Add("OK-ACCESS-TIMESTAMP", ts)
		req.Header.Add("OK-ACCESS-PASSPHRASE", rc.passPhrase)
	}
	return req, nil
}
