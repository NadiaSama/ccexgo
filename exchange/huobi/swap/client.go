package swap

import (
	"bytes"
	"context"
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

func (rc *RestClient) PrivatePostReq(ctx context.Context, endPoint string, sr Serializer, dst interface{}) error {
	raw, err := sr.Serialize()
	if err != nil {
		return errors.WithMessage(err, "serialize fail")
	}
	buf := bytes.NewBuffer(raw)
	return rc.Request(ctx, http.MethodPost, endPoint, nil, buf, true, dst)
}
