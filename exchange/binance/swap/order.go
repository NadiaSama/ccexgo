package swap

import (
	"context"
	"net/http"
	"strconv"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/NadiaSama/ccexgo/misc/tconv"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	AddOrderReq struct {
		*binance.RestReq
	}

	OrderResp struct {
		binance.APIError                 //in case of error
		ClientOrderID    string          `json:"newOrderId"`
		CumQty           decimal.Decimal `json:"cumQty"`
		CumQuote         decimal.Decimal `json:"cumQuote"`
		ExecutedQty      decimal.Decimal `json:"executedQty"`
		OrderID          int64           `json:"orderId"`
		AvgPrice         decimal.Decimal `json:"avgPrice"`
		OrigQty          decimal.Decimal `json:"origQty"`
		Price            decimal.Decimal `json:"price"`
		ReduceOnly       bool            `json:"reduceOnly"`
		Side             string          `json:"side"`
		PositionSide     string          `json:"positionSide"`
		Status           string          `json:"status"`
		StopPrice        decimal.Decimal `json:"stopPrice"`
		ClosePosition    bool            `json:"closePosition"`
		Symbol           string          `json:"symbol"`
		TimeInForce      string          `json:"timeInForce"`
		Type             string          `json:"type"`
		OrigType         string          `json:"origType"`
		ActivatePrice    decimal.Decimal `json:"activatePrice"`
		PriceRate        decimal.Decimal `json:"priceRate"`
		UpdateTime       int64           `json:"updateTime"`
		WorkingType      string          `json:"workingType"`
		PriceProtect     bool            `json:"priceProtect"`
	}

	OrderReq struct {
		*binance.RestReq
	}
)

const (
	OrderEndPoint     = "/fapi/v1/order"
	PositionSideBoth  = "both"
	PositionSideLong  = "LONG"
	PositionSideShort = "SHORT"
	SideBuy           = "BUY"
	SideSell          = "SELL"
	OrderTypeMarket   = "MARKET"
	OrderTypeLimit    = "LIMIT"
	TimeInForce       = "GTC"
)

var (
	OrderType2ExType = map[string]exchange.OrderType{
		OrderTypeLimit:  exchange.OrderTypeLimit,
		OrderTypeMarket: exchange.OrderTypeMarket,
	}

	ExType2OrderType = map[exchange.OrderType]string{
		exchange.OrderTypeLimit:  OrderTypeLimit,
		exchange.OrderTypeMarket: OrderTypeMarket,
	}

	Status2ExStatus = map[string]exchange.OrderStatus{
		"NEW":               exchange.OrderStatusOpen,
		"PARTIALLY_FILLED ": exchange.OrderStatusOpen,
		"FILLED":            exchange.OrderStatusDone,
		"CANCELED":          exchange.OrderStatusCancel,
		"REJECTED":          exchange.OrderStatusFailed,
		"EXPIRED":           exchange.OrderStatusFailed,
	}
)

//NewAddOrderReq according symbol, side, type
func NewAddOrderReq(symbol string, side string, typ string) *AddOrderReq {
	req := binance.NewRestReq()
	req.AddFields("symbol", symbol)
	req.AddFields("side", side)
	req.AddFields("type", typ)

	return &AddOrderReq{
		RestReq: req,
	}
}

func (req *AddOrderReq) TimeInForce(tif string) *AddOrderReq {
	req.AddFields("timeInForce", tif)
	return req
}

func (req *AddOrderReq) PositionSide(side string) *AddOrderReq {
	req.AddFields("positionSide", side)
	return req
}

func (req *AddOrderReq) Price(prc decimal.Decimal) *AddOrderReq {
	req.AddFields("price", prc.String())
	return req
}

func (req *AddOrderReq) Quantity(q decimal.Decimal) *AddOrderReq {
	req.AddFields("quantity", q.String())
	return req
}

func (cl *RestClient) AddOrder(ctx context.Context, req *AddOrderReq) (*OrderResp, error) {
	values, err := req.Values()
	if err != nil {
		return nil, errors.WithMessage(err, "get param fail")
	}

	var ret OrderResp
	if err := cl.Request(ctx, http.MethodPost, OrderEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "add order fail")
	}

	return &ret, nil
}

func NewOrderReq(symbol string) *OrderReq {
	req := binance.NewRestReq()
	req.AddFields("symbol", symbol)
	return &OrderReq{
		RestReq: req,
	}
}

func (r *OrderReq) OrderID(id int64) *OrderReq {
	r.AddFields("orderId", id)
	return r
}

func (cl *RestClient) GetOrder(ctx context.Context, req *OrderReq) (*OrderResp, error) {
	var ret OrderResp
	if err := cl.GetRequest(ctx, OrderEndPoint, req, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "get order req fail")
	}

	return &ret, nil
}

func (cl *RestClient) DeleteOrder(ctx context.Context, req *OrderReq) (*OrderResp, error) {
	var ret OrderResp
	values, err := req.Values()
	if err != nil {
		return nil, errors.WithMessage(err, "get req fail")
	}

	if err := cl.Request(ctx, http.MethodDelete, OrderEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "cancel order fail")
	}

	return &ret, nil
}

