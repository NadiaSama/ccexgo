package swap

import (
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var (
	statusMap = map[int]exchange.OrderStatus{
		1:  exchange.OrderStatusOpen,
		2:  exchange.OrderStatusOpen,
		3:  exchange.OrderStatusOpen,
		4:  exchange.OrderStatusOpen,
		5:  exchange.OrderStatusOpen,
		6:  exchange.OrderStatusDone,
		7:  exchange.OrderStatusCancel,
		11: exchange.OrderStatusOpen,
	}
)

func ParseOrder(resp *Response) (*exchange.Order, error) {
	symbol, err := ParseSymbol(resp.ContractCode)
	if err != nil {
		return nil, err
	}

	var created time.Time
	if resp.CreatedAt != 0 {
		created = huobi.ParseTS(resp.CreatedAt)
	}

	st, ok := statusMap[resp.Status]
	if !ok {
		return nil, errors.Errorf("unkown orderstatus %s", resp.Status)
	}

	var side exchange.OrderSide
	if resp.Direction == "buy" {
		if resp.Offset == "open" {
			side = exchange.OrderSideBuy
		} else {
			side = exchange.OrderSideCloseShort
		}
	} else {
		if resp.Offset == "open" {
			side = exchange.OrderSideSell
		} else {
			side = exchange.OrderSideCloseLong
		}
	}

	return &exchange.Order{
		Symbol:      symbol,
		ID:          exchange.NewIntID(resp.OrderID),
		ClientID:    exchange.NewIntID(resp.ClientOrderID),
		Amount:      decimal.NewFromFloat(resp.Volume),
		Price:       decimal.NewFromFloat(resp.Price),
		Filled:      decimal.NewFromFloat(float64(resp.TradeVolume)),
		AvgPrice:    decimal.NewFromFloat(resp.TradeAvgPrice),
		Fee:         decimal.NewFromFloat(resp.Fee),
		FeeCurrency: resp.FeeAsset,
		Created:     created,
		Status:      st,
		Side:        side,
		Raw:         resp,
	}, nil
}
