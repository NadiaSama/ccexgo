package spot

import (
	"encoding/json"
	"fmt"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Trade struct {
		InstrumentID string `json:"instrument_id"`
		Price        string `json:"price"`
		Side         string `json:"side"`
		Size         string `json:"size"`
		Timestamp    string `json:"timestamp"`
		TradeID      string `json:"trade_id"`
	}

	TradeChannel struct {
		symbol exchange.Symbol
	}
)

const (
	TradeTable = "spot/trade"
)

func init() {
	okex.SubscribeCB(TradeTable, parseTradeCB)
}

func NewTradeChannel(sym exchange.Symbol) *TradeChannel {
	return &TradeChannel{
		symbol: sym,
	}
}

func (s *TradeChannel) String() string {
	return fmt.Sprintf("%s:%s", TradeTable, s.symbol.String())
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
