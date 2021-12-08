package okex5

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type (
	SubAccount struct {
		Enable  bool   `json:"enable"`
		SubAcct string `json:"subAcct"`
		Label   string `json:"label"`
		Mobile  string `json:"mobile"`
		GAuth   bool   `json:"gAuth"`
		TS      string `json:"ts"`
	}

	SubAccountBalancesResp struct {
		AdjEq       string                    `json:"adjEq"`
		TotalEq     string                    `json:"totalEq"`
		IsoEq       string                    `json:"isoEq"`
		OrdFroz     string                    `json:"ordFroz"`
		Imr         string                    `json:"imr"`
		Mmr         string                    `json:"mmr"`
		MgnRatio    string                    `json:"mgnRatio"`
		NotionalUsd string                    `json:"notionalUsd"`
		UTime       string                    `json:"uTime"`
		Details     []SubAccountBalanceDetail `json:"details"`
	}

	SubAccountBalanceDetail struct {
		Ccy           string `json:"ccy"`
		Eq            string `json:"eq"`
		CashBal       string `json:"cashBal"`
		UTime         string `json:"uTime"`
		IsoEq         string `json:"isoEq"`
		AvailEq       string `json:"availEq"`
		DisEq         string `json:"disEq"`
		AvailBal      string `json:"availBal"`
		FrozenBal     string `json:"frozenBal"`
		OrdFrozen     string `json:"ordFrozen"`
		Liab          string `json:"liab"`
		Upl           string `json:"uplLiab"`
		CrossLiab     string `json:"crossLiab"`
		IsoLiab       string `json:"isoLiab"`
		MgnRatio      string `json:"mgnRatio"`
		Interest      string `json:"interest"`
		Twap          string `json:"twap"`
		MaxLoan       string `json:"maxLoan"`
		EqUsd         string `json:"eqUsd"`
		NotionalLever string `json:"notionalLever"`
	}

	SubAccountBalancesReq struct {
		*GetRequest
	}

	SubAccountBillsReq struct {
		*GetRequest
	}

	SubAccountBill struct {
		BillID  string `json:"billId"`
		Type    string `json:"type"`
		Ccy     string `json:"ccy"`
		Amt     string `json:"amt"`
		SubAcct string `json:"subAcct"`
		TS      string `json:"ts"`
	}
)

const (
	SubAccountListEndpoint     = "/api/v5/users/subaccount/list"
	SubAccountBillsEndPoint    = "/api/v5/asset/subaccount/bills"
	SubAccountBalancesEndPoint = "/api/v5/account/subaccount/balances"
)

func (rc *RestClient) SubAccounts(ctx context.Context) ([]SubAccount, error) {
	var ret []SubAccount
	if err := rc.Request(ctx, http.MethodGet, SubAccountListEndpoint, nil, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch subaccount fail")
	}

	return ret, nil
}

func (rc *RestClient) SubAccountBalances(ctx context.Context, req *SubAccountBalancesReq) (*SubAccountBalancesResp, error) {
	var ret []SubAccountBalancesResp
	if err := rc.Request(ctx, http.MethodGet, SubAccountBalancesEndPoint, req.Values(), nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch balances fail")
	}

	return &ret[0], nil
}

func (rc *RestClient) SubAccountBills(ctx context.Context, req *SubAccountBillsReq) ([]SubAccountBill, error) {
	var ret []SubAccountBill
	if err := rc.Request(ctx, http.MethodGet, SubAccountBillsEndPoint, req.Values(), nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch subAccount bills fail")
	}

	return ret, nil
}

func NewBalancesReq(subAcct string) *SubAccountBalancesReq {
	gt := NewGetRequest()
	gt.Add("subAcct", subAcct)
	return &SubAccountBalancesReq{
		GetRequest: gt,
	}
}

func NewBillsReq() *SubAccountBillsReq {
	return &SubAccountBillsReq{
		GetRequest: NewGetRequest(),
	}
}

func (br *SubAccountBillsReq) Ccy(ccy string) *SubAccountBillsReq {
	br.Add("ccy", ccy)
	return br
}

func (br *SubAccountBillsReq) Type(typ string) *SubAccountBillsReq {
	br.Add("type", typ)
	return br
}

func (br *SubAccountBillsReq) SubAcct(at string) *SubAccountBillsReq {
	br.Add("subAcct", at)
	return br
}

func (br *SubAccountBillsReq) AfterTime(ts time.Time) *SubAccountBillsReq {
	br.Add("after", strconv.FormatInt(ts.Unix()*1000, 10))
	return br
}

func (br *SubAccountBillsReq) BeforeTime(ts time.Time) *SubAccountBillsReq {
	br.Add("before", strconv.FormatInt(ts.Unix()*1000, 10))
	return br
}

func (br *SubAccountBillsReq) Limit(lmt int) *SubAccountBillsReq {
	br.Add("limit", strconv.Itoa(lmt))
	return br
}
