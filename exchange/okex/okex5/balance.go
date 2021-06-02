package okex5

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

type (
	AccountDetial struct {
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
		Upl           string `json:"upl"`
		UplLiab       string `json:"uplLiab"`
		CrossLiab     string `json:"crossLiab"`
		IsoLiab       string `json:"isoLiab"`
		MgnRatio      string `json:"mgnRatio"`
		Interest      string `json:"interest"`
		Twap          string `json:"twap"`
		MaxLoan       string `json:"maxLoan"`
		EqUsd         string `json:"eqUsd"`
		NotionalLever string `json:"notionalLever"`
	}

	AccountBalance struct {
		UTime       string          `json:"uTime"`
		TotalEq     string          `json:"totalEq"`
		IsoEq       string          `json:"isoEq"`
		AdjEq       string          `json:"adjEq"`
		OrdFroz     string          `json:"ordFroz"`
		Imr         string          `json:"imr"`
		Nmr         string          `json:"nmr"`
		MgnRatio    string          `json:"mgnRatio"`
		NotionalUSD string          `json:"notionalUsd"`
		Details     []AccountDetial `json:"details"`
	}
)

const (
	AccountBalanceEndPoint = "/api/v5/account/balance"
)

func (r *RestClient) AccountBalance(ctx context.Context, currency ...string) (*AccountBalance, error) {
	ret := []AccountBalance{}
	values := url.Values{}
	if len(currency) != 0 {
		ccy := strings.Join(currency, ",")
		values.Add("ccy", ccy)
	}

	if err := r.Request(ctx, http.MethodGet, AccountBalanceEndPoint, values, nil, true, &ret); err != nil {
		return nil, err
	}

	return &ret[0], nil
}
