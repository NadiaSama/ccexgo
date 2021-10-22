package spot

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Account struct {
		ID      int64  `json:"id"`
		Type    string `json:"type"`
		State   string `json:"state"`
		SubType string `json:"subtype"`
	}

	BalanceReq struct {
		AccountID int
	}

	Balance struct {
		Currency string `json:"currency"`
		Type     string `json:"type"`
		Balance  string `json:"balance"`
		SeqNum   string `json:"seq-num"`
	}
	BalanceResp struct {
		ID    int       `json:"id"`
		Type  string    `json:"type"`
		State string    `json:"state"`
		List  []Balance `json:"list"`
	}
)

const (
	AccountsEndPoint = "/v1/account/accounts"
	TypeFrozen       = "frozen"
	TypeTrade        = "trade"
)

//Init spot account id for Balance request
func (rc *RestClient) Init(ctx context.Context) error {
	accounts, err := rc.Accounts(ctx)
	if err != nil {
		return err
	}

	for _, ac := range accounts {
		if ac.Type == "spot" {
			rc.spotAccountID = int(ac.ID)
			return nil
		}
	}
	return errors.Errorf("no spot account")
}

func (rc *RestClient) Accounts(ctx context.Context) ([]Account, error) {
	var ret []Account

	if err := rc.Request(ctx, http.MethodGet, AccountsEndPoint, nil, nil, true, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (rc *RestClient) Balance(ctx context.Context, req *BalanceReq) (*BalanceResp, error) {
	endPoint := fmt.Sprintf("%s/%d/balance", AccountsEndPoint, req.AccountID)

	var ret BalanceResp
	if err := rc.RestClient.Request(ctx, http.MethodGet, endPoint, nil, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch balance fail")
	}

	return &ret, nil
}

func (rc *RestClient) FetchBalance(ctx context.Context, currencies ...string) (*exchange.Balances, error) {
	if rc.spotAccountID == 0 {
		return nil, errors.Errorf("client not init yet")
	}

	resp, err := rc.Balance(ctx, &BalanceReq{
		AccountID: rc.spotAccountID,
	})

	if err != nil {
		return nil, err
	}

	m := map[string]*exchange.Balance{}
	for _, b := range resp.List {
		currency := strings.ToUpper(b.Currency)
		amount, err := decimal.NewFromString(b.Balance)
		if err != nil {
			return nil, errors.WithMessagef(err, "invalid balance for currency '%s'", b.Currency)
		}

		bal, ok := m[currency]
		if !ok {
			bal = &exchange.Balance{
				Currency: currency,
			}
			m[currency] = bal
		}

		if b.Type == TypeFrozen {
			bal.Frozen = bal.Frozen.Add(amount)
			bal.Total = bal.Total.Add(amount)
		} else if b.Type == TypeTrade {
			bal.Free = bal.Free.Add(amount)
			bal.Total = bal.Total.Add(amount)
		} else {
			return nil, errors.Errorf("unsupport balance type currency '%s' type '%s'", b.Currency, b.Type)
		}
	}

	ret := exchange.NewBalances()
	ret.Raw = resp

	if len(currencies) != 0 {
		for _, c := range currencies {
			c = strings.ToUpper(c)
			bal, ok := m[c]
			if !ok {
				ret.Balances[c] = &exchange.Balance{
					Currency: c,
				}
			} else {
				ret.Balances[c] = bal
			}
		}
	} else {
		for k, v := range m {
			ret.Balances[k] = v
		}
	}

	return ret, nil
}
