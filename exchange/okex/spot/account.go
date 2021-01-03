package spot

import (
	"context"
	"net/http"

	"github.com/shopspring/decimal"
)

type (
	AccountResp struct {
		ID        string          `json:"id"`
		Available decimal.Decimal `json:"available"`
		Balance   decimal.Decimal `json:"balance"`
		Currency  string          `json:"currency"`
		Hold      decimal.Decimal `json:"hold"`
		Frozen    decimal.Decimal `json:"frozen"`
		Holds     decimal.Decimal `json:"holds"`
	}

	Accounts []AccountResp
)

const (
	AccountEndPoint = "/api/spot/v3/accounts"
)

func (rc *RestClient) FetchAccounts(ctx context.Context) (Accounts, error) {
	var ret Accounts
	if err := rc.Request(ctx, http.MethodGet, AccountEndPoint, nil, nil, true, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}
