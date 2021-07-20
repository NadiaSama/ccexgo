package swap

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	ordersChannel struct {
		symbol string
	}

	OrderNotifyTrade struct {
		TradeFee      float64 `json:"trade_fee"`
		FeeAsset      string  `json:"fee_asset"`
		TradeID       int64   `json:"trade_id"`
		ID            string  `json:"id"`
		TradeVolume   int     `json:"trade_volume"`
		TradePrice    float64 `json:"trade_price"`
		TradeTurnOver float64 `json:"trade_turnover"`
		CreatedAt     int64   `json:"created_at"`
		Profit        float64 `json:"profit"`
		RealProfit    float64 `json:"real_profit"`
	}

	OrderNotify struct {
		Op             string             `json:"op"`
		Topic          string             `json:"topic"`
		TS             int64              `json:"ts"` //ping message and auth message ts field have different type
		Symbol         string             `json:"symbol"`
		ContractCode   string             `json:"contract_code"`
		Volume         float64            `json:"volume"`
		Price          float64            `json:"price"`
		OrderPriceType string             `json:"order_price_type"`
		Direction      string             `json:"direction"`
		Offset         string             `json:"offset"`
		Status         int                `json:"status"`
		LeverRate      int                `json:"lever_rate"`
		OrderID        int64              `json:"order_id"`
		OrderIDStr     string             `json:"order_id_str"`
		ClientOrderID  int64              `json:"client_order_id"`
		OrderType      int                `json:"order_type"`
		CreatedAt      int64              `json:"created_at"`
		TradeVolume    int                `json:"trade_volume"`
		TradeTurnOver  float64            `json:"trade_turnover"`
		Fee            float64            `json:"fee"`
		FeeAsset       string             `json:"fee_asset"`
		TradeAvgPrice  float64            `json:"trade_avg_price"`
		CanceledAt     int64              `json:"canceled_at"`
		RealProfit     float64            `json:"real_profit"`
		Trades         []OrderNotifyTrade `json:"trades"`
	}
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

func NewOrdersChannel(symbol string) exchange.Channel {
	return &ordersChannel{
		symbol: symbol,
	}
}

func (oc *ordersChannel) String() string {
	return fmt.Sprintf("orders.%s", oc.symbol)
}

func ParseOrder(raw []byte) (*exchange.Order, error) {
	var resp OrderNotify
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	symbol, err := ParseSymbol(resp.ContractCode)
	if err != nil {
		return nil, err
	}

	var created time.Time
	if resp.CreatedAt != 0 {
		created = huobi.ParseTS(resp.CreatedAt)
	}

	var updated time.Time
	if resp.CanceledAt != 0 {
		updated = huobi.ParseTS(resp.CanceledAt)
	} else if len(resp.Trades) != 0 {
		updated = huobi.ParseTS(resp.Trades[0].CreatedAt)
	}

	st, ok := statusMap[resp.Status]
	if !ok {
		return nil, errors.Errorf("unkown orderstatus %s", resp.Status)
	}

	var side exchange.OrderSide
	if resp.Direction == OrderDirectionBuy {
		if resp.Offset == OrderOffsetOpen {
			side = exchange.OrderSideBuy
		} else if resp.Offset == OrderOffsetClose {
			side = exchange.OrderSideCloseShort
		} else {
			return nil, errors.Errorf("unkown order offset '%s'", resp.Offset)
		}
	} else if resp.Direction == OrderDirectionSell {
		if resp.Offset == OrderOffsetOpen {
			side = exchange.OrderSideSell
		} else if resp.Offset == OrderOffsetClose {
			side = exchange.OrderSideCloseLong
		} else {
			return nil, errors.Errorf("unkown order offset '%s'", resp.Offset)
		}
	} else {
		return nil, errors.Errorf("unkown order direction '%s'", resp.Direction)
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
		Updated:     updated,
		Status:      st,
		Side:        side,
		Raw:         &resp,
	}, nil
}
