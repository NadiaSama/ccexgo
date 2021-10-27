package swap

import (
	"context"
	"encoding/json"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	SwapFeeReq struct {
		ContractCode string `json:"contract_code"`
	}

	SwapFee struct {
		Symbol        string `json:"symbol"`
		ContractCode  string `json:"contract_code"`
		OpenMakerFee  string `json:"open_maker_fee"`
		OpenTakerFee  string `json:"open_taker_fee"`
		CloseMakerFee string `json:"close_maker_fee"`
		CloseTakerFee string `json:"close_taker_fee"`
		FeeAsset      string `json:"fee_asset"`
	}

	SwapAccountInfoReq struct {
		contractCode string
	}

	SwapAccountInfo struct {
		Symbol            string          `json:"symbol"`
		ContractCode      string          `json:"contract_code"`
		MarginBalance     decimal.Decimal `json:"margin_balance"`
		MarginFrozen      decimal.Decimal `json:"margin_frozen"`
		MarginAvailable   decimal.Decimal `json:"margin_available"`
		ProfitReal        decimal.Decimal `json:"profit_real"`
		ProfitUnreal      decimal.Decimal `json:"profit_unreal"`
		RiskRate          decimal.Decimal `json:"risk_rate"`
		LiquidationPrice  decimal.Decimal `json:"liquidation_price"`
		WithdrawAvailable decimal.Decimal `json:"withdraw_available"`
		LeverRate         decimal.Decimal `json:"lever_rate"`
		AdjustFactor      decimal.Decimal `json:"adjust_factor"`
	}
)

const (
	SwapFeeEndPoint         = "/swap-api/v1/swap_fee"
	SwapAccountInfoEndPoint = "/swap-api/v1/swap_account_info"
)

func NewSwapFeeReq(symbol string) *SwapFeeReq {
	return &SwapFeeReq{
		ContractCode: symbol,
	}
}

func (rc *RestClient) SwapFee(ctx context.Context, req *SwapFeeReq) ([]SwapFee, error) {
	var ret []SwapFee
	if err := rc.PrivatePostReq(ctx, SwapFeeEndPoint, req, &ret); err != nil {
		return nil, errors.WithMessage(err, "request swap_fee fail")
	}

	return ret, nil
}

func NewSwapAccountInfoReq() *SwapAccountInfoReq {
	return &SwapAccountInfoReq{}
}

func (sai *SwapAccountInfoReq) ContractCode(code string) *SwapAccountInfoReq {
	sai.contractCode = code
	return sai
}

func (sai *SwapAccountInfoReq) Serialize() ([]byte, error) {
	m := map[string]string{}
	if sai.contractCode != "" {
		m["contract_code"] = sai.contractCode
	}
	return json.Marshal(m)
}

func (rc *RestClient) SwapAccountInfo(ctx context.Context, req *SwapAccountInfoReq) ([]SwapAccountInfo, error) {
	var ret []SwapAccountInfo
	if err := rc.PrivatePostReq(ctx, SwapAccountInfoEndPoint, req, &ret); err != nil {
		return nil, errors.WithMessage(err, "request swap_account_info fail")
	}
	return ret, nil
}

func (rc *RestClient) FetchBalance(ctx context.Context, currencies ...string) (*exchange.Balances, error) {
	accountInfo, err := rc.SwapAccountInfo(ctx, NewSwapAccountInfoReq())
	if err != nil {
		return nil, errors.WithMessage(err, "fetch swap account info fail")
	}

	bm := map[string]*exchange.Balance{}

	for i := range accountInfo {
		a := accountInfo[i]
		bm[a.Symbol] = &exchange.Balance{
			Currency: a.Symbol,
			Total:    a.MarginBalance,
			Frozen:   a.MarginFrozen,
			Free:     a.MarginAvailable,
		}
	}
	balances := exchange.NewBalances()
	balances.Raw = accountInfo

	if len(currencies) != 0 {
		for _, c := range currencies {
			b, ok := bm[exchange.CurrencyFormat(c)]
			if !ok {
				return nil, errors.Errorf("currency '%s' not support", c)
			}
			balances.Add(b)
		}
		return balances, nil
	}

	for _, val := range bm {
		balances.Add(val)
	}
	return balances, nil
}
