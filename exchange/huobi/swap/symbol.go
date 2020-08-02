package swap

import (
	"context"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
	"github.com/pkg/errors"
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

func newSymbol(index string) *Symbol {
	return &Symbol{
		exchange.NewBaseSwapSymbol(index),
	}
}

func (s *Symbol) String() string {
	return s.Index()
}

//TODO ParseSwapSymbol? ParseSymbol?
func (rc *RestClient) initSymbol(ctx context.Context) error {
	var ci ContractInfo
	if err := rc.Request(ctx, http.MethodGet, ContractEndPoint, nil, nil, false, &ci); err != nil {
		return errors.WithMessagef(err, "get contract info fail")
	}

	if ci.Status != huobi.StatusOK {
		return errors.Errorf("got huobi contract info fail")
	}

	for _, data := range ci.Data {
		rc.swapCodeMap[data.Symbol] = newSymbol(data.ContractCode)
	}
	return nil
}

func (rc *RestClient) GetSwapContract(symbol string) (*Symbol, error) {
	val, ok := rc.swapCodeMap[symbol]
	if !ok {
		return nil, errors.Errorf("unkown symbol '%s'", symbol)
	}
	return val, nil
}
