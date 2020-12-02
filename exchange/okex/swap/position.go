package swap

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Position struct {
		LiquidationPrice decimal.Decimal `json:"liquidation_price"`
		Position         decimal.Decimal `json:"position"`
		AvailPosition    decimal.Decimal `json:"avail_position"`
		Margin           decimal.Decimal `json:"margin"`
		AvgCost          decimal.Decimal `json:"avg_cost"`
		SettlementPrice  decimal.Decimal `json:"settlement_price"`
		InstrumentID     string          `json:"instrument_id"`
		Leverage         decimal.Decimal `json:"leverage"`
		RealizedPNL      decimal.Decimal `json:"realized_pnl"`
		Side             string          `json:"side"`
		Timestamp        string          `json:"timestamp"`
		MaintMarginRatio decimal.Decimal `json:"maint_margin_ratio"`
		SettlePNL        decimal.Decimal `json:"settle_pnl"`
		Last             decimal.Decimal `json:"last"`
		UnRealizedPNL    decimal.Decimal `json:"unrealized_pnl"`
	}

	MarginPosition struct {
		MarginMode string     `json:"margin_mode"`
		Timestamp  string     `json:"timestamp"`
		Holding    []Position `json:"holding"`
	}
)

const (
	MarginModeFixed   = "fixed"
	MarginModeCrossed = "crossed"
)

var (
	positionMap map[string]exchange.PositionSide = map[string]exchange.PositionSide{
		"long":  exchange.PositionSideLong,
		"short": exchange.PositionSideShort,
	}
	modeMap map[string]exchange.PositionMode = map[string]exchange.PositionMode{
		"fixed":   exchange.PositionModeFixed,
		"crossed": exchange.PositionModeCross,
	}
)

func (rc *RestClient) FetchPosition(ctx context.Context, sym ...exchange.Symbol) ([]*exchange.Position, error) {
	if len(sym) > 1 {
		return nil, errors.Errorf("only 1 symbol is support")
	}
	var mps []MarginPosition
	var uri string

	if len(sym) == 1 {
		uri = fmt.Sprintf("/api/swap/v3/%s/position", sym[0].String())
	} else {
		uri = "/api/swap/v3/position"
	}

	if err := rc.Request(ctx, http.MethodGet, uri, nil, nil, true, &mps); err != nil {
		return nil, errors.WithMessage(err, "okex request error")
	}

	ret := []*exchange.Position{}
	for _, mp := range mps {
		posMode, ok := modeMap[mp.MarginMode]
		if !ok {
			return nil, errors.Errorf("unkown margin_mode '%s'", mp.MarginMode)
		}

		for _, pos := range mp.Holding {
			p, err := pos.Transform(posMode)
			if err != nil {
				return nil, err
			}
			ret = append(ret, p)
		}
	}
	return ret, nil
}

func (pos *Position) Transform(posMode exchange.PositionMode) (*exchange.Position, error) {
	sym, err := okex.ParseSwapSymbol(pos.InstrumentID)
	if err != nil {
		return nil, err
	}

	side, ok := positionMap[pos.Side]
	if !ok {
		return nil, errors.Errorf("unkown position side '%s'", pos.Side)
	}
	ct, err := okex.ParseTime(pos.Timestamp)
	if err != nil {
		return nil, errors.WithMessagef(err, "invalid okex timestmap %s", pos.Timestamp)
	}

	return &exchange.Position{
		Symbol:           sym,
		Mode:             posMode,
		Side:             side,
		LiquidationPrice: pos.LiquidationPrice,
		AvgOpenPrice:     pos.AvgCost,
		CreateTime:       ct,
		Margin:           pos.Margin,
		MarginMaintRatio: pos.MaintMarginRatio,
		Position:         pos.Position,
		AvailPosition:    pos.AvailPosition,
		RealizedPNL:      pos.RealizedPNL,
		UNRealizedPNL:    pos.UnRealizedPNL,
		Leverage:         pos.Leverage,
		Raw:              pos,
	}, nil
}
