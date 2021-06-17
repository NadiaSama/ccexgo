package okex5

import (
	"context"
	"net/http"
	"net/url"
)

type (
	Bill struct {
		InstType  InstType `json:"instType"`
		BillID    string   `json:"billId"`
		Type      string   `json:"type"`
		SubType   string   `json:"subType"`
		Ts        string   `json:"ts"`
		BalChg    string   `json:"balChg"`
		PosBalChg string   `json:"posBalChg"`
		Bal       string   `json:"bal"`
		PosBal    string   `json:"posBal"`
		Sz        string   `json:"sz"`
		Ccy       string   `json:"ccy"`
		Pnl       string   `json:"pnl"`
		Fee       string   `json:"fee"`
		MgnMode   MgnMode  `json:"mgnMode"`
		InstID    string   `json:"instId"`
		OrdID     string   `json:"ordId"`
		From      string   `json:"from"`
		To        string   `json:"to"`
		Notes     string   `json:"notes"`
	}

	BillReq struct {
		InstType InstType
		Ccy      string
		MgnMode  MgnMode
		CtType   CtType
		Type     string
		SubType  string
		After    string
		Before   string
		Limit    string
	}
)

const (
	BillsEndPoint = "/api/v5/account/bills"
)

func (rc *RestClient) Bills(ctx context.Context, req *BillReq) ([]Bill, error) {
	values := url.Values{}
	if req.InstType != "" && req.InstType != InstTypeAny {
		values.Add("instType", string(req.InstType))
	}
	if req.Ccy != "" {
		values.Add("ccy", req.Ccy)
	}

	if req.MgnMode != "" {
		values.Add("mgnMode", string(req.MgnMode))
	}

	if req.CtType != CtTypeNone {
		values.Add("ctType", string(req.CtType))
	}

	if req.Type != "" {
		values.Add("type", req.Type)
	}

	if req.SubType != "" {
		values.Add("subType", req.SubType)
	}

	if req.Before != "" {
		values.Add("before", req.Before)
	}

	if req.After != "" {
		values.Add("after", req.After)
	}

	if req.Limit != "" {
		values.Add("limit", req.Limit)
	}

	var ret []Bill
	if err := rc.Request(ctx, http.MethodGet, BillsEndPoint, values, nil, true, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}