func (resp *OrderResp) Transfer() (*exchange.Order, error) {
	symbol, err := ParseSymbol(resp.Symbol)
	if err != nil {
		return nil, errors.WithMessage(err, "parse symbol fail")
	}

	var (
		side   exchange.OrderSide
		typ    exchange.OrderType
		status exchange.OrderStatus
	)

	typ, ok := OrderType2ExType[resp.Type]
	if !ok {
		return nil, errors.Errorf("unknown resp type=%+v", resp.Type)
	}

	status, ok = Status2ExStatus[resp.Status]
	if !ok {
		return nil, errors.Errorf("unknown resp status=%s", resp.Status)
	}

	if resp.PositionSide == PositionSideBoth {
		if resp.Side == SideBuy {
			side = exchange.OrderSideBuy
		} else if resp.Side == SideSell {
			side = exchange.OrderSideSell
		} else {
			return nil, errors.Errorf("unknown side=%s", resp.Side)
		}
	} else if resp.PositionSide == PositionSideLong {
		if resp.Side == SideBuy {
			side = exchange.OrderSideBuy
		} else if resp.Side == SideSell {
			side = exchange.OrderSideCloseLong
		} else {
			return nil, errors.Errorf("unknown side=%s", resp.Side)
		}
	} else if resp.PositionSide == PositionSideShort {
		if resp.Side == SideBuy {
			side = exchange.OrderSideCloseShort
		} else if resp.Side == SideSell {
			side = exchange.OrderSideSell
		} else {
			return nil, errors.Errorf("unknown side=%s", resp.Side)
		}
	} else {
		return nil, errors.Errorf("unknown positionSide=%s", resp.PositionSide)
	}

	return &exchange.Order{
		ID:       exchange.NewIntID(resp.OrderID),
		Symbol:   symbol,
		Amount:   resp.OrigQty,
		Price:    resp.Price,
		Type:     typ,
		Side:     side,
		AvgPrice: resp.AvgPrice,
		Status:   status,
		Updated:  tconv.Milli2Time(resp.UpdateTime),
		Filled:   resp.ExecutedQty,
		Raw:      resp,
	}, nil
}

func (cl *RestClient) CreateOrder(ctx context.Context, req *exchange.OrderRequest) (*exchange.Order, error) {
	if cl.side == nil {
		return nil, errors.Errorf("positionSide not init")
	}

	typ, ok := ExType2OrderType[req.Type]
	if !ok {
		return nil, errors.Errorf("unknown type=%s", req.Type)
	}

	var (
		side         string
		positionSide string
	)

	//dual side position
	if cl.side.DualSidePosition {
		switch req.Side {
		case exchange.OrderSideBuy:
			side = SideBuy
			positionSide = PositionSideLong

		case exchange.OrderSideSell:
			side = SideSell
			positionSide = PositionSideShort

		case exchange.OrderSideCloseLong:
			side = SideSell
			positionSide = PositionSideLong

		case exchange.OrderSideCloseShort:
			side = SideBuy
			positionSide = PositionSideShort

		default:
			return nil, errors.Errorf("unknown side=%s", req.Side)
		}
	} else {
		positionSide = "BOTH"
		switch req.Side {
		case exchange.OrderSideBuy:
			fallthrough
		case exchange.OrderSideCloseShort:
			side = SideBuy

		case exchange.OrderSideSell:
			fallthrough
		case exchange.OrderSideCloseLong:
			side = SideSell

		default:
			return nil, errors.Errorf("unknown side=%s", req.Side)
		}
	}

	or := NewAddOrderReq(req.Symbol.String(), side, typ)
	if typ == OrderTypeLimit {
		or.Price(req.Price)
		or.TimeInForce(TimeInForce)
	}
	or.Quantity(req.Amount)
	or.PositionSide(positionSide)

	resp, err := cl.AddOrder(ctx, or)
	if err != nil {
		return nil, errors.WithMessage(err, "add order fail")
	}

	return resp.Transfer()
}

func (cl *RestClient) FetchOrder(ctx context.Context, order *exchange.Order) (*exchange.Order, error) {
	id, err := strconv.ParseInt(order.ID.String(), 10, 64)
	if err != nil {
		return nil, errors.WithMessagef(err, "bad orderID=%s", order.ID.String())
	}
	req := NewOrderReq(order.Symbol.String()).OrderID(id)

	resp, err := cl.GetOrder(ctx, req)
	if err != nil {
		return nil, errors.WithMessagef(err, "get order fail ID=%s", order.ID.String())
	}

	return resp.Transfer()
}

func (cl *RestClient) CancelOrder(ctx context.Context, order *exchange.Order) (*exchange.Order, error) {
	id, err := strconv.ParseInt(order.ID.String(), 10, 64)
	if err != nil {
		return nil, errors.WithMessagef(err, "bad orderID=%s", order.ID.String())
	}
	req := NewOrderReq(order.Symbol.String()).OrderID(id)

	resp, err := cl.DeleteOrder(ctx, req)
	if err != nil {
		return nil, errors.WithMessagef(err, "get order fail ID=%s", order.ID.String())
	}

	return resp.Transfer()

}
