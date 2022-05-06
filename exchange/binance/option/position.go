package option

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Position struct {
		EntryPrice         decimal.Decimal `json:"entryPrice"`
		Symbol             string          `json:"symbol"`
		Side               string          `json:"side"`
		Leverage           int             `json:"leverage"`
		Quantity           decimal.Decimal `json:"quantity"`
		ReducibleQty       decimal.Decimal `json:"reducibleQty"`
		MarkValue          decimal.Decimal `json:"markValue"`
		AutoReducePriority int             `json:"autoReducePriority"`
		Ror                decimal.Decimal `json:"ror"`
		UnrealizedPNL      decimal.Decimal `json:"unrealizedPNL"`
		MarkPrice          decimal.Decimal `json:"markPrice"`
		StrikePrice        decimal.Decimal `json:"strikePrice"`
		ExpiryDate         int64           `json:"expiryDate"`
	}

	PositionReq struct {
		*binance.RestReq
	}
)

const (
	PositionEndPoint  = "/vapi/v1/position"
	PositionSideShort = "SHORT"
	PositionSideLong  = "LONG"
)

func NewPositionReq() *PositionReq {
	return &PositionReq{
		RestReq: binance.NewRestReq(),
	}
}

func (pr *PositionReq) Symbol(sym string) *PositionReq {
	pr.AddFields("symbol", sym)
	return pr
}

func (rc *RestClient) Position(ctx context.Context, req *PositionReq) ([]Position, error) {
	var ret []Position

	if err := rc.GetRequest(ctx, PositionEndPoint, req, true, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (rc *RestClient) FetchPosition(ctx context.Context, sym ...exchange.Symbol) ([]exchange.Position, error) {
	if len(sym) != 0 && len(sym) != 1 {
		return nil, errors.Errorf("at most 1 symbol is support")
	}

	req := NewPositionReq()
	if len(sym) != 0 {
		req.Symbol(sym[0].String())
	}

	pos, err := rc.Position(ctx, req)
	if err != nil {
		return nil, errors.WithMessage(err, "fetch position fail")
	}

	ret := make([]exchange.Position, len(pos))
	for i, p := range pos {
		sp, err := p.Transfer()
		if err != nil {
			return nil, errors.WithMessage(err, "parse position fail")
		}

		ret[i] = *sp
	}
	return ret, nil
}

func (p *Position) Transfer() (*exchange.Position, error) {
	var posSide exchange.PositionSide

	if p.Side == PositionSideLong {
		posSide = exchange.PositionSideLong
	} else if p.Side == PositionSideShort {
		posSide = exchange.PositionSideShort
	} else {
		return nil, errors.Errorf("unknown side='%s'", p.Side)
	}

	sym, err := ParseSymbol(p.Symbol)
	if err != nil {
		return nil, errors.WithMessage(err, "parse symbol fail")
	}

	return &exchange.Position{
		Symbol:        sym,
		Side:          posSide,
		AvgOpenPrice:  p.EntryPrice,
		Position:      p.Quantity.Abs(),
		AvailPosition: p.ReducibleQty.Abs(),
		UNRealizedPNL: p.UnrealizedPNL,
		RealizedPNL:   p.Ror,
		Leverage:      decimal.NewFromInt(int64(p.Leverage)),
		Raw:           p,
	}, nil
}
