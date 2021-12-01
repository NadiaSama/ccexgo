package swap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/NadiaSama/ccexgo/internal/rpc"
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

	Trade struct {
		InstrumentID string `json:"instrument_id"`
		Price        string `json:"price"`
		Side         string `json:"side"`
		Size         string `json:"size"`
		Timestamp    string `json:"timestamp"`
		TradeID      string `json:"trade_id"`
	}
	TradeChannel struct {
		exchange.SwapSymbol
	}
)

const (
	FillsEndPoint = "/api/swap/v3/fills"
	ExecTypeTaker = "T"
	ExecTypeMaker = "M"

	TradeTable = "swap/trade"
)

func init() {
	okex.SubscribeCB(TradeTable, parseTradeCB)
}

func NewTradeChannel(symbol exchange.SwapSymbol) *TradeChannel {
	return &TradeChannel{
		SwapSymbol: symbol,
	}
}

func (s *TradeChannel) String() string {
	return fmt.Sprintf("%s:%s", TradeTable, s.SwapSymbol)
}

func parseTradeCB(table string, action string, data json.RawMessage) (*rpc.Notify, error) {
	var trades []Trade
	if err := json.Unmarshal(data, &trades); err != nil {
		return nil, errors.WithMessage(err, "parse trades fail")
	}

	ret := make([]exchange.PublicTrade, 0, len(trades))
	for _, t := range trades {
		pt, err := t.Transform()
		if err != nil {
			return nil, err
		}

		ret = append(ret, *pt)
	}
	return &rpc.Notify{
		Method: table,
		Params: ret,
	}, nil
}

func (t *Trade) Transform() (*exchange.PublicTrade, error) {
	sym, err := ParseSymbol(t.InstrumentID)
	if err != nil {
		return nil, err
	}

	price, err := decimal.NewFromString(t.Price)
	if err != nil {
		return nil, errors.WithMessagef(err, "parse price fail price='%s'", t.Price)
	}

	amount, err := decimal.NewFromString(t.Size)
	if err != nil {
		return nil, errors.WithMessagef(err, "parse size fail size='%s'", t.Size)
	}

	var side exchange.OrderSide
	if t.Side == "buy" {
		side = exchange.OrderSideBuy
	} else if t.Side == "sell" {
		side = exchange.OrderSideSell
	} else {
		return nil, errors.Errorf("unkown orderside '%s'", t.Side)
	}

	ts, err := okex.ParseTime(t.Timestamp)
	if err != nil {
		return nil, err
	}

	return &exchange.PublicTrade{
		Symbol: sym,
		Price:  price,
		Amount: amount,
		Side:   side,
		ID:     t.TradeID,
		Time:   ts,
		Raw:    t,
	}, nil
}

func (rc *RestClient) Fills(ctx context.Context, instrumentID string, orderID string, before, after, limit string) ([]Fill, error) {
	values := okex.FillsParam(instrumentID, orderID, before, after, limit)
	var ret []Fill
	if err := rc.Request(ctx, http.MethodGet, FillsEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch fills fail")
	}
	return ret, nil
}

func (rc *RestClient) Trades(ctx context.Context, req *exchange.TradeReqParam) ([]exchange.Trade, error) {
	fills, err := rc.Fills(ctx, req.Symbol.String(), "", req.StartID, req.EndID, strconv.Itoa(req.Limit))
	if err != nil {
		return nil, err
	}

	var ret []exchange.Trade
	for i := range fills {
		fill := fills[i]
		trade, err := fill.Parse()
		if err != nil {
			return nil, errors.WithMessage(err, "parse fill error")
		}
		ret = append(ret, *trade)
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

	var isMaker bool
	if f.ExecType == ExecTypeTaker {
		isMaker = false
	} else if f.ExecType == ExecTypeMaker {
		isMaker = true
	} else {
		return nil, errors.Errorf("unkown execType '%s'", f.ExecType)
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
		IsMaker:     isMaker,
		Raw:         *f,
	}, nil
}
