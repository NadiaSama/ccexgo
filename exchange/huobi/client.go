package huobi

import (
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
	"strings"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/misc/request"
	"github.com/pkg/errors"
)

const (
	Huobi            = "huobi"
	signatureMethod  = "HmacSHA256"
	signatureVersion = "2"
	scheme           = "https"
	StatusOK         = "ok"
	CodeOK           = 200
)

type (
	RestClient struct {
		key     string
		secret  string
		apiHost string
	}

	RestResponse struct {
		Status string      `json:"status"`
		Code   int         `json:"code"`
		Data   interface{} `json:"data"`
	}
)

func NewRestClient(key, secret, host string) *RestClient {
	return &RestClient{
		key:     key,
		secret:  secret,
		apiHost: host,
	}
}

func (rc *RestClient) RequestWithRawResp(ctx context.Context, method string, endPoint string, param url.Values, body io.Reader, sign bool, dst interface{}) error {
	req, err := rc.buildRequest(ctx, method, rc.apiHost, endPoint, param, body, sign)
	if err != nil {
		return err
	}
	return request.DoReqWithCtx(req, func(resp *http.Response, ierr error) error {
		if ierr != nil {
			return ierr
		}
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if err := json.Unmarshal(content, dst); err != nil {
			return errors.WithMessagef(err, "unmarshal %s fail", string(content))
		}

		return nil
	})
}

func (rc *RestClient) Request(ctx context.Context, method string, endPoint string, param url.Values, body io.Reader, sign bool, dst interface{}) error {
	req, err := rc.buildRequest(ctx, method, rc.apiHost, endPoint, param, body, sign)
	if err != nil {
		return err
	}
	return request.DoReqWithCtx(req, func(resp *http.Response, ierr error) error {
		if ierr != nil {
			return ierr
		}
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		rr := RestResponse{
			Data: dst,
		}
		if err := json.Unmarshal(content, &rr); err != nil {
			return errors.WithMessagef(err, "unmarshal %s fail", string(content))
		}

		if (rr.Status != "" && rr.Status != "ok") || (rr.Code != 0 && rr.Code != 200) {
			return errors.Errorf("rest return error %s", string(content))
		}
		return nil
	})
}

func (rc *RestClient) Property() exchange.Property {
	return exchange.Property{
		Trades: &exchange.TradesProp{
			MaxDuration: time.Hour * 48,
			SuportID:    false,
			SupportTime: true,
		},
	}
}

func (rc *RestClient) buildRequest(ctx context.Context, method, host string, endPoint string, values url.Values, body io.Reader, sign bool) (*http.Request, error) {
	var query string
	if sign {
		if values == nil {
			values = url.Values{}
		}
		ts := time.Now().UTC()
		values.Add("AccessKeyId", rc.key)
		values.Add("SignatureMethod", signatureMethod)
		values.Add("SignatureVersion", signatureVersion)
		values.Add("Timestamp", ts.Format("2006-01-02T15:04:05"))
		query = values.Encode()
		sig := rc.signature(method, host, endPoint, query)
		query = fmt.Sprintf("%s&Signature=%s", query, url.QueryEscape(sig))
	} else {
		query = values.Encode()
	}
	u := url.URL{Scheme: scheme, Path: endPoint, RawQuery: query, Host: host}

	if method == http.MethodPost || method == http.MethodPut {
		req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", "application/json")
		return req, nil

	} else if method == http.MethodGet {
		return http.NewRequestWithContext(ctx, method, u.String(), nil)
	} else {
		return nil, errors.Errorf("unsuport method %s", method)
	}
}

func (rc *RestClient) signature(method, host, path, query string) string {
	return Signature(rc.secret, method, host, path, query)
}

func Signature(secret, method, host, path, query string) string {
	fields := []string{method, host, path, query}
	raw := strings.Join(fields, "\n")

	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(raw))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
