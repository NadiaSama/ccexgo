package ftx

import (
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Order struct {
		CreatedAt     string          `json:"createdAt"`
		FilledSize    decimal.Decimal `json:"filledSize"`
		Future        string          `json:"future"`
		ID            int64           `json:"id"`
		Market        string          `json:"market"`
		Price         decimal.Decimal `json:"price"`
		AvgFillPrice  decimal.Decimal `json:"avgFillPrice"`
		RemainingSize decimal.Decimal `json:"remainingSize"`
		Side          string          `json:"side"`
		Size          decimal.Decimal `json:"size"`
		Status        string          `json:"status"`
		Type          string          `json:"type"`
		ReduceOnly    bool            `json:"reduceOnly"`
		IOC           bool            `json:"ioc"`
		PostOnly      bool            `json:"postOnly"`
		ClientID      string          `json:"clientId"`
	}
)

const (
	ftxOrderNew   = "new"
	ftxOrderOpen  = "open"
	ftxOrderClose = "closed"
)

var (
	typeMap map[string]exchange.OrderType = map[string]exchange.OrderType{
		"limit":  exchange.OrderTypeLimit,
		"market": exchange.OrderTypeMarket,
	}

	sideMap map[string]exchange.OrderSide = map[string]exchange.OrderSide{
		"buy":  exchange.OrderSideBuy,
		"sell": exchange.OrderSideSell,
	}
)

func (rc *RestClient) parseOrder(o *Order, isFutures bool) (*exchange.Order, error) {
	ct, err := time.Parse("2006-01-02T15:04:05.000000Z", o.CreatedAt)
	if err != nil {
		return nil, errors.WithMessagef(err, "bad create time '%s'", o.CreatedAt)
	}
	var os exchange.OrderStatus
	if o.Status == ftxOrderNew || o.Status == ftxOrderOpen {
		os = exchange.OrderStatusOpen
	} else {
		if o.FilledSize == o.Size {
			os = exchange.OrderStatusDone
		} else {
			os = exchange.OrderStatusCancel
		}
	}
	var symbol exchange.Symbol
	if isFutures {
		s, err := rc.ParseFutureSymbol(o.Market)
		if err != nil {
			return nil, err
		}
		symbol = s
	} else {
		s, err := rc.ParseSwapSymbol(o.Market)
		if err != nil {
			return nil, err
		}
		symbol = s
	}

	typ, ok := typeMap[o.Type]
	if !ok {
		return nil, errors.Errorf("unkown order type '%s'", o.Type)
	}

	side, ok := sideMap[o.Side]
	if !ok {
		return nil, errors.Errorf("unkown order side '%s'", o.Side)
	}

	order := exchange.Order{
		ID:       exchange.NewIntID(o.ID),
		Symbol:   symbol,
		Amount:   o.Size,
		Filled:   o.FilledSize,
		Price:    o.Price,
		AvgPrice: o.AvgFillPrice,
		Created:  ct,
		Updated:  ct,
		Status:   os,
		Side:     side,
		Type:     typ,
	}
}
