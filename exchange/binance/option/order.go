package option

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	GetOrderReq struct {
		*binance.RestReq
	}

	PostOrdreReq struct {
		*binance.RestReq
	}

	OrderResp struct {
		ID            string          `json:"id"`
		Symbol        string          `json:"symbol"`
		Price         decimal.Decimal `json:"price"`
		Quantity      decimal.Decimal `json:"quantity"`
		ExecutedQty   decimal.Decimal `json:"executedQty"`
		Fee           decimal.Decimal `json:"fee"`
		Side          string          `json:"side"`
		Type          string          `json:"type"`
		TimeInForce   string          `json:"timeInForce"`
		CreateDate    int64           `json:"createDate"`
		Status        string          `json:"status"`
		AvgPrice      decimal.Decimal `json:"avgPrice"`
		Source        string          `json:"source"`
		ReduceOnly    bool            `json:"reduceOnly"`
		ClientOrderID string          `json:"clientOrderId"`
	}

	OrderStatus int
)

const (
	OrderEndPoint   = "/vapi/v1/order"
	OrderSideBuy    = "BUY"
	OrderSideSell   = "SELL"
	OrderTypeLimit  = "LIMIT"
	OrderTypeMarket = "MARKET"
)

const (
	OrderStatusReceived OrderStatus = iota
	OrderStatusUntriggered
	OrderStatusAccepted
	OrderStatusRejected
	OrderStatusPartiallyFilled
	OrderStatusFilled
	OrderStatusCancelling
	OrderStatusCancelled
)

var (
	omStatusToStr = map[OrderStatus]string{
		OrderStatusReceived:        "RECEIVED",
		OrderStatusUntriggered:     "UNTRIGGERED",
		OrderStatusAccepted:        "ACCEPTED",
		OrderStatusRejected:        "REJECTED",
		OrderStatusPartiallyFilled: "PARTIALLY_FILLED",
		OrderStatusFilled:          "FILLED",
		OrderStatusCancelling:      "CANCELLING",
		OrderStatusCancelled:       "CANCELLED",
	}

	bnOrderStatusToStatus = map[OrderStatus]exchange.OrderStatus{
		OrderStatusReceived:        exchange.OrderStatusOpen,
		OrderStatusUntriggered:     exchange.OrderStatusOpen,
		OrderStatusAccepted:        exchange.OrderStatusOpen,
		OrderStatusRejected:        exchange.OrderStatusFailed,
		OrderStatusPartiallyFilled: exchange.OrderStatusOpen,
		OrderStatusFilled:          exchange.OrderStatusDone,
		OrderStatusCancelling:      exchange.OrderStatusFailed,
		OrderStatusCancelled:       exchange.OrderStatusFailed,
	}

	bnOrderStatusStrToStatus = map[string]exchange.OrderStatus{}

	bnOrderSideToSide = map[string]exchange.OrderSide{
		OrderSideBuy:  exchange.OrderSideBuy,
		OrderSideSell: exchange.OrderSideSell,
	}
	sideToBnOrderSide = map[exchange.OrderSide]string{}

	bnOrderTypeToType = map[string]exchange.OrderType{
		OrderTypeLimit:  exchange.OrderTypeLimit,
		OrderTypeMarket: exchange.OrderTypeMarket,
	}
	typeToBnOrderType = map[exchange.OrderType]string{}
)

func init() {
	for k, v := range bnOrderStatusToStatus {
		bnOrderStatusStrToStatus[k.String()] = v
	}

	for k, v := range bnOrderSideToSide {
		sideToBnOrderSide[v] = k
	}

	for k, v := range bnOrderTypeToType {
		typeToBnOrderType[v] = k
	}
}

func (os OrderStatus) String() string {
	ret, ok := omStatusToStr[os]
	if !ok {
		panic(fmt.Sprintf("unkown os=%d", os))
	}
	return ret
}

//NewPostOrderReq build create order request, the amount and price param will be formatted according to symbol precision
func NewPostOrderReq(symbol string, side string, typ string, amount float64, price float64) (*PostOrdreReq, error) {
	s, err := ParseSymbol(symbol)
	if err != nil {
		return nil, errors.WithMessage(err, "parse symbol fail")
	}

	at := decimal.NewFromFloat(amount)
	pr := decimal.NewFromFloat(price)

	req := binance.NewRestReq()
	req.AddFields("symbol", symbol)
	req.AddFields("side", side)
	req.AddFields("type", typ)
	req.AddFields("quantity", at.StringFixed(int32(s.AmountExponent())))
	req.AddFields("price", pr.StringFixed(int32(s.PriceExponent())))

	return &PostOrdreReq{
		RestReq: req,
	}, nil
}

func (or *PostOrdreReq) ReduceOnly(r bool) *PostOrdreReq {
	or.AddFields("reduceOnly", r)
	return or
}

func NewDeleteOrderReq(symbol string, orderID string) *GetOrderReq {
	req := NewGetOrderReq(symbol)
	req.OrderID(orderID)
	return req
}

func NewGetOrderReq(symbol string) *GetOrderReq {
	req := binance.NewRestReq()
	req.AddFields("symbol", symbol)

	return &GetOrderReq{
		RestReq: req,
	}
}

func (gor *GetOrderReq) OrderID(orderId string) *GetOrderReq {
	gor.AddFields("orderId", orderId)
	return gor
}

