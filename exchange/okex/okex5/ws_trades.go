package okex5

import (
	"encoding/json"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Trade struct {
		InstID  string    `json:"instId"`
		TradeID string    `json:"tradeId"`
		Px      string    `json:"px"`
		Sz      string    `json:"sz"`
		Side    OrderSide `json:"side"`
		Ts      string    `json:"ts"`
	}
)

const (
	TradesChannel = "trades"
)

func init() {
	parseCBMap[TradesChannel] = parseTrades
}

func NewTradesChannel(instID string) *Okex5Channel {
	return &Okex5Channel{
		Channel: TradesChannel,
		InstID:  instID,
	}
}

func parseTrades(data *wsResp) (*rpc.Notify, error) {
	var t []Trade
	if err := json.Unmarshal(data.Data, &t); err != nil {
		return nil, err
	}

	trades := []*exchange.Trade{}
	for _, i := range t {
		trade, err := i.Parse()
		if err != nil {
			return nil, errors.WithMessage(err, "parse trade fail")
		}

		trades = append(trades, trade)
	}

	return &rpc.Notify{
		Method: data.Arg.Channel,
		Params: trades,
	}, nil
}

func (t *Trade) Parse() (*exchange.Trade, error) {
	sym, err := ParseSymbol(t.InstID)
	if err != nil {
		return nil, errors.WithMessage(err, "parse symbol fail")
	}

	prc, err := decimal.NewFromString(t.Px)
	if err != nil {
		return nil, errors.WithMessage(err, "parse price fail")
	}

	amt, err := decimal.NewFromString(t.Sz)
	if err != nil {
		return nil, errors.WithMessage(err, "parse size fail")
	}

	ts, err := ParseTimestamp(t.Ts)
	if err != nil {
		return nil, errors.WithMessage(err, "parse ts fail")
	}

	var side exchange.OrderSide
	if t.Side == OrderSideBuy {
		side = exchange.OrderSideBuy
	} else if t.Side == OrderSideSell {
		side = exchange.OrderSideSell
	}

	return &exchange.Trade{
		ID:     t.TradeID,
		Symbol: sym,
		Price:  prc,
		Amount: amt,
		Time:   ts,
		Side:   side,
		Raw:    t,
	}, nil
}
