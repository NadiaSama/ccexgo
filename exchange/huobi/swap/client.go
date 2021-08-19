package swap

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange/huobi"
	"github.com/pkg/errors"
)

type (
	RestClient struct {
		*huobi.RestClient
	}

	Serializer interface {
		Serialize() ([]byte, error)
	}
)

const (
	SwapHost    = "api.hbdm.com"
	SwapProHost = "api.huobi.pro"
)

func NewRestClient(key string, secret string) *RestClient {
	return NewRestClientWithHost(key, secret, SwapHost)
}

func NewRestClientWithHost(key, secret, host string) *RestClient {
	return &RestClient{
		RestClient: huobi.NewRestClient(key, secret, host),
	}
}

//PrivatePostReq send post request to huobi swap api. the request body is generate from req param
//vai json.Marshal() or Serialize()
func (rc *RestClient) PrivatePostReq(ctx context.Context, endPoint string, req interface{}, dst interface{}) error {
	var (
		raw []byte
		err error
	)

	if sr, ok := req.(Serializer); ok {
		raw, err = sr.Serialize()
	} else {
		raw, err = json.Marshal(req)
	}
	if err != nil {
		return errors.WithMessage(err, "serialize fail")
	}
	buf := bytes.NewBuffer(raw)
	return rc.Request(ctx, http.MethodPost, endPoint, nil, buf, true, dst)
}
