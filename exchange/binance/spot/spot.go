package spot

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	AccountBalance struct {
		Asset  string          `json:"asset"`
		Free   decimal.Decimal `json:"free"`
		Locked decimal.Decimal `json:"locked"`
	}

	AccountResp struct {
		MakerCommision int              `json:"makerCommision"`
		Balances       []AccountBalance `json:"balances"`
	}

	AccountReq struct {
		*binance.RestReq
	}
)

const (
	AccountEndPoint = "/api/v3/account"
)

func NewAccountReq() *AccountReq {
	return &AccountReq{
		RestReq: binance.NewRestReq(),
	}
}

func (rc *RestClient) Account(ctx context.Context, req *AccountReq) (*AccountResp, error) {
	var ret AccountResp
	if err := rc.GetRequest(ctx, AccountEndPoint, req, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "execute request fail")
	}

	return &ret, nil
}
