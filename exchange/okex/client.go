package okex

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/misc/request"
	"github.com/pkg/errors"
)

type (
	RestClient struct {
		key        string
		secret     string
		passPhrase string
		apiHost    string
		test       bool
	}
)

const (
	okexRestHost = "www.okx.com"
)

func NewRestClient(key, secret, passPhrase string) *RestClient {
	return &RestClient{
		key:        key,
		secret:     secret,
		passPhrase: passPhrase,
		apiHost:    okexRestHost,
		test:       false,
	}
}

func NewTESTRestClient(key, secret, passPhrase string) *RestClient {
	return &RestClient{
		key:        key,
		secret:     secret,
		passPhrase: passPhrase,
		apiHost:    okexRestHost,
		test:       true,
	}
}

func (rc *RestClient) Request(ctx context.Context, method, endPoint string, params url.Values, body io.Reader, sign bool, dst interface{}) error {
	p := map[string]string{}
	for k, v := range params {
		if len(v) != 0 {
			p[k] = v[0]
		}
	}

	return rc.request(ctx, method, endPoint, p, body, sign, dst)
}

func (rc *RestClient) Property() exchange.Property {
	return exchange.Property{
		Trades: &exchange.TradesProp{
			SuportID:    true,
			SupportTime: false,
		},
		Finance: &exchange.FinanceProp{
			SuportID:    true,
			SupportTime: false,
		},
	}
}
func (rc *RestClient) request(ctx context.Context, method, endPoint string, param map[string]string, body io.Reader, sign bool, dst interface{}) error {
	req, err := rc.buildRequest(ctx, method, endPoint, param, body, sign)
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
			return errors.WithMessage(err, "unmarshal response fail")
		}
		return nil
	})
}

func (rc *RestClient) buildRequest(ctx context.Context, method, endPoint string, param map[string]string, data io.Reader, sign bool) (*http.Request, error) {
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
		b, e := ioutil.ReadAll(data)
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

	if rc.test {
		req.Header.Add("x-simulated-trading", "1")
	}
	return req, nil
}
