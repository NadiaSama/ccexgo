package ftx

import (
	"context"
	"fmt"
	"net/http"
)

type (
	Balance struct {
		Coin  string  `json:"coin"`
		Free  float64 `json:"free"`
		Total float64 `json:"total"`
	}
)

const (
	walletEndPoint = "/wallet"
)

func (rc *RestClient) Balances(ctx context.Context) ([]Balance, error) {
	var ret []Balance
	endPoint := fmt.Sprintf("%s/balances", walletEndPoint)

	if err := rc.request(ctx, http.MethodGet, endPoint, nil, nil, true, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}
