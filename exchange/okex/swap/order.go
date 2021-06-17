package swap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const (
	orderTable    = "swap/order"
	orderEndPoint = "/api/swap/v3/order"

	orderTypeNormal = "0"
	orderTypeMaker  = "1"
	orderTypeFOK    = "2"
	orderTypeIOC    = "3"
	orderTypeMarket = "4"
)

type (
	OrderChannel struct {
		sym exchange.SwapSymbol
	}

	Order struct {
		InstrumentID string          `json:"instrument_id"`
		Size         decimal.Decimal `json:"size"`
		Timestamp    string          `json:"timestamp"`
		FilledQty    decimal.Decimal `json:"filled_qty"`
		Fee          decimal.Decimal `json:"fee"`
		OrderID      string          `json:"order_id"`
		ClientID     string          `json:"client_id"`
		Price        decimal.Decimal `json:"price"`
		PriceAvg     decimal.Decimal `json:"price_avg"`
		Type         string          `json:"type"`
		ContractVal  decimal.Decimal `json:"contract_val"`
		OrderType    string          `json:"order_type"`
		State        string          `json:"state"`
	}

	orderRequest struct {
		ClientOID    string `json:"client_oid"`
		Size         string `json:"size"`
		Type         string `json:"type"`
		OrderType    string `json:"order_type"`
		MatchPrice   string `json:"match_price"`
		Price        string `json:"price"`
		InstrumentID string `json:"instrument_id"`
	}

	orderResponse struct {
		OrderID      string `json:"order_id"`
		ClientOID    string `json:"client_oid"`
		ErrorCode    string `json:"error_code"`
		ErrorMessage string `json:"error_message"`
		Result       string `json:"result"`
	}
)

var (
	sideMap map[string]exchange.OrderSide = map[string]exchange.OrderSide{
		"1": exchange.OrderSideBuy,
		"2": exchange.OrderSideSell,
		"3": exchange.OrderSideCloseLong,
		"4": exchange.OrderSideCloseShort,
	}

	rSideMap map[exchange.OrderSide]string = map[exchange.OrderSide]string{
		exchange.OrderSideBuy:        "1",
		exchange.OrderSideSell:       "2",
		exchange.OrderSideCloseLong:  "3",
		exchange.OrderSideCloseShort: "4",
	}

	statusMap map[string]exchange.OrderStatus = map[string]exchange.OrderStatus{
		"-1": exchange.OrderStatusCancel,
		"0":  exchange.OrderStatusOpen,
		"1":  exchange.OrderStatusOpen,
		"2":  exchange.OrderStatusDone,
		"3":  exchange.OrderStatusOpen,
		"4":  exchange.OrderStatusOpen,
	}
)

func NewOrderChannel(symbol exchange.SwapSymbol) exchange.Channel {
	return &OrderChannel{
		sym: symbol,
	}
}

func (oc *OrderChannel) String() string {
	return fmt.Sprintf("%s:%s", orderTable, oc.sym.String())
}

func init() {
	okex.SubscribeCB(orderTable, parseOrderCB)
}

func parseOrderCB(table string, action string, raw json.RawMessage) (*rpc.Notify, error) {
	var orders []Order
	if err := json.Unmarshal(raw, &orders); err != nil {
		return nil, err
	}

	var os []*exchange.Order
	for _, o := range orders {
		order, err := o.Transform()
		if err != nil {
			return nil, err
		}

		os = append(os, order)
	}
	return &rpc.Notify{
		Method: orderTable,
		Params: os,
	}, nil
}

func (o *Order) Transform() (*exchange.Order, error) {
	sym, err := ParseSymbol(o.InstrumentID)
	if err != nil {
		return nil, err
	}

	side, ok := sideMap[o.Type]
	if !ok {
		return nil, errors.Errorf("unkown order side '%s'", o.Type)
	}

	status, ok := statusMap[o.State]
	if !ok {
		return nil, errors.Errorf("unkown order state '%s'", o.State)
	}

	ts, err := okex.ParseTime(o.Timestamp)
	if err != nil {
		return nil, err
	}

	return &exchange.Order{
		Symbol:   sym,
		Side:     side,
		ID:       exchange.NewStrID(o.OrderID),
		ClientID: exchange.NewStrID(o.ClientID),
		Status:   status,
		Amount:   o.Size,
		Filled:   o.FilledQty,
		Price:    o.Price,
		AvgPrice: o.PriceAvg,
		Fee:      o.Fee,
		Updated:  ts,
		Raw:      o,
	}, nil
}

func (rc *RestClient) CreateOrder(ctx context.Context, req *exchange.OrderRequest, options ...exchange.OrderReqOption) (*exchange.Order, error) {
	oReq := orderRequest{
		Size:         req.Amount.String(),
		InstrumentID: req.Symbol.String(),
		MatchPrice:   "0",
		Type:         rSideMap[req.Side],
		OrderType:    orderTypeNormal,
	}

	if req.ClientID != nil {
		oReq.ClientOID = req.ClientID.String()
	}

	if req.Type == exchange.OrderTypeMarket {
		oReq.OrderType = orderTypeMarket
	} else {
		oReq.Price = req.Price.String()
	}

	if len(options) > 1 {
		return nil, errors.Errorf("okex create order only one option is support")
	}

	if oReq.OrderType != orderTypeNormal && len(options) != 0 {
		return nil, errors.Errorf("okex cannot creat order with type=%s and options", oReq.OrderType)
	}

	for _, opt := range options {
		switch t := opt.(type) {
		case *exchange.PostOnlyOption:
			oReq.OrderType = orderTypeMaker

		case *exchange.TimeInForceOption:
			if t.Flag == exchange.TimeInForceFOK {
				oReq.OrderType = orderTypeFOK
			} else if t.Flag == exchange.TimeInForceIOC {
				oReq.OrderType = orderTypeIOC
			}
		}
	}

	b, _ := json.Marshal(&oReq)
	body := bytes.NewBuffer(b)
	resp := orderResponse{}
	if err := rc.Request(ctx, http.MethodPost, orderEndPoint, nil, body, true, &resp); err != nil {
		return nil, err
	}

	if err := resp.Error(); err != nil {
		return nil, err
	}

	return &exchange.Order{
		ID:     exchange.NewStrID(resp.OrderID),
		Symbol: req.Symbol,
	}, nil
}

func (rc *RestClient) CancelOrder(ctx context.Context, order *exchange.Order) error {
	endPoint := fmt.Sprintf("/api/swap/v3/cancel_order/%s/%s", order.Symbol.String(), order.ID.String())
	var resp orderResponse
	if err := rc.Request(ctx, http.MethodPost, endPoint, nil, bytes.NewBuffer([]byte{}), true, &resp); err != nil {
		return err
	}

	return resp.Error()
}

func (or *orderResponse) Error() error {
	if or.ErrorCode != "0" {
		return errors.Errorf("okex order response error code=%s msg=\"%s\"",
			or.ErrorCode, or.ErrorMessage)
	}
	return nil
}
