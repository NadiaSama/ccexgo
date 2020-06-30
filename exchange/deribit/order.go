package deribit

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/misc/tconv"
	"github.com/pkg/errors"
)

type (
	orderParam struct {
		AuthToken
		InstrumentName string  `json:"instrument_name"`
		Amount         float64 `json:"amount"`
		Price          float64 `json:"price"`
		Type           string  `json:"type"`
		PostOnly       bool    `json:"post_only,omitempty"`
		TimeInForce    string  `json:"time_in_force,omitempty"`
	}

	orderResult struct {
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
		InstrumentName       string  `json:"instrument_name"`
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

	tifMap map[exchange.TimeInForceFlag]string = map[exchange.TimeInForceFlag]string{
		exchange.TimeInForceFOK: "fill_or_kill",
		exchange.TimeInForceGTC: "good_til_cancelled",
		exchange.TimeInForceIOC: "immediate_or_cancel",
	}
)

func (c *Client) OptionCreateOrder(ctx context.Context, req *exchange.OrderRequest, opts ...exchange.OrderReqOption) (*exchange.Order, error) {
	var method string
	if req.Side == exchange.OrderSideBuy {
		method = "/private/buy"
	} else {
		method = "/private/sell"
	}

	param := &orderParam{
		Amount:         req.Amount,
		Price:          req.Price,
		InstrumentName: req.Symbol.String(),
		Type:           type2Str[req.Type],
	}

	for _, opt := range opts {
		switch msg := opt.(type) {
		case *exchange.PostOnlyOption:
			param.PostOnly = msg.PostOnly

		case *exchange.TimeInForceOption:
			val, ok := tifMap[msg.Flag]
			if !ok {
				return nil, exchange.NewBadArg("invalid TimeInForceOption", msg)
			}
			param.TimeInForce = val

		default:
			return nil, exchange.NewBadArg("unsupport option value", msg)
		}
	}

	var or orderResult
	if err := c.call(ctx, method, param, &or, true); err != nil {
		return nil, err
	}

	return or.Order.transform()
}

func (c *Client) OptionFetchOrder(ctx context.Context, order *exchange.Order) (*exchange.Order, error) {
	param := map[string]interface{}{
		"order_id": order.ID,
	}
	var r Order
	if err := c.call(ctx, "/private/get_order_state", param, &r, true); err != nil {
		return nil, err
	}

	return r.transform()
}

func (c *Client) OptionCancelOrder(ctx context.Context, order *exchange.Order) (*exchange.Order, error) {
	param := map[string]interface{}{
		"order_id": order.ID,
	}

	var r Order
	if err := c.call(ctx, "/private/cancel", param, &r, true); err != nil {
		return nil, err
	}
	return r.transform()
}

func (order *Order) transform() (*exchange.Order, error) {
	create := tconv.Milli2Time(order.CreationTimestamp)
	update := tconv.Milli2Time(order.LastUpdatedTimestamp)
	sym, err := ParseOptionSymbol(order.InstrumentName)
	if err != nil {
		return nil, errors.WithMessagef(err, "parse symbol %s fail", order.InstrumentName)
	}
	return &exchange.Order{
		ID:       order.OrderID,
		Amount:   order.Amount,
		Price:    order.Price,
		AvgPrice: order.AveragePrice,
		Status:   statusMap[order.OrderState],
		Side:     directionMap[order.Direction],
		Created:  create,
		Updated:  update,
		Symbol:   sym,
		Filled:   order.FilledAmont,
	}, nil
}
