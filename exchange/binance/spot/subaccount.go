package spot

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange/binance"
)

type (
	SubAccount struct {
		Email      string `json:"email"`
		IsFreeze   bool   `json:"isFreeze"`
		CreateTime int64  `json:"createTime"`
	}

	SubAccountResp struct {
		SubAccounts []SubAccount `json:"subAccounts"`
	}

	SubAccountReq struct {
		*binance.RestReq
	}

	SubAccountTransferHistoryReq struct {
		*binance.RestReq
	}

	SubAccountTransferHistory struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Asset  string `json:"asset"`
		Qty    string `json:"qty"`
		Status string `json:"status"`
		TranID int    `json:"tranId"`
		Time   int64  `json:"time"`
	}

	SubAccountFuturesInternalTransferReq struct {
		*binance.RestReq
	}

	SubAccountFuturesInternalTransferResp struct {
		Success     bool                                `json:"success"`
		FuturesType int                                 `json:"futuresType"`
		Transfers   []SubAccountFuturesInternalTransfer `json:"transfers"`
	}

	SubAccountFuturesInternalTransfer struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Asset  string `json:"asset"`
		Qty    string `json:"qty"`
		TranID int    `json:"tranId"`
		Time   int64  `json:"time"`
	}

	SubAccountAssetReq struct {
		*binance.RestReq
	}

	SubAccountBalance struct {
		Asset  string  `json:"asset"`
		Free   float64 `json:"free"`
		Locked float64 `json:"locked"`
	}

	SubAccountAssetResp struct {
		Balnaces []SubAccountBalance `json:"balances"`
	}
)

const (
	SubAccountListEndPoint                    = "/sapi/v1/sub-account/list"
	SubAccountTransferHistoryEndPoint         = "/sapi/v1/sub-account/sub/transfer/history"
	SubAccountFuturesInternalTransferEndPoint = "/sapi/v1/sub-account/futures/internalTransfer"
	SubAccountAssetEndPoint                   = "/sapi/v3/sub-account/assets"
)

func NewSubAccountReq() *SubAccountReq {
	return &SubAccountReq{
		RestReq: binance.NewRestReq(),
	}
}

func (sr *SubAccountReq) Email(email string) *SubAccountReq {
	sr.AddFields("email", email)
	return sr
}

func (sr *SubAccountReq) IsFreeze(freeze bool) *SubAccountReq {
	sr.AddFields("isFreeze", freeze)
	return sr
}

func (sr *SubAccountReq) Page(page int) *SubAccountReq {
	sr.AddFields("page", page)
	return sr
}

func (sr *SubAccountReq) Limit(limit int) *SubAccountReq {
	sr.AddFields("limit", limit)
	return sr
}

func NewSubAccountTransferHistoryReq() *SubAccountTransferHistoryReq {
	return &SubAccountTransferHistoryReq{
		RestReq: binance.NewRestReq(),
	}
}

func (sr *SubAccountTransferHistoryReq) FromEmail(email string) *SubAccountTransferHistoryReq {
	sr.AddFields("fromEmail", email)
	return sr
}

func (sr *SubAccountTransferHistoryReq) ToEmail(email string) *SubAccountTransferHistoryReq {
	sr.AddFields("toEmail", email)
	return sr
}

func (sr *SubAccountTransferHistoryReq) StartTime(ts int64) *SubAccountTransferHistoryReq {
	sr.AddFields("startTime", ts)
	return sr
}

func (sr *SubAccountTransferHistoryReq) EndTime(ts int64) *SubAccountTransferHistoryReq {
	sr.AddFields("endTime", ts)
	return sr
}

func (sr *SubAccountTransferHistoryReq) Page(page int) *SubAccountTransferHistoryReq {
	sr.AddFields("page", page)
	return sr
}

func (sr *SubAccountTransferHistoryReq) Limit(limit int) *SubAccountTransferHistoryReq {
	sr.AddFields("limit", limit)
	return sr
}

func NewSubAccountFuturesInternalTransferReq(email string, typ int) *SubAccountFuturesInternalTransferReq {
	sr := binance.NewRestReq()
	sr.AddFields("email", email)
	sr.AddFields("futuresType", typ)
	return &SubAccountFuturesInternalTransferReq{
		RestReq: sr,
	}
}

func (sr *SubAccountFuturesInternalTransferReq) StartTime(ts int64) *SubAccountFuturesInternalTransferReq {
	sr.AddFields("startTime", ts)
	return sr
}

func (sr *SubAccountFuturesInternalTransferReq) EndTime(ts int64) *SubAccountFuturesInternalTransferReq {
	sr.AddFields("endTime", ts)
	return sr
}

func (sr *SubAccountFuturesInternalTransferReq) Limit(lmt int) *SubAccountFuturesInternalTransferReq {
	sr.AddFields("limit", lmt)
	return sr
}

func (sr *SubAccountFuturesInternalTransferReq) Page(page int) *SubAccountFuturesInternalTransferReq {
	sr.AddFields("page", page)
	return sr
}

func NewSubAccountAssetReq(email string) *SubAccountAssetReq {
	rt := binance.NewRestReq()
	rt.AddFields("email", email)
	return &SubAccountAssetReq{
		RestReq: rt,
	}
}

func (rc *RestClient) SubAccountList(ctx context.Context, req *SubAccountReq) (*SubAccountResp, error) {
	var ret SubAccountResp
	if err := rc.GetRequest(ctx, SubAccountListEndPoint, req, true, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (rc *RestClient) SubAccountTransferHistory(ctx context.Context, req *SubAccountTransferHistoryReq) ([]SubAccountTransferHistory, error) {
	var ret []SubAccountTransferHistory
	if err := rc.GetRequest(ctx, SubAccountTransferHistoryEndPoint, req, true, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (rc *RestClient) SubAccountFuturesInternalTransfer(ctx context.Context, req *SubAccountFuturesInternalTransferReq) (*SubAccountFuturesInternalTransferResp, error) {
	var ret SubAccountFuturesInternalTransferResp
	if err := rc.GetRequest(ctx, SubAccountFuturesInternalTransferEndPoint, req, true, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (rc *RestClient) SubAccountAsset(ctx context.Context, req *SubAccountAssetReq) (*SubAccountAssetResp, error) {
	var ret SubAccountAssetResp
	if err := rc.GetRequest(ctx, SubAccountAssetEndPoint, req, true, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
