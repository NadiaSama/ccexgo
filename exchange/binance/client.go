package binance

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/misc/request"
)

type (
	//Binance Rest client instance
	RestClient struct {
		key         string
		secret      string
		apiHost     string
		pair2Symbol map[string]exchange.SpotSymbol
	}
)

func NewRestClient(key, secret, host string) *RestClient {
	ret := &RestClient{
		key:         key,
		secret:      secret,
		apiHost:     host,
		pair2Symbol: make(map[string]exchange.SpotSymbol),
	}
	return ret
}

func (rc *RestClient) Init(ctx context.Context) error {
	return rc.initPair(ctx)
}

func (rc *RestClient) signature(param string) string {
	h := hmac.New(sha256.New, []byte(rc.secret))
	h.Write([]byte(param))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (rc *RestClient) request(ctx context.Context, endPoint string, param map[string]string, signed bool, dst interface{}) error {
	if signed {
		ts := timeStamp()
		if param == nil {
			param = map[string]string{}
		}
		param["timestamp"] = fmt.Sprintf("%d", ts)
	}
	query := rc.buildQuery(param, signed)

	u := url.URL{Scheme: "https", Host: rc.apiHost, Path: endPoint, RawQuery: query}
	req, _ := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	req.Header.Add("X-MBX-APIKEY", rc.key)

	rerr := request.DoReqWithCtx(req, func(resp *http.Response, ierr error) error {
		if ierr != nil {
			return ierr
		}
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			var ret error
			if strings.HasPrefix(endPoint, "/wapi") {
				var we WAPIError
				if ret = json.Unmarshal(content, &we); ret != nil {
					return ret
				}
				ret = &we
				return ret

			} else {
				var ae APIError
				if ret := json.Unmarshal(content, &ae); ret != nil {
					return ret
				}
				ret = &ae
				return ret
			}
		}

		if err := json.Unmarshal(content, dst); err != nil {
			return err
		}
		return nil
	})
	return rerr
}

func timeStamp() int64 {
	now := time.Now()
	return now.UnixNano() / 1e6
}

func (rc *RestClient) buildQuery(param map[string]string, signed bool) string {
	fields := make([]string, 0)
	for k, v := range param {
		fields = append(fields, fmt.Sprintf("%s=%s", k, v))
	}
	query := strings.Join(fields, "&")

	if signed {
		sig := rc.signature(query)
		query = fmt.Sprintf("%s&signature=%s", query, sig)
	}
	return query
}
