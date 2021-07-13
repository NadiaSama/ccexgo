package swap

import (
	"context"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Symbol struct {
		*exchange.BaseSwapSymbol
	}

	Data struct {
		Symbol         string  `json:"symbol"`
		ContractCode   string  `json:"contract_code"`
		ContractSize   float64 `json:"contract_size"`
		PriceTick      float64 `json:"price_tick"`
		CreateDate     string  `json:"create_date"`
		ContractStatus int     `json:"contract_status"`
		SettlementDate string  `json:"settlement_date"`
	}

	ContractInfo struct {
		Status string `json:"status"`
		TS     int    `json:"ts"`
		Data   []Data `json:"data"`
	}
)

const (
	ContractEndPoint = "/swap-api/v1/swap_contract_info"
)

var (
	contractMap = map[string]*Symbol{}
)

func (s *Symbol) String() string {
	return s.Index()
}

func Init(ctx context.Context) error {
	rc := NewRestClient("", "")

	var ci ContractInfo
	if err := rc.RequestWithRawResp(ctx, http.MethodGet, ContractEndPoint, nil, nil, false, &ci); err != nil {
		return errors.WithMessagef(err, "get contract info fail")
	}

	if ci.Status != huobi.StatusOK {
		return errors.Errorf("got huobi contract info fail")
	}

	for i := range ci.Data {
		data := ci.Data[i]
		cv := decimal.NewFromFloat(data.ContractSize)
		symbol := &Symbol{
			exchange.NewBaseSwapSymbolWithCfg(data.ContractCode, cv, exchange.SymbolConfig{
				AmountPrecision: decimal.NewFromFloat(1.0),
				PricePrecision:  decimal.NewFromFloat(data.PriceTick),
				AmountMin:       decimal.NewFromInt(1),
				AmountMax:       decimal.Zero,
			}, &data),
		}

		contractMap[data.ContractCode] = symbol
	}

	return nil
}

func ParseSymbol(sym string) (exchange.SwapSymbol, error) {
	s, ok := contractMap[sym]
	if !ok {
		return nil, errors.Errorf("unsupport symbol %s", sym)
	}

	return s, nil
}
