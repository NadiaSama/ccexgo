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
	signatureMethod  = "HmacSHA256"
	signatureVersion = "2"
	scheme           = "https"
	statusOK         = "ok"
	CodeOK           = 200
)

type (
	RestClient struct {
		key         string
		secret      string
		apiHost     string
		pair2Symbol map[string]exchange.SpotSymbol
	}
)

func NewRestClient(key, secret, host string) *RestClient {
	return &RestClient{
		key:         key,
		secret:      secret,
		apiHost:     host,
		pair2Symbol: make(map[string]exchange.SpotSymbol),
	}
}

func (rc *RestClient) Init(ctx context.Context) error {
	return rc.initSymbol(ctx)
}

//Request build and send huobi raw request
func (rc *RestClient) Request(ctx context.Context, method string, endPoint string, param url.Values, body io.Reader, sign bool, dst interface{}) error {
	req, err := rc.buildRequest(ctx, method, endPoint, param, body, sign)
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

func (rc *RestClient) buildRequest(ctx context.Context, method, endPoint string, values url.Values, body io.Reader, sign bool) (*http.Request, error) {
	var query string
	if sign {
		ts := time.Now().UTC()
		values.Add("AccessKeyId", rc.key)
		values.Add("SignatureMethod", signatureMethod)
		values.Add("SignatureVersion", signatureVersion)
		values.Add("Timestamp", ts.Format("2006-01-02T15:04:05"))
		query = values.Encode()
		sig := rc.signature(method, rc.apiHost, endPoint, query)
		query = fmt.Sprintf("%s&Signature=%s", query, url.QueryEscape(sig))
	} else {
		query = values.Encode()
	}
	u := url.URL{Scheme: scheme, Path: endPoint, RawQuery: query, Host: rc.apiHost}

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
	fields := []string{method, host, path, query}
	raw := strings.Join(fields, "\n")

	hash := hmac.New(sha256.New, []byte(rc.secret))
	hash.Write([]byte(raw))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
