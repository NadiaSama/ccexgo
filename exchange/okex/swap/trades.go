package swap

import (
	"context"
	"net/http"
	"strconv"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Fill struct {
		TradeID      string          `json:"trade_id"`
		FillID       string          `json:"fill_id"`
		InstrumentID string          `json:"instrument_id"`
		OrderID      string          `json:"order_id"`
		Price        decimal.Decimal `json:"price"`
		OrderQty     decimal.Decimal `json:"order_qty"`
		Fee          decimal.Decimal `json:"fee"`
		Timestamp    string          `json:"timestamp"`
		ExecType     string          `json:"exec_type"`
		Side         string          `json:"side"`
		OrderSide    string          `json:"order_side"`
		Type         string          `json:"type"`
	}
)

const (
	FillsEndPoint = "/api/swap/v3/fills"
)

func (rc *RestClient) Fills(ctx context.Context, instrumentID string, orderID string, before, after, limit string) ([]Fill, error) {
	values := okex.FillsParam(instrumentID, orderID, before, after, limit)
	var ret []Fill
	if err := rc.Request(ctx, http.MethodGet, FillsEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch fills fail")
	}
	return ret, nil
}

func (rc *RestClient) Trades(ctx context.Context, req *exchange.TradeReqParam) ([]*exchange.Trade, error) {
	fills, err := rc.Fills(ctx, req.Symbol.String(), "", req.StartID, req.EndID, strconv.Itoa(req.Limit))
	if err != nil {
		return nil, err
	}

	var ret []*exchange.Trade
	for i := range fills {
		fill := fills[i]
		trade, err := fill.Parse()
		if err != nil {
			return nil, errors.WithMessage(err, "parse fill error")
		}
		ret = append(ret, trade)
	}
	return ret, nil
}

func (f *Fill) Parse() (*exchange.Trade, error) {
	s, err := ParseSymbol(f.InstrumentID)
	if err != nil {
		return nil, err
	}
	t, err := okex.ParseTime(f.Timestamp)
	if err != nil {
		return nil, err
	}

	var side exchange.OrderSide
	if f.Side == "long" {
		if f.OrderSide == "buy" {
			side = exchange.OrderSideBuy
		} else if f.OrderSide == "sell" {
			side = exchange.OrderSideCloseLong
		} else {
			return nil, errors.Errorf("unkown order side '%s'", f.OrderSide)
		}
	} else if f.Side == "short" {
		if f.OrderSide == "buy" {
			side = exchange.OrderSideCloseShort
		} else if f.OrderSide == "sell" {
			side = exchange.OrderSideSell
		} else {
			return nil, errors.Errorf("unkown order side '%s'", f.OrderSide)
		}
	} else {
		return nil, errors.Errorf("unkown side '%s'", f.Side)
	}

	return &exchange.Trade{
		ID:          f.TradeID,
		OrderID:     f.OrderID,
		Price:       f.Price,
		Amount:      f.OrderQty,
		Fee:         f.Fee,
		FeeCurrency: "USDT",
		Symbol:      s,
		Time:        t,
		Side:        side,
		Raw:         *f,
	}, nil
}
