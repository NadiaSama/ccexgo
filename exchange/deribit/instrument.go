package deribit

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	InstrumentsRequest struct {
		currency string
		kind     string
		expired  bool
	}

	InstrumentResult struct {
		TickSize            decimal.Decimal `json:"tick_size"`
		TakerCommision      decimal.Decimal `json:"taker_commision"`
		MakerCommision      decimal.Decimal `json:"maker_commision"`
		Strike              decimal.Decimal `json:"strike"`
		SettlementPeriod    string          `json:"settlement_period"`
		QuoteCurrency       string          `json:"quote_currency"`
		BaseCurreny         string          `json:"base_currency"`
		MinTradeAmount      decimal.Decimal `json:"min_trade_amount"`
		Kind                string          `json:"kind"`
		IsActive            bool            `json:"is_active"`
		InstrumentName      string          `json:"instrument_name"`
		ExpirationTimestamp int64           `json:"expiration_timestamp"`
		CreationTimestamp   int64           `json:"creation_timestamp"`
		ContractSize        decimal.Decimal `json:"contract_size"`
		OptionType          string          `json:"option_type"`
	}
)

const (
	InstrumentsKindOption = "option"
	InstrumentsKindFuture = "future"

	InstrumentsEndPoint = "/public/get_instruments"
)

func NewInstrumentsRequest(currency string) *InstrumentsRequest {
	return &InstrumentsRequest{
		currency: currency,
	}
}

func (ir *InstrumentsRequest) Expired() *InstrumentsRequest {
	ir.expired = true
	return ir
}

func (ir *InstrumentsRequest) Kind(kind string) *InstrumentsRequest {
	ir.kind = kind
	return ir
}

func (c *RestClient) Instruments(ctx context.Context, ir *InstrumentsRequest) ([]InstrumentResult, error) {
	values := url.Values{}
	values.Add("currency", ir.currency)
	if ir.kind != "" {
		values.Add("kind", ir.kind)
	}

	if ir.expired {
		values.Add("expired", "true")
	}

	var ret []InstrumentResult
	if err := c.Request(ctx, http.MethodGet, InstrumentsEndPoint, values, nil, false, &ret); err != nil {
		return nil, errors.WithMessage(err, "get instruments fail")
	}
	return ret, nil
}

func (c *RestClient) OptionSymbols(ctx context.Context, currency string) ([]exchange.OptionSymbol, error) {
	req := NewInstrumentsRequest(currency).Kind(InstrumentsKindOption)
	irs, err := c.Instruments(ctx, req)
	if err != nil {
		return nil, err
	}

	ret := make([]exchange.OptionSymbol, 0)

	for i := range irs {
		ir := irs[i]
		sym, err := ir.Symbol()
		if err != nil {
			return nil, err
		}

		osym, ok := sym.(exchange.OptionSymbol)
		if !ok {
			return nil, errors.Errorf("invalid symbol %+v", sym)
		}

		ret = append(ret, osym)
	}
	return ret, nil
}

//Symbols fetch deribit option and futures(future + swap) symbols
func (c *RestClient) Symbols(ctx context.Context, currency string) ([]exchange.Symbol, error) {
	req := NewInstrumentsRequest(currency)
	irs, err := c.Instruments(ctx, req)
	if err != nil {
		return nil, err
	}

	ret := make([]exchange.Symbol, 0, len(irs))
	for i := range irs {
		ir := irs[i]
		sym, err := ir.Symbol()
		if err != nil {
			return nil, err
		}

		ret = append(ret, sym)
	}
	return ret, nil
}

func (i *InstrumentResult) Symbol() (exchange.Symbol, error) {
	cfg := exchange.SymbolConfig{
		AmountMax:       decimal.Zero,
		AmountMin:       i.MinTradeAmount,
		PricePrecision:  i.TickSize,
		AmountPrecision: i.MinTradeAmount,
		ValueMin:        decimal.Zero,
	}
	var st time.Time
	if i.Kind != KindFuture || i.SettlementPeriod != SettlePeriodPerpetual {
		st = time.Unix(i.ExpirationTimestamp/1e3, 0)
	}
	if i.Kind == KindOption {
		var t exchange.OptionType
		if i.OptionType == OptionTypeCall {
			t = exchange.OptionTypeCall
		} else if i.OptionType == OptionTypePut {
			t = exchange.OptionTypePut
		} else {
			return nil, errors.Errorf("unkown option type %s", i.OptionType)
		}
		ret := &OptionSymbol{
			exchange.NewBaseOptionSymbol(i.BaseCurreny, st, i.Strike, t, cfg, i),
		}
		return ret, nil

	} else if i.Kind == KindFuture {
		if i.SettlementPeriod == SettlePeriodPerpetual {
			return &SwapSymbol{
				exchange.NewBaseSwapSymbolWithCfg(i.BaseCurreny, i.ContractSize, cfg, i),
			}, nil
		}

		var ft exchange.FutureType
		diff := time.Until(st)
		if i.SettlementPeriod == SettlePeriodWeek {
			if (diff / time.Second / 86400) < 7 {
				ft = exchange.FutureTypeCW
			} else {
				ft = exchange.FutureTypeNW
			}
		} else {
			if (diff / time.Second / 86400) < 31 {
				ft = exchange.FutureTypeCQ
			} else if (diff / time.Second / 86400) < 90 {
				ft = exchange.FutureTypeNQ
			} else {
				ft = exchange.FutureTypeNNQ
			}
		}
		return &FuturesSymbol{
			exchange.NewBaseFuturesSymbolWithCfg(i.BaseCurreny, st, ft, cfg, i),
		}, nil
	}
	return nil, errors.Errorf("unkown kind '%s'", i.Kind)
}
