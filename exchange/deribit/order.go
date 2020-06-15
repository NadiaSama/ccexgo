package deribit

import "github.com/NadiaSama/ccexgo/exchange"

type (
	OrderParam struct {
		InstrumentName string  `json:"instrument_name"`
		Amount         float64 `json:"amount"`
		Price          float64 `json:"price"`
		Type           string  `json:"limit"`
	}

	OrderResult struct {
		Price                float64 `json:"price"`
		Amount               float64 `json:"amount"`
		AveragePrice         float64 `json:"average_price"`
		OrderState           string  `json:"order_state"`
		OrderID              string  `json:"order_id"`
		LastUpdatedTimestamp int     `json:"last_updated_timestamp"`
		CreationTimestamp    int     `json:"creation_timestamp"`
		Commision            float64 `json:"commision"`
	}
)

var (
	type2Str map[exchange.OrderType]string = map[exchange.OrderType]string{
		exchange.OrderTypeLimit:      "limit",
		exchange.OrderTypeMarket:     "market",
		exchange.OrderTypeStopLimit:  "stop_limit",
		exchange.OrderTypeStopMarket: "stop_market",
	}
)

func (c *Client) optionCreateOrder(op exchange.OptionSymbol, side exchange.OrderSide,
	price, amount float64, typ exchange.OrderType, options interface{}) (*exchange.Order, error) {
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

	return nil, nil
}
