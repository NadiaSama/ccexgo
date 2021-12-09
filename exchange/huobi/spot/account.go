package spot

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

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

	AccountHistoryReq struct {
		*exchange.RestReq
	}

	AccountHistoryData struct {
		AccountID    int    `json:"account-id"`
		Currency     string `json:"currency"`
		RecordID     int    `json:"record-id"`
		TransactAmt  string `json:"transact-amt"`
		TransactType string `json:"transact-type"`
		AvailBalance string `json:"avail-balance"`
		AcctBalance  string `json:"acct-balance"`
		TransactTime int    `json:"transact-time"`
	}
	AccountHistoryResp struct {
		Status  string               `json:"status"`
		NextID  int                  `json:"next-id"`
		ErrMsg  string               `json:"err-msg"`
		ErrCode string               `json:"err-code"`
		Data    []AccountHistoryData `json:"data"`
	}

	AccountLedgerReq struct {
		*exchange.RestReq
	}

	AccountLedgerData struct {
		AccountID    int     `json:"accountId"`
		Currency     string  `json:"currency"`
		TransactAmt  float64 `json:"transactAmt"`
		TransactType string  `json:"transactType"`
		TransferType string  `json:"transferType"`
		TransactID   int     `json:"transactId"`
		TransactTime int64   `json:"transactTime"`
		Transferer   int     `json:"transferer"`
		Transferee   int     `json:"transferee"`
	}

	AccountLedgerResp struct {
		Code    int                 `json:"code"`
		Message string              `json:"message"`
		NextID  int                 `json:"nextId"`
		OK      bool                `json:"ok"`
		Data    []AccountLedgerData `json:"data"`
	}
)

const (
	AccountsEndPoint       = "/v1/account/accounts"
	AccountHistoryEndPoint = "/v1/account/history"
	AccountLedgerEndPoint  = "/v2/account/ledger"
	TypeFrozen             = "frozen"
	TypeTrade              = "trade"
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

func NewAccountHistoryReq(uid int) *AccountHistoryReq {
	req := exchange.NewRestReq()
	req.AddFields("account-id", uid)
	return &AccountHistoryReq{
		RestReq: req,
	}
}

func (ar *AccountHistoryReq) Currency(cy string) *AccountHistoryReq {
	ar.AddFields("currency", cy)
	return ar
}

func (ar *AccountHistoryReq) TransactTypes(types ...string) *AccountHistoryReq {
	typ := strings.Join(types, ",")
	ar.AddFields("transact-types", typ)
	return ar
}

func (ar *AccountHistoryReq) AddTime(ts time.Time) *AccountHistoryReq {
	ar.AddFields("start-time", ts.Unix()*1000)
	return ar
}

func (ar *AccountHistoryReq) EndTime(ts time.Time) *AccountHistoryReq {
	ar.AddFields("end-time", ts.Unix()*1000)
	return ar
}

func (ar *AccountHistoryReq) Sort(sort string) *AccountHistoryReq {
	ar.AddFields("sort", sort)
	return ar
}

func (ar *AccountHistoryReq) Size(size int) *AccountHistoryReq {
	ar.AddFields("size", size)
	return ar
}

func (ar *AccountHistoryReq) FromID(id int) *AccountHistoryReq {
	ar.AddFields("from-id", id)
	return ar
}

func (rc *RestClient) AccountHistory(ctx context.Context, req *AccountHistoryReq) (*AccountHistoryResp, error) {
	values, err := req.Values()
	if err != nil {
		return nil, errors.WithMessage(err, "build request param fail")
	}

	var ret AccountHistoryResp
	if err := rc.RequestWithRawResp(ctx, http.MethodGet, AccountHistoryEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "request failed")
	}

	if ret.Status != "ok" {
		return nil, errors.Errorf("reqeust response failed %+v", ret)
	}
	return &ret, nil
}

func NewAccountLedgerReq(uid int) *AccountLedgerReq {
	req := exchange.NewRestReq()
	req.AddFields("accountId", uid)
	return &AccountLedgerReq{
		RestReq: req,
	}
}

func (ar *AccountLedgerReq) Currency(cy string) *AccountLedgerReq {
	ar.AddFields("currency", cy)
	return ar
}

func (ar *AccountLedgerReq) TransactTypes(typ string) *AccountLedgerReq {
	ar.AddFields("transactTypes", typ)
	return ar
}

func (ar *AccountLedgerReq) StartTime(st time.Time) *AccountLedgerReq {
	ar.AddFields("startTime", st.Unix()*1000)
	return ar
}

func (ar *AccountLedgerReq) EndTime(et time.Time) *AccountLedgerReq {
	ar.AddFields("endTime", et.Unix()*1000)
	return ar
}

func (ar *AccountLedgerReq) Sort(s string) *AccountLedgerReq {
	ar.AddFields("sort", s)
	return ar
}

func (ar *AccountLedgerReq) Limit(limit int) *AccountLedgerReq {
	ar.AddFields("limit", limit)
	return ar
}

func (ar *AccountLedgerReq) FromID(id int) *AccountLedgerReq {
	ar.AddFields("fromId", id)
	return ar
}

func (rc *RestClient) AccountLedger(ctx context.Context, req *AccountLedgerReq) (*AccountLedgerResp, error) {
	values, err := req.Values()
	if err != nil {
		return nil, errors.WithMessage(err, "build request param fail")
	}

	var ret AccountLedgerResp
	if err := rc.RequestWithRawResp(ctx, http.MethodGet, AccountLedgerEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "build request fail")
	}

	if !ret.OK {
		return nil, errors.Errorf("response error %+v", ret)
	}

	return &ret, nil
}
