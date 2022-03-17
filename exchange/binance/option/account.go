package option

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Account struct {
		Currency       string          `json:"currency"`
		Equity         decimal.Decimal `json:"equity"`
		Available      decimal.Decimal `json:"available"`
		OverMargin     decimal.Decimal `json:"overMargin"`
		PositionMargin decimal.Decimal `json:"positionMargin"`
		UnrealizedPNL  decimal.Decimal `json:"unrealizedPNL"`
		MaintMargin    decimal.Decimal `json:"maintMargin"`
		Balance        decimal.Decimal `json:"balance"`
	}
)

const (
	AccountEndPoint = "/vapi/v1/account"
)

func (rc *RestClient) Account(ctx context.Context) ([]Account, error) {
	var ac []Account

	if err := rc.GetRequest(ctx, AccountEndPoint, binance.NewRestReq(), true, &ac); err != nil {
		return nil, err
	}
	return ac, nil
}

func (rc *RestClient) FetchBalance(ctx context.Context, currencies ...string) ([]exchange.Balance, error) {
	account, err := rc.Account(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "fetch account fail")
	}

	ret := make([]exchange.Balance, len(account))
	for i, ac := range account {
		ret[i] = exchange.Balance{
			Currency: ac.Currency,
			Equitity: ac.Equity,
			Total:    ac.Balance,
			Free:     ac.Available,
		}
	}
	return ret, nil
}
