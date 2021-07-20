package spot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	PlaceReq struct {
		data map[string]string
	}

	OrdersReq struct {
		OrderID string
	}

	SubmitCancelReq OrdersReq
	MatchResultReq  OrdersReq

	PlaceResp struct {
		Status  string `json:"status"`
		Data    string `json:"data"`
		ErrCode string `json:"err-code"`
		ErrMsg  string `json:"err-msg"`
	}

	OrdersRespDetail struct {
		ID               int64  `json:"id"`
		Symbol           string `json:"symbol"`
		AccountID        int64  `json:"account-id"`
		Amount           string `json:"amount"`
		Price            string `json:"price"`
		CreatedAt        int64  `json:"created-at"`
		Type             string `json:"type"`
		FilledAmount     string `json:"filled-amount"`
		FilledCashAmount string `json:"filled-cash-amount"`
		FilledFees       string `json:"filled-fees"`
		FinishedAt       int64  `json:"finished-at"`
		UserID           int64  `json:"user-id"`
		Source           string `json:"source"`
		State            string `json:"state"`
		CanceledAt       int64  `json:"canceled-at"`
	}

	OrdersResp struct {
		Data OrdersRespDetail `json:"data"`
	}
)

const (
	PlaceOrderEndPoint = "/v1/order/orders/place"
)

func NewOrdersReq(orderID string) *OrdersReq {
	return &OrdersReq{
		OrderID: orderID,
	}
}

func NewSubmitCancelReq(orderID string) *SubmitCancelReq {
	return &SubmitCancelReq{
		OrderID: orderID,
	}
}

func NewMatchResultReq(orderID string) *MatchResultReq {
	return &MatchResultReq{
		OrderID: orderID,
	}
}

func NewPlaceReq(accountID string, symbol string, typ string, amount string) *PlaceReq {
	ret := &PlaceReq{
		data: make(map[string]string),
	}
	ret.data["account-id"] = accountID
	ret.data["symbol"] = symbol
	ret.data["type"] = typ
	ret.data["amount"] = amount

	return ret
}

func (pr *PlaceReq) Price(price string) *PlaceReq {
	pr.data["price"] = price
	return pr
}

func (pr *PlaceReq) Source(source string) *PlaceReq {
	pr.data["source"] = source
	return pr
}

func (pr *PlaceReq) ClientOrderID(coi string) *PlaceReq {
	pr.data["client-order-id"] = coi
	return pr
}

func (pr *PlaceReq) StopPrice(sp string) *PlaceReq {
	pr.data["stop-price"] = sp
	return pr
}

func (pr *PlaceReq) Operator(op string) *PlaceReq {
	pr.data["operator"] = op
	return pr
}

func (pr *PlaceReq) Serialize() ([]byte, error) {
	b, e := json.Marshal(pr.data)
	if e != nil {
		return nil, errors.WithMessage(e, "marshal json fail")
	}
	return b, nil
}

func (rc *RestClient) Place(ctx context.Context, req *PlaceReq) (*PlaceResp, error) {
	msg, err := req.Serialize()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(msg)
	var resp PlaceResp
	if err := rc.RequestWithRawResp(ctx, http.MethodPost, PlaceOrderEndPoint, nil, buf, true, &resp); err != nil {
		return nil, err
	}

	if resp.Status != "ok" {
		return nil, errors.Errorf("place order error %+v", resp)
	}

	return &resp, nil
}

func (rc *RestClient) Orders(ctx context.Context, req *OrdersReq) (*OrdersResp, error) {
	url := fmt.Sprintf("/v1/order/orders/%s", req.OrderID)

	var resp OrdersResp
	if err := rc.RequestWithRawResp(ctx, http.MethodGet, url, nil, nil, true, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (rc *RestClient) SubmitCancel(ctx context.Context, req *SubmitCancelReq) (*PlaceResp, error) {
	url := fmt.Sprintf("/v1/order/orders/%s/submitcancel", req.OrderID)

	var resp PlaceResp
	if err := rc.RequestWithRawResp(ctx, http.MethodPost, url, nil, nil, true, &resp); err != nil {
		return nil, err
	}

	if resp.Status != "ok" {
		return nil, errors.Errorf("cancel order fail %+v", resp)
	}

	return &resp, nil
}

func (rc *RestClient) MatchResult(ctx context.Context, req *MatchResultReq) ([]MatchResult, error) {
	endPoint := fmt.Sprintf("/v1/order/orders/%s/matchresults", req.OrderID)
	var resp []MatchResult

	if err := rc.Request(ctx, http.MethodGet, endPoint, nil, nil, true, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (rc *RestClient) FetchOrder(ctx context.Context, order *exchange.Order) (*exchange.Order, error) {
	req := NewOrdersReq(order.ID.String())
	resp, err := rc.Orders(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Transform()
}

func (r *OrdersResp) Transform() (*exchange.Order, error) {
	symbol, err := ParseSymbol(r.Data.Symbol)
	if err != nil {
		return nil, err
	}

	amount, err := decimal.NewFromString(r.Data.Amount)
	if err != nil {
		return nil, err
	}
	price, err := decimal.NewFromString(r.Data.Price)
	if err != nil {
		return nil, err
	}
	filled, err := decimal.NewFromString(r.Data.FilledAmount)
	if err != nil {
		return nil, err
	}
	fees, err := decimal.NewFromString(r.Data.FilledFees)
	if err != nil {
		return nil, err
	}
	cost, err := decimal.NewFromString(r.Data.FilledCashAmount)
	var avgPrice decimal.Decimal
	if !filled.IsZero() {
		avgPrice = cost.Div(amount)
	}

	var ut time.Time
	ct := huobi.ParseTS(r.Data.CreatedAt)
	if r.Data.CanceledAt != 0 {
		ut = huobi.ParseTS(r.Data.CanceledAt)
	} else {
		ut = huobi.ParseTS(r.Data.FinishedAt)
	}

	status, err := ParseOrderStatus(r.Data.State)
	if err != nil {
		return nil, err
	}

	side, typ, err := ParseOrderType(r.Data.Type)
	if err != nil {
		return nil, err
	}

	return &exchange.Order{
		ID:       exchange.NewIntID(r.Data.ID),
		Symbol:   symbol,
		Price:    price,
		Amount:   amount,
		Filled:   filled,
		AvgPrice: avgPrice,
		Fee:      fees,
		Created:  ct,
		Updated:  ut,
		Side:     side,
		Status:   status,
		Type:     typ,
		Raw:      r,
	}, nil
}
