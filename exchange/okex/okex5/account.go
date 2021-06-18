package okex5

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
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

func (rc *RestClient) Finance(ctx context.Context, req *exchange.FinanceReqParam) ([]exchange.Finance, error) {
	var it InstType
	switch req.Symbol.(type) {
	case exchange.MarginSymbol:
		it = InstTypeMargin

	case exchange.OptionSymbol:
		it = InstTypeOption

	case exchange.SwapSymbol:
		it = InstTypeSwap

	case exchange.SpotSymbol:
		it = InstTypeSpot

	}

	param := BillReq{
		InstType: it,
		Before:   req.StartID,
		After:    req.EndID,
	}
	if req.Limit != 0 {
		param.Limit = strconv.Itoa(req.Limit)
	}

	if req.Type == exchange.FinanceTypeFunding {
		param.Type = "8"
	} else if req.Type == exchange.FinanceTypeInterest {
		param.Type = "7"
	} else {
		return nil, errors.Errorf("other type is not support yet")
	}

	bills, err := rc.Bills(ctx, &param)
	if err != nil {
		return nil, err
	}

	s := req.Symbol.String()
	ret := []exchange.Finance{}
	for _, b := range bills {
		if s != b.InstID {
			continue
		}

		f, err := b.Parse()
		if err != nil {
			return nil, err
		}

		ret = append(ret, *f)
	}
	return ret, nil
}

func (b *Bill) Parse() (*exchange.Finance, error) {
	ts, err := ParseTimestamp(b.Ts)
	if err != nil {
		return nil, err
	}

	var symbol exchange.Symbol
	switch b.InstType {
	case InstTypeSpot:
		symbol, err = ParseSpotSymbol(b.InstID)
		if err != nil {
			return nil, err
		}

	case InstTypeMargin:
		symbol, err = ParseMarginSymbol(b.InstID)
		if err != nil {
			return nil, err
		}

	case InstTypeSwap:
		symbol, err = ParseSwapSymbol(b.InstID)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.Errorf("not support instType %s", b.InstType)
	}

	ret := &exchange.Finance{
		ID:       b.BillID,
		Symbol:   symbol,
		Time:     ts,
		Currency: b.Ccy,
		Raw:      *b,
	}

	if b.Type == "8" {
		ret.Type = exchange.FinanceTypeFunding
		amount, err := decimal.NewFromString(b.BalChg)
		if err != nil {
			return nil, errors.WithMessagef(err, "invalid balChg %s", b.BalChg)
		}
		ret.Amount = amount

	} else if b.Type == "7" {
		ret.Type = exchange.FinanceTypeInterest
		amount, err := decimal.NewFromString(b.Sz)
		if err != nil {
			return nil, errors.WithMessagef(err, "invalid sz %s", b.Sz)
		}
		ret.Amount = amount

	} else {
		return nil, errors.Errorf("unsupport type '%s'", b.Type)
	}

	return ret, nil
}
