package spot

import (
	"context"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	DepositWithdrawlReq struct {
		*exchange.RestReq
	}

	DepositWithdrawlRecord struct {
		ID         int     `json:"id"`
		Type       string  `json:"type"`
		Currency   string  `json:"currency"`
		TxHash     string  `json:"tx-hash"`
		Chain      string  `json:"chain"`
		Amount     float64 `json:"amount"`
		Address    string  `json:"address"`
		AddressTag string  `json:"address-tag"`
		Fee        float64 `json:"fee"`
		State      string  `json:"state"`
		CreatedAt  int64   `json:"created-at"`
		UpdatedAt  int64   `json:"updated-at"`
		ErrorCode  string  `json:"error-code"`
		ErrorMsg   string  `json:"error-msg"`
	}
)

const (
	DepositWitdrawlEndPoint = "/v1/query/deposit-withdraw"
)

func NewDepositWithdrawlReq(typ string) *DepositWithdrawlReq {
	return &DepositWithdrawlReq{
		RestReq: exchange.NewRestReq().AddFields("type", typ),
	}
}

func (dr *DepositWithdrawlReq) Type(typ string) *DepositWithdrawlReq {
	dr.AddFields("type", typ)
	return dr
}

func (dr *DepositWithdrawlReq) Direct(d string) *DepositWithdrawlReq {
	dr.AddFields("direct", d)
	return dr
}

func (cl *RestClient) DepositWitdrawl(ctx context.Context, req *DepositWithdrawlReq) ([]DepositWithdrawlRecord, error) {
	var ret []DepositWithdrawlRecord
	values, err := req.Values()
	if err != nil {
		return nil, errors.WithMessage(err, "build request fail")
	}
	if err := cl.Request(ctx, http.MethodGet, DepositWitdrawlEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "request fail")
	}

	return ret, nil
}
