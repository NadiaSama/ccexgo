package spot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	OrderParam struct {
		Type         string `json:"type"`
		Side         string `json:"side"`
		InstrumentID string `json:"instrument_id"`
		Size         string `json:"size"`
		ClinetOID    string `json:"client_oid"`
		Price        string `json:"price,omitempty"`
		OrderType    string `json:"order_type"`
		Notional     string `json:"notional,omitempty"`
	}

	OrderResponse struct {
		OrderID      string `json:"order_id"`
		ClientID     string `json:"client_id"`
		Result       bool   `json:"result"`
		ErrorCode    string `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	}

	FetchOrderResponse struct {
		OrderID        string `json:"order_id"`
		ClientOID      string `json:"client_oid"`
		Price          string `json:"price"`
		Size           string `json:"size"`
		OrderType      string `json:"order_type"`
		Notional       string `json:"notional"`
		InstrumentID   string `json:"instrument_id"`
		Side           string `json:"side"`
		Type           string `json:"type"`
		Timestamp      string `json:"timestamp"`
		FilledSize     string `json:"filled_size"`
		FilledNotional string `json:"filled_notional"`
		State          string `json:"state"`
		PriceAvg       string `json:"price_avg"`
		FeeCurrency    string `json:"fee_currency"`
		Fee            string `json:"fee"`
		RebateCurrency string `json:"rebate_currency"`
		Rebate         string `json:"rebate"`
	}
)

const (
	OrderSideBuy      = "buy"
	OrderSideSell     = "sell"
	OrderTypeMarket   = "market"
	OrderTypeLimit    = "limit"
	OrderTypeNoraml   = "0"
	OrderTypePostOnly = "1"
	OrderTypeFOK      = "2"
	OrderTypeIOC      = "3"

	CreateOrderEndPoint = "/api/spot/v3/orders"
)

var (
	stateMap = map[string]exchange.OrderStatus{
		"-2": exchange.OrderStatusFailed,
		"-1": exchange.OrderStatusCancel,
		"0":  exchange.OrderStatusOpen,
		"1":  exchange.OrderStatusOpen,
		"2":  exchange.OrderStatusDone,
		"3":  exchange.OrderStatusOpen,
		"4":  exchange.OrderStatusOpen,
	}
)

//CreateOrder create a spot order
func (rc *RestClient) CreateOrder(ctx context.Context, req *exchange.OrderRequest, options ...exchange.OrderReqOption) (*exchange.Order, error) {
	op := OrderParam{
		InstrumentID: req.Symbol.String(),
		Size:         req.Amount.String(),
	}
	if req.ClientID != nil {
		op.ClinetOID = req.ClientID.String()
	}

	if req.Type == exchange.OrderTypeLimit {
		op.Type = OrderTypeLimit
		op.Price = req.Price.String()
	} else if req.Type == exchange.OrderTypeMarket {
		if len(options) != 0 {
			return nil, errors.Errorf("okex market order do not support options")
		}
		op.Type = OrderTypeMarket
		op.Notional = req.Price.Mul(req.Amount).String()
	} else {
		return nil, errors.Errorf("unsupport orderType %s", req.Type)
	}

	if req.Side == exchange.OrderSideBuy {
		op.Side = OrderSideBuy
	} else if req.Side == exchange.OrderSideSell {
		op.Side = OrderSideSell
	} else {
		return nil, errors.Errorf("unsupport orderSide %s", req.Side)
	}

	if len(options) > 1 {
		return nil, errors.Errorf("multiple option do not support")
	}

	op.OrderType = OrderTypeNoraml

	for _, option := range options {
		switch t := option.(type) {
		case *exchange.PostOnlyOption:
			if t.PostOnly {
				op.OrderType = OrderTypeMarket
			}

		case *exchange.TimeInForceOption:
			if t.Flag == exchange.TimeInForceFOK {
				op.OrderType = OrderTypeFOK
			} else if t.Flag == exchange.TimeInForceIOC {
				op.OrderType = OrderTypeIOC
			} else {
				return nil, errors.Errorf("unsuport timeinfor option %s", t.Flag)
			}
		}
	}

	var resp OrderResponse
	b, _ := json.Marshal(op)
	buf := bytes.NewBuffer(b)
	if err := rc.Request(ctx, http.MethodPost, CreateOrderEndPoint, nil, buf, true, &resp); err != nil {
		return nil, err
	}
	return resp.Transform(req.Symbol)
}

func (rc *RestClient) FetchOrder(ctx context.Context, order *exchange.Order) (*exchange.Order, error) {
	u := fmt.Sprintf("/api/spot/v3/orders/%s", order.ID.String())
	params := url.Values{}
	params.Add("instrument_id", order.Symbol.String())

	var resp FetchOrderResponse
	if err := rc.Request(ctx, http.MethodGet, u, params, nil, true, &resp); err != nil {
		return nil, err
	}

	return resp.Transform()
}

func (rc *RestClient) CancelOrder(ctx context.Context, order *exchange.Order) error {
	u := fmt.Sprintf("/api/spot/v3/cancel_orders/%s", order.ID.String())
	params := url.Values{}
	params.Add("instrument_id", order.Symbol.String())

	var resp OrderResponse
	if err := rc.Request(ctx, http.MethodPost, u, params, nil, true, &resp); err != nil {
		return err
	}

	if !resp.Result {
		return errors.Errorf("cancel order error error_code=%s error_message='%s'", resp.ErrorCode, resp.ErrorMessage)
	}
	return nil
}

func (resp *OrderResponse) Transform(sym exchange.Symbol) (*exchange.Order, error) {
	if !resp.Result {
		return nil, errors.Errorf("create order fail error_code=%s error_message='%s'", resp.ErrorCode, resp.ErrorMessage)
	}
	return &exchange.Order{
		Symbol:   sym,
		ID:       exchange.NewStrID(resp.OrderID),
		ClientID: exchange.NewStrID(resp.ClientID),
		Status:   exchange.OrderStatusOpen,
		Raw:      resp,
	}, nil
}

func (fResp *FetchOrderResponse) Transform() (*exchange.Order, error) {
	var (
		side exchange.OrderSide
	)

	if fResp.Side == OrderSideBuy {
		side = exchange.OrderSideBuy
	} else if fResp.Side == OrderSideSell {
		side = exchange.OrderSideSell
	} else {
		return nil, errors.Errorf("unkown order side '%s'", fResp.Side)
	}

	fee, err := toDecimal(fResp.Fee)
	if err != nil {
		return nil, err
	}
	price, err := toDecimal(fResp.Price)
	if err != nil {
		return nil, err
	}
	priceAvg, err := toDecimal(fResp.PriceAvg)
	if err != nil {
		return nil, err
	}

	var (
		typ    exchange.OrderType
		filled decimal.Decimal
		amount decimal.Decimal
	)
	if fResp.Type == OrderTypeLimit {
		typ = exchange.OrderTypeLimit
		filled, err = toDecimal(fResp.FilledSize)
		if err != nil {
			return nil, err
		}
		amount, err = toDecimal(fResp.Size)
		if err != nil {
			return nil, err
		}
	} else if fResp.Type == OrderTypeMarket {
		typ = exchange.OrderTypeMarket
		fn, err := toDecimal(fResp.FilledNotional)
		if err != nil {
			return nil, err
		}
		n, err := toDecimal(fResp.Notional)
		filled = fn.Div(priceAvg)
		amount = n.Div(priceAvg)
	} else {
		return nil, errors.Errorf("unkown order type '%s'", fResp.Type)
	}

	st, err := okex.ParseTime(fResp.Timestamp)
	if err != nil {
		return nil, err
	}

	sym, err := ParseSymbol(fResp.InstrumentID)
	if err != nil {
		return nil, err
	}

	status, ok := stateMap[fResp.State]
	if !ok {
		return nil, errors.Errorf("unkown order state '%s'", fResp.State)
	}
	ret := &exchange.Order{
		Symbol:      sym,
		ID:          exchange.NewStrID(fResp.OrderID),
		ClientID:    exchange.NewStrID(fResp.ClientOID),
		Price:       price,
		AvgPrice:    priceAvg,
		Filled:      filled,
		Amount:      amount,
		Status:      status,
		Type:        typ,
		Created:     st,
		Updated:     time.Now(),
		Side:        side,
		Fee:         fee.Abs(),
		FeeCurrency: fResp.FeeCurrency,
		Raw:         fResp,
	}

	return ret, nil
}

func toDecimal(str string) (decimal.Decimal, error) {
	if str == "" {
		return decimal.Zero, nil
	}

	ret, err := decimal.NewFromString(str)
	return ret, err
}
