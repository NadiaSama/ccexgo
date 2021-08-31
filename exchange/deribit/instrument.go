package deribit

import (
	"context"
	"net/http"
	"net/url"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	InstrumentsRequest struct {
		currency string
		kind     string
		expired  bool
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

	for _, ir := range irs {
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
