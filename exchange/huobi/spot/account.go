package spot

import (
	"context"
	"net/http"
)

type (
	Account struct {
		ID      int64  `json:"id"`
		Type    string `json:"type"`
		State   string `json:"state"`
		SubType string `json:"subtype"`
	}
)

const (
	AccountsEndPoint = "/v1/account/accounts"
)

func (rc *RestClient) Accounts(ctx context.Context) ([]Account, error) {
	var ret []Account

	if err := rc.Request(ctx, http.MethodGet, AccountsEndPoint, nil, nil, true, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}
