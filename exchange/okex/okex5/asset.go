package okex5

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
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

	AssetBillReq struct {
		*GetRequest
	}

	AssetBillRecord struct {
		BillID string `json:"billId"`
		Ccy    string `json:"ccy"`
		BalChg string `json:"balChg"`
		Bal    string `json:"bal"`
		Type   string `json:"type"`
		Ts     string `json:"ts"`
	}

	WithdrawlHistoryReq struct {
		*GetRequest
	}

	WithdrawlHistory struct {
		Chain string `json:"chain"`
		Fee   string `json:"fee"`
		Ccy   string `json:"ccy"`
		Amt   string `json:"amt"`
		TxID  string `json:"txId"`
		From  string `json:"from"`
		To    string `json:"to"`
		State string `json:"state"`
		TS    string `json:"ts"`
		WdID  string `json:"wdId"`
	}
)

const (
	TransferEndPoint         = "/api/v5/asset/transfer"
	CurrenciesEndPoint       = "/api/v5/asset/currencies"
	AssetBillsEndPoint       = "/api/v5/asset/bills"
	WithdrawlHistoryEndPoint = "/api/v5/asset/withdrawal-history"
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

func NewAssetBillReq() *AssetBillReq {
	return &AssetBillReq{
		GetRequest: NewGetRequest(),
	}
}
func (ar *AssetBillReq) Ccy(ccy string) *AssetBillReq {
	ar.Add("ccy", ccy)
	return ar
}

func (ar *AssetBillReq) Type(typ string) *AssetBillReq {
	ar.Add("type", typ)
	return ar
}

func (ar *AssetBillReq) BeforeTime(ts time.Time) *AssetBillReq {
	ar.Add("before", fmt.Sprintf("%d", ts.Unix()*1e3))
	return ar
}

func (ar *AssetBillReq) AfterTime(ts time.Time) *AssetBillReq {
	ar.Add("after", fmt.Sprintf("%d", ts.Unix()*1e3))
	return ar
}

func (ar *AssetBillReq) Limit(l string) *AssetBillReq {
	ar.Add("limit", l)
	return ar
}

func (c *RestClient) AssetBills(ctx context.Context, req *AssetBillReq) ([]AssetBillRecord, error) {
	var ret []AssetBillRecord
	if err := c.Request(ctx, http.MethodGet, BillsEndPoint, req.Values(), nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "request bills end point fail")
	}

	return ret, nil
}

func NewWithdrawlHistoryReq() *WithdrawlHistoryReq {
	return &WithdrawlHistoryReq{
		NewGetRequest(),
	}
}

func (wr *WithdrawlHistoryReq) Ccy(ccy string) *WithdrawlHistoryReq {
	wr.Add("ccy", ccy)
	return wr
}

func (wr *WithdrawlHistoryReq) AfterTime(ts time.Time) *WithdrawlHistoryReq {
	wr.Add("after", fmt.Sprintf("%d", ts.Unix()/1e6))
	return wr
}

func (c *RestClient) WithdrawlHistory(ctx context.Context, req *WithdrawlHistoryReq) ([]WithdrawlHistory, error) {
	var ret []WithdrawlHistory
	if err := c.Request(ctx, http.MethodGet, WithdrawlHistoryEndPoint, req.Values(), nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "request withdrawl history fail")
	}

	return ret, nil
}