func (gor *GetOrderReq) ClientOrderID(cid string) *GetOrderReq {
	gor.AddFields("clientOrderId", cid)
	return gor
}

func (resp *OrderResp) Transfer() (*exchange.Order, error) {
	sym, err := ParseSymbol(resp.Symbol)
	if err != nil {
		return nil, errors.WithMessage(err, "parse symbol fail")
	}

	status, ok := bnOrderStatusStrToStatus[resp.Status]
	if !ok {
		return nil, errors.Errorf("unkown orderStatus='%s'", resp.Status)
	}

	typ, ok := bnOrderTypeToType[resp.Type]
	if !ok {
		return nil, errors.Errorf("unkown orderType='%s'", resp.Type)
	}

	side, ok := bnOrderSideToSide[resp.Side]
	if !ok {
		return nil, errors.Errorf("unkown orderSide='%s'", resp.Side)
	}

	ts := time.Unix(resp.CreateDate/1e3, (resp.CreateDate%1e3)*1e6)

	return &exchange.Order{
		ID:       exchange.NewStrID(resp.ID),
		ClientID: exchange.NewStrID(resp.ClientOrderID),
		Symbol:   sym,
		Status:   status,
		Side:     side,
		Type:     typ,
		Price:    resp.Price,
		AvgPrice: resp.AvgPrice,
		Amount:   resp.Quantity,
		Filled:   resp.ExecutedQty,
		Fee:      resp.Fee,
		Created:  ts,
		Raw:      resp,
	}, nil
}

func (rc *RestClient) PostOrder(ctx context.Context, req *PostOrdreReq) (*OrderResp, error) {
	var ret OrderResp

	resp := RestResp{
		Data: &ret,
	}

	values, err := req.Values()
	if err != nil {
		return nil, errors.WithMessage(err, "get values fail")
	}

	if err := rc.Request(ctx, http.MethodPost, OrderEndPoint, values, nil, true, &resp); err != nil {
		return nil, errors.WithMessage(err, "request fail")
	}

	if resp.Code != 0 {
		return nil, errors.Errorf("invalid resp code=%d msg=%s", resp.Code, resp.Msg)
	}

	return &ret, nil
}

func (rc *RestClient) GetOrder(ctx context.Context, req *GetOrderReq) (*OrderResp, error) {
	var ret OrderResp

	if err := rc.GetRequest(ctx, OrderEndPoint, req, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "get request fail")
	}

	return &ret, nil
}

func (rc *RestClient) DeleteOrder(ctx context.Context, req *GetOrderReq) (*OrderResp, error) {
	var ret OrderResp
	resp := RestResp{
		Data: &ret,
	}

	values, err := req.Values()
	if err != nil {
		return nil, errors.WithMessage(err, "get values fail")
	}

	if err := rc.Request(ctx, http.MethodDelete, OrderEndPoint, values, nil, true, &resp); err != nil {
		return nil, errors.WithMessage(err, "request fail")
	}

	if resp.Code != 0 {
		return nil, errors.Errorf("invalid resp code=%d msg=%s", resp.Code, resp.Msg)
	}

	return &ret, nil
}

func (rc *RestClient) CreateOrder(ctx context.Context, req *exchange.OrderRequest, options ...exchange.OrderReqOption) (*exchange.Order, error) {
	if len(options) != 0 {
		return nil, errors.Errorf("options is not support yet")
	}

	side, ok := sideToBnOrderSide[req.Side]
	if !ok {
		return nil, errors.Errorf("unknown side='%d'", req.Side)
	}
	typ, ok := typeToBnOrderType[req.Type]
	if !ok {
		return nil, errors.Errorf("unknown typ='%d'", req.Type)
	}
	price, _ := req.Price.Float64()
	amt, _ := req.Amount.Float64()

	or, err := NewPostOrderReq(req.Symbol.String(), side, typ, amt, price)
	if err != nil {
		return nil, errors.WithMessage(err, "create order req fail")
	}

	resp, err := rc.PostOrder(ctx, or)
	if err != nil {
		return nil, errors.WithMessage(err, "postOrder fail")
	}
	return resp.Transfer()
}

func (rc *RestClient) CancelOrder(ctx context.Context, order *exchange.Order) (*exchange.Order, error) {
	resp, err := rc.DeleteOrder(ctx, NewDeleteOrderReq(order.Symbol.String(), order.ID.String()))
	if err != nil {
		return nil, errors.WithMessage(err, "delete order fail")
	}

	return resp.Transfer()
}

func (rc *RestClient) FetchOrder(ctx context.Context, order *exchange.Order) (*exchange.Order, error) {
	resp, err := rc.GetOrder(ctx, NewGetOrderReq(order.Symbol.String()).OrderID(order.ID.String()))
	if err != nil {
		return nil, errors.WithMessage(err, "fetch order fail")
	}

	return resp.Transfer()
}

//ParseOrderResp extract orderResp info from order.Raw field
//the order param must be get via CreateOrder, FetchOrder, CancelOrder
func ParseOrderResp(order *exchange.Order) (*OrderResp, error) {
	if order.Raw == nil {
		return nil, errors.Errorf("no raw info")
	}

	resp, ok := order.Raw.(*OrderResp)
	if !ok {
		return nil, errors.Errorf("incorrect raw type")
	}

	return resp, nil
}
