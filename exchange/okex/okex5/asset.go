package okex5

import (
	"context"
	"net/http"
)

type (
	TransferParam struct {
		Ccy      string `json:"ccy"`
		Amt      string `json:"amt"`
		Type     string `json:"type"`
		From     string `json:"from"`
		To       string `json:"to"`
		SubAcct  string `json:"subAcct"`
		InstID   string `json:"instId"`
		ToInstID string `json:"toInstId"`
	}

	TransferResp struct {
		TransID string `json:"transId"`
		Ccy     string `json:"ccy"`
		From    string `json:"from"`
		Amt     string `json:"amt"`
		To      string `json:"to"`
	}

	CurrencyResp struct {
		Ccy         string `json:"ccy"`
		Name        string `json:"name"`
		Chain       string `json:"chain"`
		CanDep      bool   `json:"canDep"`
		CanWd       bool   `json:"canWd"`
		CanInternal bool   `json:"canInternal"`
		MinWd       string `json:"minWd"`
		MinFee      string `json:"minFee"`
		MaxFee      string `json:"maxFee"`
	}
)

const (
	TransferEndPoint   = "/api/v5/asset/transfer"
	CurrenciesEndPoint = "/api/v5/asset/currencies"
)

func (c *RestClient) Transfer(ctx context.Context, param *TransferParam) (*TransferResp, error) {
	var resp TransferResp
	if err := c.doPostJSON(ctx, TransferEndPoint, param, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *RestClient) Currencies(ctx context.Context) ([]CurrencyResp, error) {
	var resp []CurrencyResp
	if err := c.Request(ctx, http.MethodGet, CurrenciesEndPoint, nil, nil, true, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
