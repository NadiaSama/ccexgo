package deribit

import (
	"context"

	"github.com/NadiaSama/ccexgo/misc"

	"github.com/NadiaSama/ccexgo/exchange"
)

type (
	OrderParam struct {
		AuthToken
		InstrumentName string  `json:"instrument_name"`
		Amount         float64 `json:"amount"`
		Price          float64 `json:"price"`
		Type           string  `json:"limit"`
	}

	OrderResult struct {
		Order Order `json:"order"`
	}

	Order struct {
		Price                float64 `json:"price"`
		Amount               float64 `json:"amount"`
		AveragePrice         float64 `json:"average_price"`
		OrderState           string  `json:"order_state"`
		OrderID              string  `json:"order_id"`
		LastUpdatedTimestamp int64   `json:"last_update_timestamp"`
		CreationTimestamp    int64   `json:"creation_timestamp"`
		Commision            float64 `json:"commision"`
		Direction            string  `json:"direction"`
		FilledAmont          float64 `json:"filled_amount"`
	}
)

var (
	type2Str map[exchange.OrderType]string = map[exchange.OrderType]string{
		exchange.OrderTypeLimit:      "limit",
		exchange.OrderTypeMarket:     "market",
		exchange.OrderTypeStopLimit:  "stop_limit",
		exchange.OrderTypeStopMarket: "stop_market",
	}

	statusMap map[string]exchange.OrderStatus = map[string]exchange.OrderStatus{
		"open":        exchange.OrderStatusOpen,
		"rejected":    exchange.OrderStatusCancel,
		"cancelled":   exchange.OrderStatusCancel,
		"filled":      exchange.OrderStatusDone,
		"untriggered": exchange.OrderStatusOpen,
	}

	directionMap map[string]exchange.OrderSide = map[string]exchange.OrderSide{
		"buy":  exchange.OrderSideBuy,
		"sell": exchange.OrderSideSell,
	}
)

func (c *Client) OptionCreateOrder(ctx context.Context, op exchange.OptionSymbol, side exchange.OrderSide,
	price, amount float64, typ exchange.OrderType, options ...interface{}) (*exchange.Order, error) {
	var method string
	if side == exchange.OrderSideBuy {
		method = "/private/buy"
	} else {
		method = "/private/sell"
	}

	param := &OrderParam{
		Amount:         amount,
		Price:          price,
		InstrumentName: op.String(),
		Type:           type2Str[typ],
	}
	var or OrderResult
	if err := c.call(ctx, method, param, &or, true); err != nil {
		return nil, err
	}

	return or.Order.transform(), nil
}

func (c *Client) OptionFetchOrder(ctx context.Context, order *exchange.Order) (*exchange.Order, error) {
	param := map[string]interface{}{
		"order_id": order.ID,
	}
	var r Order
	if err := c.call(ctx, "/private/get_order_state", param, &r, true); err != nil {
		return nil, err
	}

	return r.transform(), nil
}

func (c *Client) OptionCancelOrder(ctx context.Context, order *exchange.Order) (*exchange.Order, error) {
	param := map[string]interface{}{
		"order_id": order.ID,
	}

	var r Order
	if err := c.call(ctx, "/private/cancel", param, &r, true); err != nil {
		return nil, err
	}
	return r.transform(), nil
}

func (order *Order) transform() *exchange.Order {
	create := misc.Milli2Time(order.CreationTimestamp)
	update := misc.Milli2Time(order.LastUpdatedTimestamp)
	return &exchange.Order{
		ID:       order.OrderID,
		Amount:   order.Amount,
		Price:    order.Price,
		AvgPrice: order.AveragePrice,
		Status:   statusMap[order.OrderState],
		Side:     directionMap[order.Direction],
		Created:  create,
		Updated:  update,
		Filled:   order.FilledAmont,
	}
}
