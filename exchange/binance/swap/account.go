package swap

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	AccountAsset struct {
		Asset                  string          `json:"asset"`
		WalletBalance          decimal.Decimal `json:"walletBalance"`
		UnrealizedProfit       decimal.Decimal `json:"unrealizedProfit"`
		MarginBalance          decimal.Decimal `json:"marginBalance"`
		InitialMargin          decimal.Decimal `json:"initialMargin"`
		PositionInitialMargin  decimal.Decimal `json:"positionInitialMargin"`
		OpenOrderInitialMargin decimal.Decimal `json:"openOrderInitialMargin"`
		CrossWalletBalance     decimal.Decimal `json:"crossWalletBalance"`
		CrossUnPNL             decimal.Decimal `json:"crossUnPnl"`
		AvailableBalance       decimal.Decimal `json:"availableBalance"`
		MaxWithdrawAmount      decimal.Decimal `json:"maxWithdrawAmount"`
		MarginAvailable        bool            `json:"marginAvailable"`
		UpdateTime             int64           `json:"updateTime"`
	}

	AccountPosition struct {
		Symbol                 string          `json:"symbol"`
		InitialMargin          decimal.Decimal `json:"initialMargin"`
		MaintMargin            decimal.Decimal `json:"maintMargin"`
		UnrealizedProfit       decimal.Decimal `json:"unrealizedProfit"`
		PositionInitialMargin  decimal.Decimal `json:"positionInitialMargin"`
		OpenOrderInitialMargin decimal.Decimal `json:"openOrderInitialMargin"`
		Leverage               decimal.Decimal `json:"leverage"`
		Isolated               bool            `json:"isolated"`
		EntryPrice             decimal.Decimal `json:"entryPrice"`
		MaxNotional            decimal.Decimal `json:"maxNotional"`
		PositionSide           string          `json:"positionSide"`
		PositionAmt            decimal.Decimal `json:"positionAmt"`
		UpdateTime             int64           `json:"updateTime"`
	}

	AccountResp struct {
		binance.APIError                             //in case of error
		FeeTier                    int               `json:"feeTier"`
		CanTrade                   bool              `json:"canTrade"`
		CanDeposit                 bool              `json:"canDeposit"`
		CanWithdraw                bool              `json:"canWithdraw"`
		UpdateTime                 int64             `json:"updateTime"`
		TotalInitialMargin         decimal.Decimal   `json:"totalInitialMargin"`
		TotalMaintMargin           decimal.Decimal   `json:"totalMaintMargin"`
		TotalWalletBalance         decimal.Decimal   `json:"totalWalletBalance"`
		TotalUnrealizedProfit      decimal.Decimal   `json:"totalUnrealizedPorfit"`
		TotalMarginBalance         decimal.Decimal   `json:"totalMarginBalance"`
		TotalPositionInitialMargin decimal.Decimal   `json:"totalPositionInitialMargin"`
		Assets                     []AccountAsset    `json:"assets"`
		Positions                  []AccountPosition `json:"positions"`
	}

	AccountReq struct {
		*binance.RestReq
	}
)

const (
	AccountEndPoint = "/fapi/v2/account"
)

func NewAccountReq() *AccountReq {
	return &AccountReq{
		RestReq: binance.NewRestReq(),
	}
}

func (rc *RestClient) Account(ctx context.Context, req *AccountReq) (*AccountResp, error) {
	var ret AccountResp
	if err := rc.GetRequest(ctx, AccountEndPoint, req, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "get account fail")
	}

	return &ret, nil
}

func (ar *AccountResp) GetAsset(currency string) (*AccountAsset, error) {
	for i := range ar.Assets {
		as := &ar.Assets[i]
		if as.Asset == currency {
			return as, nil
		}
	}

	return nil, errors.Errorf("unknown currency=%s", currency)
}

func (ar *AccountResp) GetPosition(symbol string) ([]AccountPosition, error) {
	ret := []AccountPosition{}
	for _, pos := range ar.Positions {
		if pos.Symbol == symbol {
			ret = append(ret, pos)
		}
	}

	if len(ret) == 0 {
		return nil, errors.Errorf("unknown position=%s", symbol)
	}
	return ret, nil
}
