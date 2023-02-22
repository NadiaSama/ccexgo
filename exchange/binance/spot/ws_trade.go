package spot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

type (
	TradeNotify struct {
		Event       string `json:"e"`
		EventTS     int64  `json:"E"`
		Symbol      string `json:"s"`
		TradeID     int64  `json:"t"`
		Price       string `json:"p"`
		Quantity    string `json:"q"`
		BuyOrderID  int64  `json:"b"`
		SellOrderID int64  `json:"a"`
		TradeTS     int64  `json:"T"`
		Taker       bool   `json:"m"`
	}

	TradeChannel struct {
		sym string
	}
)

func NewTradeChannel(symbol string) *TradeChannel {
	return &TradeChannel{
		sym: strings.ToLower(symbol),
	}
}

func (tc *TradeChannel) String() string {
	return fmt.Sprintf("%s@trade", tc.sym)
}

func ParseTradeNotify(gs *gjson.Result) *TradeNotify {
	return &TradeNotify{
		Event:       gs.Get("e").String(),
		EventTS:     gs.Get("E").Int(),
		Symbol:      gs.Get("s").String(),
		TradeID:     gs.Get("t").Int(),
		Price:       gs.Get("p").String(),
		Quantity:    gs.Get("q").String(),
		BuyOrderID:  gs.Get("b").Int(),
		SellOrderID: gs.Get("a").Int(),
		TradeTS:     gs.Get("T").Int(),
		Taker:       gs.Get("m").Bool(),
	}
}

func (tn *TradeNotify) Parse() ([]*exchange.Trade, error) {
	sym, err := ParseSymbol(tn.Symbol)
	if err != nil {
		return nil, errors.WithMessage(err, "invalid symbol")
	}

	price, err := decimal.NewFromString(tn.Price)
	if err != nil {
		return nil, errors.WithMessage(err, "invalid price")
	}
	amount, err := decimal.NewFromString(tn.Quantity)
	if err != nil {
		return nil, errors.WithMessage(err, "invalid amount")
	}

	buy := &exchange.Trade{
		ID:      strconv.FormatInt(tn.TradeID, 10),
		OrderID: strconv.FormatInt(tn.BuyOrderID, 10),
		Symbol:  sym,
		Price:   price,
		Amount:  amount,
		Side:    exchange.OrderSideBuy,
		Time:    binance.Milli2Time(tn.TradeTS),
		IsMaker: tn.Taker,
		Raw:     tn,
	}

	sell := &exchange.Trade{
		ID:      strconv.FormatInt(tn.TradeID, 10),
		OrderID: strconv.FormatInt(tn.SellOrderID, 10),
		Symbol:  sym,
		Price:   price,
		Amount:  amount,
		Side:    exchange.OrderSideBuy,
		Time:    binance.Milli2Time(tn.TradeTS),
		IsMaker: !tn.Taker,
		Raw:     tn,
	}

	return []*exchange.Trade{buy, sell}, nil
}
